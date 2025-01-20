package ntt

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

// EncodeNativeTokenTransfer is a utility that encodes a native token transfer.
func EncodeNativeTokenTransfer(ntt NativeTokenTransfer) (bz []byte) {
	bz = append(bz, NativeTokenTransferPrefix...)
	bz = append(bz, 6)
	bz = binary.BigEndian.AppendUint64(bz, ntt.Amount)
	bz = append(bz, ntt.SourceToken...)
	bz = append(bz, ntt.To...)
	bz = binary.BigEndian.AppendUint16(bz, ntt.ToChain)

	if len(ntt.AdditionalPayload) > 0 {
		bz = binary.BigEndian.AppendUint16(bz, uint16(len(ntt.AdditionalPayload)))
		bz = append(bz, ntt.AdditionalPayload...)
	}

	return
}

// ParseNativeTokenTransfer is a utility that parses a native token transfer.
func ParseNativeTokenTransfer(bz []byte) (transfer NativeTokenTransfer, err error) {
	if len(bz) < 79 {
		return NativeTokenTransfer{}, errors.New("native token transfer is malformed")
	}

	offset := 0

	prefix := bz[offset : offset+4]
	if !bytes.Equal(prefix, NativeTokenTransferPrefix) {
		return NativeTokenTransfer{}, fmt.Errorf(
			"native token transfer prefix is invalid: expected %s, got %s",
			hex.EncodeToString(NativeTokenTransferPrefix),
			hex.EncodeToString(prefix),
		)
	}
	offset += 4

	// NOTE: We ignore the number of decimals as it's assumed to be 6.
	offset += 1

	transfer.Amount = binary.BigEndian.Uint64(bz[offset : offset+8])
	offset += 8

	transfer.SourceToken = bz[offset : offset+32]
	offset += 32

	transfer.To = bz[offset : offset+32]
	offset += 32

	transfer.ToChain = binary.BigEndian.Uint16(bz[offset : offset+2])
	offset += 2

	if len(bz) > 79 {
		if len(bz) < 81 {
			return NativeTokenTransfer{}, errors.New("native token transfer is malformed")
		}

		payloadLength := binary.BigEndian.Uint16(bz[offset : offset+2])
		offset += 2
		transfer.AdditionalPayload = bz[offset:]

		if len(transfer.AdditionalPayload) != int(payloadLength) {
			return NativeTokenTransfer{}, errors.New("native token transfer additional payload is invalid length")
		}
	}

	return
}
