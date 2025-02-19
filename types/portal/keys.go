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

package portal

const SubmoduleName = "dollar/portal"

// NOTE: These variables are initialized when creating the module keeper.
var (
	PaddedTransceiverAddress = make([]byte, 32)
	TransceiverAddress       = ""
	PaddedManagerAddress     = make([]byte, 32)
	ManagerAddress           = ""
	RawToken                 = make([]byte, 32)
)

var (
	OwnerKey                    = []byte("portal/owner")
	PausedKey                   = []byte("portal/paused")
	PeerPrefix                  = []byte("portal/peer/")
	SupportedBridgingPathPrefix = []byte("portal/bridging_path/")
	NonceKey                    = []byte("portal/nonce")
)
