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

package keeper

import (
	"context"
	"fmt"
	"strconv"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	hyperlaneutil "github.com/bcp-innovations/hyperlane-cosmos/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	solomachine "github.com/cosmos/ibc-go/v8/modules/light-clients/06-solomachine"
	tendermint "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"

	"dollar.noble.xyz/v2/types"
	"dollar.noble.xyz/v2/types/v2"
)

var _ v2.QueryServer = &queryServerV2{}

type queryServerV2 struct {
	*Keeper
}

func NewQueryServerV2(keeper *Keeper) v2.QueryServer {
	return &queryServerV2{Keeper: keeper}
}

func (k queryServerV2) Stats(ctx context.Context, req *v2.QueryStats) (*v2.QueryStatsResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	stats, err := k.Keeper.Stats.Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get stats from state")
	}

	rawTotalExternalYield, err := k.GetTotalExternalYield(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get total external yield from state")
	}

	totalExternalYield := make(map[string]v2.QueryStatsResponse_ExternalYield)
	for key, rawAmount := range rawTotalExternalYield {
		provider, identifier := v2.ParseYieldRecipientKey(key)

		chainId := "UNKNOWN"
		switch provider {
		case v2.Provider_IBC:
			chainId = k.getIBCChainId(ctx, identifier)
		case v2.Provider_HYPERLANE:
			chainId = k.getHyperlaneChainId(ctx, identifier)
		}

		amount, _ := math.NewIntFromString(rawAmount)

		totalExternalYield[key] = v2.QueryStatsResponse_ExternalYield{
			ChainId: chainId,
			Amount:  amount,
		}
	}

	return &v2.QueryStatsResponse{
		TotalHolders:       stats.TotalHolders,
		TotalPrincipal:     stats.TotalPrincipal,
		TotalYieldAccrued:  stats.TotalYieldAccrued,
		TotalExternalYield: totalExternalYield,
	}, nil
}

func (k queryServerV2) YieldRecipients(ctx context.Context, req *v2.QueryYieldRecipients) (*v2.QueryYieldRecipientsResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	yieldRecipients, err := k.GetYieldRecipients(ctx)

	return &v2.QueryYieldRecipientsResponse{YieldRecipients: yieldRecipients}, err
}

func (k queryServerV2) YieldRecipient(ctx context.Context, req *v2.QueryYieldRecipient) (*v2.QueryYieldRecipientResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	key := collections.Join(int32(req.Provider), req.Identifier)
	yieldRecipient, err := k.Keeper.YieldRecipients.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("unable to find yield recipient for provider %s with identifier %s", req.Provider, req.Identifier)
	}

	return &v2.QueryYieldRecipientResponse{YieldRecipient: yieldRecipient}, nil
}

func (k queryServerV2) RetryAmounts(ctx context.Context, req *v2.QueryRetryAmounts) (*v2.QueryRetryAmountsResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	retryAmounts, err := k.GetRetryAmounts(ctx)

	return &v2.QueryRetryAmountsResponse{RetryAmounts: retryAmounts}, err
}

func (k queryServerV2) RetryAmount(ctx context.Context, req *v2.QueryRetryAmount) (*v2.QueryRetryAmountResponse, error) {
	if req == nil {
		return nil, types.ErrInvalidRequest
	}

	retryAmount := k.GetRetryAmount(ctx, req.Provider, req.Identifier)

	return &v2.QueryRetryAmountResponse{RetryAmount: retryAmount}, nil
}

func (k *Keeper) getIBCChainId(ctx context.Context, channelId string) string {
	_, rawClientState, _ := k.channel.GetChannelClientState(sdk.UnwrapSDKContext(ctx), transfertypes.PortID, channelId)

	switch clientState := rawClientState.(type) {
	case *solomachine.ClientState:
		return clientState.ConsensusState.Diversifier
	case *tendermint.ClientState:
		return clientState.ChainId
	default:
		return "UNKNOWN"
	}
}

func (k *Keeper) getHyperlaneChainId(ctx context.Context, identifier string) string {
	rawIdentifier, err := hyperlaneutil.DecodeHexAddress(identifier)
	if err != nil {
		return "UNKNOWN"
	}
	tokenId := rawIdentifier.GetInternalId()

	router, err := k.getHyperlaneRouter(ctx, tokenId)
	if err != nil {
		return "UNKNOWN"
	}

	return strconv.Itoa(int(router.ReceiverDomain))
}
