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
pragma solidity >=0.8.0;

import {Test, console} from "forge-std/Test.sol";
import {NobleDollar} from "../src/NobleDollar.sol";

contract NobleDollarTest is Test {
    NobleDollar public usdn;

    function setUp() public {
        // https://github.com/hyperlane-xyz/hyperlane-registry/blob/main/chains/ethereum/addresses.yaml
        usdn = new NobleDollar(
            address(0xc005dc82818d67AF737725bD4bf75435d065D239),
            address(0x9e6B1022bE9BBF5aFd152483DAD9b88911bC8611),
            address(0x1AB8c76BAD3829B46b738B61cC941b22bE82C16e)
        );
    }

    function test_TotalSupply() public view {
        assertEq(usdn.totalSupply(), 0);
    }
}
