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

package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"dollar.noble.xyz/v3/types"
	"dollar.noble.xyz/v3/types/v2"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(QueryStats())
	cmd.AddCommand(QueryYieldRecipients())
	cmd.AddCommand(QueryYieldRecipient())
	cmd.AddCommand(QueryRetryAmounts())
	cmd.AddCommand(QueryRetryAmount())

	return cmd
}

func QueryStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Execute the Stats RPC method",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := v2.NewQueryClient(clientCtx)

			res, err := queryClient.Stats(context.Background(), &v2.QueryStats{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryYieldRecipients() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "yield-recipients",
		Short: "Query all yield recipients for external chains",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := v2.NewQueryClient(clientCtx)

			res, err := queryClient.YieldRecipients(context.Background(), &v2.QueryYieldRecipients{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryYieldRecipient() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "yield-recipient [provider] [identifier]",
		Short: "Query the yield recipient for an external chain",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := v2.NewQueryClient(clientCtx)

			provider, err := parseProvider(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.YieldRecipient(context.Background(), &v2.QueryYieldRecipient{
				Provider:   provider,
				Identifier: args[1],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryRetryAmounts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "retry-amounts",
		Short: "Query all retry amounts for external chains",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := v2.NewQueryClient(clientCtx)

			res, err := queryClient.RetryAmounts(context.Background(), &v2.QueryRetryAmounts{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryRetryAmount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "retry-amount [provider] [identifier]",
		Short: "Query the retry amount for an external chain",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := v2.NewQueryClient(clientCtx)

			provider, err := parseProvider(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.RetryAmount(context.Background(), &v2.QueryRetryAmount{
				Provider:   provider,
				Identifier: args[1],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
