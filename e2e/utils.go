// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, NASD Inc. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package e2e

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	wormholetypes "github.com/noble-assets/wormhole/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/stretchr/testify/require"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	portaltypes "dollar.noble.xyz/types/portal"
	"dollar.noble.xyz/utils"
)

// Suite sets up a test suite with a single chain.
func Suite(t *testing.T) (ctx context.Context, logger *zap.Logger, chain *cosmos.CosmosChain, guardians []utils.Guardian) {
	ctx = context.Background()
	logger = zaptest.NewLogger(t)

	numValidators, numFullNodes := 1, 0

	guardian := utils.NewGuardian()
	guardians = []utils.Guardian{guardian}

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "dollar",
			Version:       "local",
			NumValidators: &numValidators,
			NumFullNodes:  &numFullNodes,
			ChainConfig: ibc.ChainConfig{
				Type:    "cosmos",
				Name:    "dollar",
				ChainID: "dollar-1",
				Images: []ibc.DockerImage{
					{
						Repository: "noble-dollar-simd",
						Version:    "local",
						UIDGID:     "1025:1025",
					},
				},
				Bin:            "simd",
				Bech32Prefix:   "noble",
				Denom:          "uusdc",
				GasPrices:      "0uusdc",
				GasAdjustment:  1.5,
				TrustingPeriod: "504h",
				ModifyGenesis: func(cc ibc.ChainConfig, genesis []byte) ([]byte, error) {
					peers := make(map[uint16]portaltypes.Peer)
					peers[uint16(vaautils.ChainIDEthereum)] = portaltypes.Peer{
						// https://github.com/m0-foundation/m-portal/blob/dbe93da561c94dfc04beec8a144b11b287957b7a/deployments/noble/1.json#L3
						Transceiver: common.FromHex("0x000000000000000000000000c7dd372c39e38bf11451ab4a8427b4ae38cef644"),
						// https://github.com/m0-foundation/m-portal/blob/dbe93da561c94dfc04beec8a144b11b287957b7a/deployments/noble/1.json#L2
						Manager: common.FromHex("0x00000000000000000000000083ae82bd4054e815fb7b189c39d9ce670369ea16"),
					}

					guardianSets := make(map[uint16]wormholetypes.GuardianSet)
					var addresses [][]byte
					for _, guardian := range guardians {
						addresses = append(addresses, guardian.Address.Bytes())
					}
					guardianSets[0] = wormholetypes.GuardianSet{
						Addresses:      addresses,
						ExpirationTime: 0,
					}

					updatedGenesis := []cosmos.GenesisKV{
						cosmos.NewGenesisKV("app_state.dollar.portal.peers", peers),
						cosmos.NewGenesisKV("app_state.wormhole.config", wormholetypes.Config{
							ChainId:          uint16(vaautils.ChainIDNoble),
							GuardianSetIndex: 0,
							GovChain:         uint16(vaautils.GovernanceChain),
							GovAddress:       vaautils.GovernanceEmitter.Bytes(),
						}),
						cosmos.NewGenesisKV("app_state.wormhole.guardian_sets", guardianSets),
					}

					return cosmos.ModifyGenesis(updatedGenesis)(cc, genesis)
				},
			},
		},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	chain = chains[0].(*cosmos.CosmosChain)

	ic := interchaintest.NewInterchain().
		AddChain(chain)

	client, network := interchaintest.DockerSetup(t)

	require.NoError(t, ic.Build(ctx, nil, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: true,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})

	return ctx, logger, chain, guardians
}
