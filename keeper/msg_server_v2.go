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

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"

	"dollar.noble.xyz/v2/types"
	"dollar.noble.xyz/v2/types/v2"
)

var _ v2.MsgServer = &msgServerV2{}

type msgServerV2 struct {
	*Keeper
}

func NewMsgServerV2(keeper *Keeper) v2.MsgServer {
	return &msgServerV2{Keeper: keeper}
}

func (k msgServerV2) SetYieldRecipient(ctx context.Context, msg *v2.MsgSetYieldRecipient) (*v2.MsgSetYieldRecipientResponse, error) {
	if msg.Signer != k.authority {
		return nil, errors.Wrapf(types.ErrInvalidAuthority, "expected %s, got %s", k.authority, msg.Signer)
	}

	switch msg.Provider {
	case v2.Provider_IBC:
		_, found := k.channel.GetChannel(sdk.UnwrapSDKContext(ctx), transfertypes.PortID, msg.Identifier)
		if !found {
			return nil, fmt.Errorf("ibc identifier does not exist: %s", msg.Identifier)
		}
	}

	key := collections.Join(int32(msg.Provider), msg.Identifier)
	err := k.YieldRecipients.Set(ctx, key, msg.Recipient)
	if err != nil {
		return nil, errors.Wrap(err, "unable to set yield recipient in state")
	}

	return &v2.MsgSetYieldRecipientResponse{}, k.event.EventManager(ctx).Emit(ctx, &v2.YieldRecipientSet{
		Provider:   msg.Provider,
		Identifier: msg.Identifier,
		Recipient:  msg.Recipient,
	})
}
