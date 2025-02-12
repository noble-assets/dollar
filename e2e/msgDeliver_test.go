package e2e

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/math"
	portaltypes "dollar.noble.xyz/types/portal"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/stretchr/testify/require"
)

// TestMsgDeliverInjection ensurs that MsgDeliverInjection is not public.
func TestMsgDeliverInjection(t *testing.T) {
	ctx, _, chain := Suite(t)

	broadcaster := cosmos.NewBroadcaster(t, chain)
	user := interchaintest.GetAndFundTestUsers(t, ctx, "wallet", math.OneInt(), chain)[0]

	bCtx, bCancel := context.WithTimeout(ctx, 20*time.Second)
	defer bCancel()

	mockVAA := "AQAAAAABAPel1AcBA57rIzaTw70Qqlta9SxhuBYByiTv3viGqwgfFq4Wfx/EN0Mb8D71aTIwBz36NUmI98Q2fCEQyFlFSqQAZ1vRXAAAAAAnEgAAAAAAAAAAAAAAAHsb16a05hwqEjrGvCy/xhRDfQRwAAAAAAAAsrwPAScUAAAAAAAAAAAAAAAAKcvx4HFm0xRGMHrgeZn6bRYiOZAAAADjmUX/EAAAAAAAAAAAAAAAABt64ZSyDFVbnZmcg190zc42pnp0AAAAAAAAAAAAAAAAG3rhlLIMVVudmZyDX3TNzjamenQAmwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABrAAAAAAAAAAAAAAAAlO0ORGvBexrFFbj2bBk6ZU0driQAWZlOVFQGAAAAAAAAJw8AAAAAAAAAAAAAAAAMlBrZTKSlLtrqvyA7Yb3RgHzuwAAAAAAAAAAAAAAAAJTtDkRrwXsaxRW49mwZOmVNHa4kJxQACAAAAO++YLGeAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD0JAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGpv1m2ycUAAAAAAAAAAAAAAAAlO0ORGvBexrFFbj2bBk6ZU0driQAAAAAAAAAAAAAAAB6ClOEd3b36UzDV0KXGssiF7DbgQAAAAAAAAAAAAAAAHoKU4R3dvfpTMNXQpcayyIXsNuBAAAAAAAAAAAAAAAAKcvx4HFm0xRGMHrgeZn6bRYiOZAA"

	_, err := cosmos.BroadcastTx(
		bCtx,
		broadcaster,
		user,
		&portaltypes.MsgDeliverInjection{
			Vaa: []byte(mockVAA),
		},
	)

	require.Error(t, err)
	require.ErrorContains(t, err, "no message handler found")

	// TODO: Query for the VAA and ensure it was not processed
}
