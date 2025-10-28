/*
 * Copyright 2025 NASD Inc. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
pragma solidity 0.8.30;

import {NobleDollar as BaseNobleDollar} from "./NobleDollar.sol";

/**
 * @title  NobleDollar
 * @author John Letey <john@noble.xyz>
 * @notice ERC20 Noble Dollar on HyperEVM.
 */
contract NobleDollar is BaseNobleDollar {
    /// @notice The address of the Hyperliquid bridge for this token
    address public bridge;

    constructor(address mailbox_, uint32 tokenId_) BaseNobleDollar(mailbox_) {
        // According to the Hyperliquid documentation, this is how you derive the
        // bridge address for a linked HyperCore <> HyperEVM token.
        //
        // "Every token has a system address on the Core, which is the address
        // with first byte 0x20 and the remaining bytes all zeros, except for
        // the token index encoded in big-endian format."
        uint160 base = uint160(0x20) << 152;
        bridge = address(uint160(base | uint160(tokenId_)));
    }

    /// @notice Claims all available yield for the bridge and transfers it to the owner.
    function claimForBridge() public {
        uint256 amount = claim(bridge);

        _update(bridge, owner(), amount);
    }
}
