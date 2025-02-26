// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, NASD Inc. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package keeper

import (
	"bytes"
	"context"
	"encoding/binary"
	"strconv"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/event"
	"cosmossdk.io/errors"
	"google.golang.org/protobuf/runtime/protoiface"

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
	if err := k.Keeper.Deliver(ctx, msg.Vaa); err != nil {
		return nil, err
	}

	return &portal.MsgDeliverResponse{}, k.event.EventManager(ctx).Emit(ctx, &portal.Delivered{
		Vaa: msg.Vaa,
	})
}

func (k portalMsgServer) Transfer(ctx context.Context, msg *portal.MsgTransfer) (*portal.MsgTransferResponse, error) {
	if k.GetPortalPaused(ctx) {
		return nil, portal.ErrPaused
	}

	peer, err := k.PortalPeers.Get(ctx, msg.DestinationChainId)
	if err != nil {
		return nil, errors.Wrapf(portal.ErrInvalidPeer, "chain %d is not configured", msg.DestinationChainId)
	}

	key := collections.Join(msg.DestinationChainId, msg.DestinationToken)
	if has, _ := k.PortalBridgingPaths.Has(ctx, key); !has {
		return nil, errors.Wrapf(portal.ErrInvalidBridgePath, "token %s is not configured for chain %d", msg.DestinationToken, msg.DestinationChainId)
	}

	if len(msg.Recipient) != 32 {
		return nil, errors.Wrap(portal.ErrInvalidRecipient, "recipient must be 32 bytes")
	}

	index, err := k.Index.Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get index from state")
	}

	additionalPayload := portal.EncodeAdditionalPayload(index, msg.DestinationToken)

	rawNativeTokenTransfer := ntt.EncodeNativeTokenTransfer(ntt.NativeTokenTransfer{
		Amount:            msg.Amount.Uint64(),
		SourceToken:       portal.RawToken,
		To:                msg.Recipient,
		ToChain:           msg.DestinationChainId,
		AdditionalPayload: additionalPayload,
	})

	sender, err := k.address.StringToBytes(msg.Signer)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to decode account %s", msg.Signer)
	}
	rawSender := make([]byte, 32)
	copy(rawSender[32-len(sender):], sender)

	nonce, err := k.IncrementPortalNonce(ctx)
	if err != nil {
		return nil, err
	}
	rawNonce := make([]byte, 4)
	binary.BigEndian.PutUint32(rawNonce, nonce)
	id := make([]byte, 32)
	copy(id[32-len(rawNonce):], rawNonce)

	managerMessage := ntt.ManagerMessage{
		Id:      id,
		Sender:  rawSender,
		Payload: rawNativeTokenTransfer,
	}

	chain, err := k.wormhole.GetChain(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get wormhole chain id")
	}
	rawManagerMessage := ntt.EncodeManagerMessage(managerMessage)
	messageId := ntt.ManagerMessageDigest(chain, managerMessage)

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

	err = k.wormhole.PostMessage(
		ctx,
		portal.TransceiverAddress,
		rawTransceiverMessage,
		nonce,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to post transfer message")
	}

	if err := k.event.EventManager(ctx).Emit(ctx, &portal.USDNTokenSent{
		SourceToken:        string(portal.RawToken),
		DestinationChainId: msg.DestinationChainId,
		DestinationToken:   msg.DestinationToken,
		Sender:             msg.Signer,
		Recipient:          msg.Recipient,
		Amount:             msg.Amount,
		Index:              index,
		MessageId:          messageId,
	}); err != nil {
		return nil, err
	}

	if err := k.event.EventManager(ctx).EmitKV(
		ctx,
		"transfer_sent",
		event.Attribute{Key: "recipient", Value: string(msg.Recipient)},
		event.Attribute{Key: "refund_address", Value: ""},
		event.Attribute{Key: "amount", Value: strconv.Itoa(int(msg.Amount.Int64()))},
		event.Attribute{Key: "fee", Value: ""},
		event.Attribute{Key: "recipient_chain", Value: strconv.Itoa(int(msg.DestinationChainId))},
		event.Attribute{Key: "msg_sequence", Value: strconv.Itoa(int(nonce))},
	); err != nil {
		return nil, err
	}

	return &portal.MsgTransferResponse{}, k.event.EventManager(ctx).EmitKV(
		ctx,
		"transfer_sent",
		event.Attribute{Key: "digest", Value: string(messageId)},
	)
}

func (k portalMsgServer) SetPausedState(ctx context.Context, msg *portal.MsgSetPausedState) (*portal.MsgSetPausedStateResponse, error) {
	if err := k.EnsureOwner(ctx, msg.Signer); err != nil {
		return nil, err
	}

	if err := k.PortalPaused.Set(ctx, msg.Paused); err != nil {
		return nil, err
	}

	event := protoiface.MessageV1(&portal.Unpaused{})
	if msg.Paused {
		event = &portal.Paused{}
	}
	return &portal.MsgSetPausedStateResponse{}, k.event.EventManager(ctx).Emit(ctx, event)
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

	peer, _ := k.PortalPeers.Get(ctx, msg.Chain)
	err = k.PortalPeers.Set(ctx, msg.Chain, portal.Peer{
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

func (k portalMsgServer) SetBridgingPath(ctx context.Context, msg *portal.MsgSetBridgingPath) (*portal.MsgSetBridgingPathResponse, error) {
	if err := k.EnsureOwner(ctx, msg.Signer); err != nil {
		return nil, err
	}

	chainID, err := k.wormhole.GetChain(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get wormhole chain id")
	}
	if msg.DestinationChainId == chainID {
		return nil, errors.Wrapf(portal.ErrInvalidBridgePath, "destination chain cannot be %d", chainID)
	}

	empty := make([]byte, 32)
	if len(msg.DestinationToken) != 32 {
		return nil, errors.Wrap(portal.ErrInvalidBridgePath, "destination token must be 32 bytes")
	}
	if bytes.Equal(msg.DestinationToken, empty) {
		return nil, errors.Wrap(portal.ErrInvalidBridgePath, "destination token must not be empty")
	}

	key := collections.Join(msg.DestinationChainId, msg.DestinationToken)
	err = k.PortalBridgingPaths.Set(ctx, key, msg.Supported)
	if err != nil {
		return nil, errors.Wrap(err, "unable to set bridging path")
	}

	return &portal.MsgSetBridgingPathResponse{}, k.event.EventManager(ctx).Emit(ctx, &portal.BridgingPathSet{
		DestinationChainId: msg.DestinationChainId,
		DestinationToken:   msg.DestinationToken,
		Supported:          msg.Supported,
	})
}

func (k portalMsgServer) TransferOwnership(ctx context.Context, msg *portal.MsgTransferOwnership) (*portal.MsgTransferOwnershipResponse, error) {
	if err := k.EnsureOwner(ctx, msg.Signer); err != nil {
		return nil, err
	}

	if _, err := k.address.StringToBytes(msg.NewOwner); err != nil {
		return nil, errors.Wrap(err, "unable to decode new owner address")
	}
	if msg.NewOwner == msg.Signer {
		return nil, portal.ErrSameOwner
	}

	if err := k.PortalOwner.Set(ctx, msg.NewOwner); err != nil {
		return nil, errors.Wrap(err, "unable to set owner in state")
	}

	return &portal.MsgTransferOwnershipResponse{}, k.event.EventManager(ctx).Emit(ctx, &portal.OwnershipTransferred{
		PreviousOwner: msg.Signer,
		NewOwner:      msg.NewOwner,
	})
}

// EnsureOwner is a utility that ensures a message was signed by the portal owner.
func (k portalMsgServer) EnsureOwner(ctx context.Context, signer string) error {
	owner, _ := k.PortalOwner.Get(ctx)
	if owner == "" {
		return portal.ErrNoOwner
	}

	if signer != owner {
		return errors.Wrapf(portal.ErrNotOwner, "expected %s, got %s", owner, signer)
	}

	return nil
}
