package ntt

import "encoding/hex"

var (
	// NativeTokenTransferPrefix is the prefix for all NativeTokenTransfer payloads.
	//
	// https://github.com/wormhole-foundation/native-token-transfers/blob/67df54701e0f4b3793b6c621719911804c9875a3/evm/src/libraries/TransceiverStructs.sol#L37-L39
	NativeTokenTransferPrefix, _ = hex.DecodeString("994e5454")
	// TransceiverMessagePrefix is the prefix for all TransceiverMessage payloads.
	//
	// https://github.com/wormhole-foundation/native-token-transfers/blob/67df54701e0f4b3793b6c621719911804c9875a3/evm/src/Transceiver/WormholeTransceiver/WormholeTransceiverState.sol#L38-L41
	TransceiverMessagePrefix, _ = hex.DecodeString("9945ff10")
)
