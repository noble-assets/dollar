package portal

import (
	"encoding/binary"
	"encoding/hex"

	"cosmossdk.io/math"

	"dollar.noble.xyz/types/portal/ntt"
)

type PayloadType int

const (
	Unknown PayloadType = iota
	Token
	Index
)

// GetPayloadType is a utility for determining the type of custom payload.
// Since we aren't implementing a registrar, the Key and List types are ignored.
//
// https://github.com/m0-foundation/m-portal/blob/ddf583b9bef971752ec1360f9b089e6fefa9c526/src/libs/PayloadEncoder.sol#L32-L43
func GetPayloadType(payload []byte) PayloadType {
	if len(payload) < 4 {
		return Unknown
	}

	switch hex.EncodeToString(payload[:4]) {
	// 0x99'N''T''T'
	case "994e5454":
		return Token
	// M0IT - M0 Index Transfer
	case "4d304954":
		return Index
	}

	return Unknown
}

// DecodeTokenPayload is a utility for decoding a custom payload of type Token.
//
// https://github.com/m0-foundation/m-portal/blob/ddf583b9bef971752ec1360f9b089e6fefa9c526/src/libs/PayloadEncoder.sol#L45-L62
func DecodeTokenPayload(payload []byte) (amount math.Int, index math.LegacyDec, recipient []byte, destination uint16) {
	ntt, _ := ntt.ParseNativeTokenTransfer(payload)

	rawIndex := binary.BigEndian.Uint64(ntt.AdditionalPayload)
	index = math.LegacyNewDec(int64(rawIndex)).QuoInt64(1e12)

	return math.NewIntFromUint64(ntt.Amount), index, ntt.To[12:], ntt.ToChain
}

// DecodeIndexPayload is a utility for decoding a custom payload of type Index.
//
// https://github.com/m0-foundation/m-portal/blob/ddf583b9bef971752ec1360f9b089e6fefa9c526/src/libs/PayloadEncoder.sol#L68-L75
func DecodeIndexPayload(payload []byte) (index math.LegacyDec, destination uint16) {
	rawIndex := binary.BigEndian.Uint64(payload[4:12])
	index = math.LegacyNewDec(int64(rawIndex)).QuoInt64(1e12)

	destination = binary.BigEndian.Uint16(payload[12:14])

	return
}
