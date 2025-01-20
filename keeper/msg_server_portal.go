package keeper

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/errors"

	"dollar.noble.xyz/types/portal"
	"dollar.noble.xyz/types/portal/ntt"
)

var _ portal.MsgServer = &portalMsgServer{}

type portalMsgServer struct {
	*Keeper
}

func NewPortalMsgServer(keeper *Keeper) portal.MsgServer {
	return &portalMsgServer{Keeper: keeper}
}

func (k portalMsgServer) Deliver(ctx context.Context, msg *portal.MsgDeliver) (*portal.MsgDeliverResponse, error) {
	vaa, err := k.wormhole.ParseAndVerifyVAA(ctx, msg.Vaa)
	if err != nil {
		return nil, err
	}

	peer, err := k.Peers.Get(ctx, uint16(vaa.EmitterChain))
	if err != nil {
		return nil, errors.Wrapf(portal.ErrInvalidPeer, "chain %d not configured", vaa.EmitterChain)
	}

	if !bytes.Equal(peer.Transceiver, vaa.EmitterAddress.Bytes()) {
		return nil, errors.Wrapf(
			portal.ErrInvalidPeer,
			"expected transceiver %s for chain %d, got %s",
			hex.EncodeToString(peer.Transceiver), vaa.EmitterChain,
			vaa.EmitterAddress.String(),
		)
	}

	transceiverMessage, err := ntt.ParseTransceiverMessage(vaa.Payload)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(peer.Manager, transceiverMessage.SourceManagerAddress) {
		return nil, errors.Wrapf(
			portal.ErrInvalidPeer,
			"expected manager %s for chain %d, got %s",
			hex.EncodeToString(peer.Manager), vaa.EmitterChain,
			hex.EncodeToString(transceiverMessage.SourceManagerAddress),
		)
	}

	if !bytes.Equal(portal.PaddedManagerAddress, transceiverMessage.RecipientManagerAddress) {
		return nil, errors.Wrapf(
			portal.ErrInvalidMessage,
			"expected recipient manager %s, got %s",
			hex.EncodeToString(portal.PaddedManagerAddress),
			hex.EncodeToString(transceiverMessage.RecipientManagerAddress),
		)
	}

	managerMessage, err := ntt.ParseManagerMessage(transceiverMessage.ManagerPayload)
	if err != nil {
		return nil, err
	}

	return &portal.MsgDeliverResponse{}, k.HandlePayload(ctx, managerMessage.Payload)
}

func (k portalMsgServer) SetPeer(ctx context.Context, msg *portal.MsgSetPeer) (*portal.MsgSetPeerResponse, error) {
	if err := k.EnsureOwner(ctx, msg.Signer); err != nil {
		return nil, err
	}

	if msg.Chain == 0 {
		return nil, errors.Wrap(portal.ErrInvalidPeer, "chain cannot be 0")
	}
	chain, err := k.wormhole.GetChain(ctx)
	if err != nil || msg.Chain == chain {
		return nil, errors.Wrapf(portal.ErrInvalidPeer, "chain cannot be %d", chain)
	}

	empty := make([]byte, 32)

	if len(msg.Transceiver) != 32 {
		return nil, errors.Wrap(portal.ErrInvalidPeer, "transceiver must be 32 bytes")
	}
	if bytes.Equal(msg.Transceiver, empty) {
		return nil, errors.Wrap(portal.ErrInvalidPeer, "transceiver must not be empty")
	}

	if len(msg.Manager) != 32 {
		return nil, errors.Wrap(portal.ErrInvalidPeer, "manager must be 32 bytes")
	}
	if bytes.Equal(msg.Manager, empty) {
		return nil, errors.Wrap(portal.ErrInvalidPeer, "manager must not be empty")
	}

	peer, _ := k.Peers.Get(ctx, msg.Chain)
	err = k.Peers.Set(ctx, msg.Chain, portal.Peer{
		Transceiver: msg.Transceiver,
		Manager:     msg.Manager,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to set peer in state")
	}

	return &portal.MsgSetPeerResponse{}, k.event.EventManager(ctx).Emit(ctx, &portal.PeerUpdated{
		Chain:          msg.Chain,
		OldTransceiver: peer.Transceiver,
		NewTransceiver: msg.Transceiver,
		OldManager:     peer.Manager,
		NewManager:     msg.Manager,
	})
}

func (k portalMsgServer) Transfer(ctx context.Context, msg *portal.MsgTransfer) (*portal.MsgTransferResponse, error) {
	peer, err := k.Peers.Get(ctx, msg.Chain)
	if err != nil {
		return nil, errors.Wrapf(portal.ErrInvalidPeer, "chain %d is not configured", msg.Chain)
	}

	if len(msg.Recipient) != 32 {
		return nil, errors.Wrap(portal.ErrInvalidRecipient, "recipient must be 32 bytes")
	}

	index, err := k.Index.Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get index from state")
	}
	additionalPayload := make([]byte, 8)
	binary.BigEndian.PutUint64(additionalPayload, uint64(index))

	rawNativeTokenTransfer := ntt.EncodeNativeTokenTransfer(ntt.NativeTokenTransfer{
		Amount:            msg.Amount.Uint64(),
		SourceToken:       portal.RawToken,
		To:                msg.Recipient,
		ToChain:           msg.Chain,
		AdditionalPayload: additionalPayload,
	})

	sender, err := k.address.StringToBytes(msg.Signer)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to decode account %s", msg.Signer)
	}
	rawSender := make([]byte, 32)
	copy(rawSender[32-len(sender):], sender)

	nonce, err := k.IncrementNonce(ctx)
	if err != nil {
		return nil, err
	}
	rawNonce := make([]byte, 4)
	binary.BigEndian.PutUint32(rawNonce, nonce)
	id := make([]byte, 32)
	copy(id[32-len(rawNonce):], rawNonce)

	rawManagerMessage := ntt.EncodeManagerMessage(ntt.ManagerMessage{
		Id:      id,
		Sender:  rawSender,
		Payload: rawNativeTokenTransfer,
	})

	rawTransceiverMessage := ntt.EncodeTransceiverMessage(ntt.TransceiverMessage{
		SourceManagerAddress:    portal.PaddedManagerAddress,
		RecipientManagerAddress: peer.Manager,
		ManagerPayload:          rawManagerMessage,
		TransceiverPayload:      nil,
	})

	err = k.Burn(ctx, sender, msg.Amount)
	if err != nil {
		return nil, errors.Wrap(err, "unable to burn coins")
	}

	return &portal.MsgTransferResponse{}, k.wormhole.PostMessage(
		ctx,
		portal.TransceiverAddress,
		rawTransceiverMessage,
		nonce,
	)
}

// EnsureOwner is a utility that ensures a message was signed by the portal owner.
func (k portalMsgServer) EnsureOwner(ctx context.Context, signer string) error {
	owner, _ := k.Owner.Get(ctx)
	if owner == "" {
		return portal.ErrNoOwner
	}

	if signer != owner {
		return errors.Wrapf(portal.ErrNotOwner, "expected %s, got %s", owner, signer)
	}

	return nil
}

// HandlePayload is a utility that handles custom payloads when delivering portal messages.
func (k portalMsgServer) HandlePayload(ctx context.Context, payload []byte) error {
	chain, _ := k.wormhole.GetChain(ctx)

	switch portal.GetPayloadType(payload) {
	case portal.Unknown:
		return nil
	case portal.Token:
		amount, index, recipient, destination := portal.DecodeTokenPayload(payload)
		if chain != destination {
			return fmt.Errorf("not destination chain: expected %d, got %d", chain, destination)
		}

		return k.Mint(ctx, recipient, amount, &index)
	case portal.Index:
		index, destination := portal.DecodeIndexPayload(payload)
		if chain != destination {
			return fmt.Errorf("not destination chain: expected %d, got %d", chain, destination)
		}

		return k.UpdateIndex(ctx, index)
	}

	return nil
}
