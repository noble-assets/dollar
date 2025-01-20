package ntt

import (
	"encoding/binary"
	"errors"
)

// EncodeManagerMessage is a utility that encodes a manager message.
func EncodeManagerMessage(msg ManagerMessage) (bz []byte) {
	bz = append(bz, msg.Id...)
	bz = append(bz, msg.Sender...)

	bz = binary.BigEndian.AppendUint16(bz, uint16(len(msg.Payload)))
	bz = append(bz, msg.Payload...)

	return
}

// ParseManagerMessage is a utility that parses a manager message.
func ParseManagerMessage(bz []byte) (msg ManagerMessage, err error) {
	if len(bz) < 66 {
		return ManagerMessage{}, errors.New("manager message is malformed")
	}

	msg.Id = bz[:32]
	msg.Sender = bz[32:64]

	payloadLength := binary.BigEndian.Uint16(bz[64:66])
	msg.Payload = bz[66:]
	if len(msg.Payload) != int(payloadLength) {
		return ManagerMessage{}, errors.New("manager message payload is invalid length")
	}

	return
}
