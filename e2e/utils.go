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

	wormholetypes "github.com/noble-assets/wormhole/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/relayer/rly"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/stretchr/testify/require"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"
	"go.uber.org/zap/zaptest"

	portaltypes "dollar.noble.xyz/v2/types/portal"
	"dollar.noble.xyz/v2/utils"
)

// Suite is a utility for spinning up a new E2E testing suite.
func Suite(t *testing.T, ibcEnabled bool) (ctx context.Context, noble *cosmos.CosmosChain, ibcSimapp *cosmos.CosmosChain, guardians []utils.Guardian) {
	ctx = context.Background()
	logger := zaptest.NewLogger(t)
	reporter := testreporter.NewNopReporter()
	execReporter := reporter.RelayerExecReporter(t)
	client, network := interchaintest.DockerSetup(t)
	var relayer *rly.CosmosRelayer

	guardian := utils.NewGuardian(t)
	guardians = []utils.Guardian{guardian}

	numValidators, numFullNodes := 1, 0

	specs := []*interchaintest.ChainSpec{
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
						Transceiver: utils.SourceTransceiverAddress,
						Manager:     utils.SourceManagerAddress,
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
	}
	if ibcEnabled {
		specs = append(specs, &interchaintest.ChainSpec{
			Name:          "ibc-go-simd",
			Version:       "v8.7.0",
			NumValidators: &numValidators,
			NumFullNodes:  &numFullNodes,
			ChainConfig: ibc.ChainConfig{
				ChainID: "ibc-go-simd-1",
			},
		})
	}
	factory := interchaintest.NewBuiltinChainFactory(logger, specs)

	chains, err := factory.Chains(t.Name())
	require.NoError(t, err)

	noble = chains[0].(*cosmos.CosmosChain)
	interchain := interchaintest.NewInterchain().AddChain(noble)
	if ibcEnabled {
		relayer = interchaintest.NewBuiltinRelayerFactory(
			ibc.CosmosRly,
			logger,
		).Build(t, client, network).(*rly.CosmosRelayer)

		ibcSimapp = chains[1].(*cosmos.CosmosChain)

		interchain = interchain.
			AddChain(ibcSimapp).
			AddRelayer(relayer, "relayer").
			AddLink(interchaintest.InterchainLink{
				Chain1:  noble,
				Chain2:  ibcSimapp,
				Relayer: relayer,
				Path:    "transfer",
			})
	}
	require.NoError(t, interchain.Build(ctx, execReporter, interchaintest.InterchainBuildOptions{
		TestName:  t.Name(),
		Client:    client,
		NetworkID: network,
	}))

	t.Cleanup(func() {
		_ = interchain.Close()
	})

	if ibcEnabled {
		require.NoError(t, relayer.StartRelayer(ctx, execReporter))
	}

	return
}
