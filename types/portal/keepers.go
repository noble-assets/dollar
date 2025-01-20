package portal

import (
	"context"

	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"
)

type WormholeKeeper interface {
	GetChain(ctx context.Context) (uint16, error)
	ParseAndVerifyVAA(ctx context.Context, bz []byte) (*vaautils.VAA, error)
	PostMessage(ctx context.Context, signer string, message []byte, nonce uint32) error
}
