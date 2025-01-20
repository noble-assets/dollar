package ntt

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

// EncodeTransceiverMessage is a utility that encodes a transceiver message.
func EncodeTransceiverMessage(msg TransceiverMessage) (bz []byte) {
	bz = append(bz, TransceiverMessagePrefix...)
	bz = append(bz, msg.SourceManagerAddress...)
	bz = append(bz, msg.RecipientManagerAddress...)

	bz = binary.BigEndian.AppendUint16(bz, uint16(len(msg.ManagerPayload)))
	bz = append(bz, msg.ManagerPayload...)

	bz = binary.BigEndian.AppendUint16(bz, uint16(len(msg.TransceiverPayload)))
	bz = append(bz, msg.TransceiverPayload...)

	return
}

// ParseTransceiverMessage is a utility that parses a transceiver message.
func ParseTransceiverMessage(bz []byte) (msg TransceiverMessage, err error) {
	offset := 0

	prefix := bz[offset : offset+4]
	if !bytes.Equal(prefix, TransceiverMessagePrefix) {
		return TransceiverMessage{}, fmt.Errorf(
			"transceiver message prefix is invalid: expected %s, got %s",
			hex.EncodeToString(TransceiverMessagePrefix),
			hex.EncodeToString(prefix),
		)
	}
	offset += 4

	msg.SourceManagerAddress = bz[offset : offset+32]
	offset += 32

	msg.RecipientManagerAddress = bz[offset : offset+32]
	offset += 32

	managerPayloadLength := int(binary.BigEndian.Uint16(bz[offset : offset+2]))
	offset += 2
	msg.ManagerPayload = bz[offset : offset+managerPayloadLength]
	offset += managerPayloadLength

	transceiverPayloadLength := int(binary.BigEndian.Uint16(bz[offset : offset+2]))
	offset += 2
	msg.TransceiverPayload = bz[offset : offset+transceiverPayloadLength]

	return
}
