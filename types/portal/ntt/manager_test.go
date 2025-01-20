package ntt_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"dollar.noble.xyz/types/portal/ntt"
)

func TestManagerMessageUtilities(t *testing.T) {
	raw, err := hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000000000000000000000000000090f8bf6a479f320ead074411a4b0e7944ea8c9c10059994e54540600000000000f42400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000153b01572d58e7d24a4611b2bb31d9e8de9535a70fa900080000000000000000")
	require.NoError(t, err, "unable to decode input")

	msg, err := ntt.ParseManagerMessage(raw)
	require.NoError(t, err, "unable to parse manager message")

	require.Equal(t, make([]byte, 32), msg.Id)
	require.Equal(t, "90f8bf6a479f320ead074411a4b0e7944ea8c9c1", hex.EncodeToString(msg.Sender[12:]))
	require.Equal(t, "994e54540600000000000f42400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000153b01572d58e7d24a4611b2bb31d9e8de9535a70fa900080000000000000000", hex.EncodeToString(msg.Payload))

	bz := ntt.EncodeManagerMessage(msg)
	require.Equal(t, raw, bz)
}
