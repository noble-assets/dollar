package ntt_test

import (
	"encoding/hex"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/stretchr/testify/require"

	"dollar.noble.xyz/types/portal/ntt"
)

// MustEncodeRecipient is a utility that parses a universal address to a Noble
// address. NOTE: This should only be used for testing purposes!
func MustEncodeRecipient(bz []byte) string {
	cdc := address.NewBech32Codec("noble")
	recipient, _ := cdc.BytesToString(bz[12:])
	return recipient
}

func TestTransceiverMessageUtilities(t *testing.T) {
	raw, err := hex.DecodeString("9945ff1000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002e859506ba229c183f8985d54fe7210923fb9bca009b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000090f8bf6a479f320ead074411a4b0e7944ea8c9c10059994e54540600000000000f42400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000153b01572d58e7d24a4611b2bb31d9e8de9535a70fa9000800000000000000000000")
	require.NoError(t, err, "unable to decode input")

	msg, err := ntt.ParseTransceiverMessage(raw)
	require.NoError(t, err, "unable to parse transceiver message")

	require.Equal(t, make([]byte, 32), msg.SourceManagerAddress)
	require.Equal(t, "noble196ze2p46y2wps0ufsh25leeppy3lhx72snqjun", MustEncodeRecipient(msg.RecipientManagerAddress))
	require.Equal(t, "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000090f8bf6a479f320ead074411a4b0e7944ea8c9c10059994e54540600000000000f42400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000153b01572d58e7d24a4611b2bb31d9e8de9535a70fa900080000000000000000", hex.EncodeToString(msg.ManagerPayload))
	require.Equal(t, make([]byte, 0), msg.TransceiverPayload)

	bz := ntt.EncodeTransceiverMessage(msg)
	require.Equal(t, raw, bz)
}
