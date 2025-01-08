package ntt

import (
	"encoding/binary"
	"errors"
)

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
