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

import {Test} from "forge-std/Test.sol";
import {console} from "forge-std/console.sol";

import {NoopIsm} from "@hyperlane/isms/NoopIsm.sol";
import {Message} from "@hyperlane/libs/Message.sol";
import {TokenMessage} from "@hyperlane/token/libs/TokenMessage.sol";
import {TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

import {NobleDollar} from "../src/NobleDollar.sol";

contract NobleDollarTest is Test {
    NobleDollar public usdn;

    address constant MAILBOX = 0xc005dc82818d67AF737725bD4bf75435d065D239;
    address constant USER1 = 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045;
    address constant USER2 = 0xF2f1ACbe0BA726fEE8d75f3E32900526874740BB;

    function setUp() public {
        NobleDollar implementation = new NobleDollar(MAILBOX);
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(implementation),
            msg.sender,
            abi.encodeWithSelector(
                NobleDollar.initialize.selector,
                address(0x9e6B1022bE9BBF5aFd152483DAD9b88911bC8611),
                address(new NoopIsm())
            )
        );
        usdn = NobleDollar(address(proxy));

        uint32[] memory domains = new uint32[](1);
        domains[0] = 1313817164;
        bytes32[] memory routers = new bytes32[](1);
        routers[0] = 0x726f757465725f61707000000000000000000000000000010000000000000000;
        usdn.enrollRemoteRouters(domains, routers);
    }

    function test() public {
        // ACT: Transfer of 1M $USDN from Noble Core to USER1.
        bytes memory mintPayload = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000004e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000d8da6bf26964af9d7eed9e03e53415d37aa96045000000000000000000000000000000000000000000000000000000e8d4a51000"
        );
        (bool mintSuccess,) = MAILBOX.call(mintPayload);

        // ASSERT: The transfer was successful, USER1 has a balance of 1M $USDN with a principal of 1M.
        assertEq(mintSuccess, true);

        assertEq(usdn.index(), 1e12);
        assertEq(usdn.totalSupply(), 1e12);
        assertEq(usdn.totalPrincipal(), 1e12);
        assertEq(usdn.balanceOf(USER1), 1e12);
        assertEq(usdn.principalOf(USER1), 1e12);
        assertEq(usdn.yield(USER1), 0);
        assertEq(usdn.balanceOf(USER2), 0);
        assertEq(usdn.principalOf(USER2), 0);
        assertEq(usdn.yield(USER2), 0);

        // ACT: Yield accrual of 111.506849 $USDN, 1 day's worth of 4.07% yield.
        bytes memory yieldPayload = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000014e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a0000000000000000000000000000000000000000000000000000000006a575a1"
        );
        (bool yieldSuccess,) = MAILBOX.call(yieldPayload);

        // ASSERT: The yield accrual was successful, USER1 has 111.506849 $USDN of claimable yield.
        assertEq(yieldSuccess, true);

        assertEq(usdn.index(), 1000111506849);
        assertEq(usdn.totalSupply(), 1000111506849);
        assertEq(usdn.totalPrincipal(), 1e12);
        assertEq(usdn.balanceOf(USER1), 1000000000000);
        assertEq(usdn.principalOf(USER1), 1000000000000);
        assertEq(usdn.yield(USER1), 111506849);
        assertEq(usdn.balanceOf(USER2), 0);
        assertEq(usdn.principalOf(USER2), 0);
        assertEq(usdn.yield(USER2), 0);

        // ACT: Transfer of 500k $USDN from USER1 to USER2.
        vm.prank(USER1);
        usdn.transfer(USER2, 5e11);

        // ASSERT: The transfer was successful.
        assertEq(usdn.index(), 1000111506849);
        assertEq(usdn.totalSupply(), 1000111506849);
        assertEq(usdn.totalPrincipal(), 1e12);
        assertEq(usdn.balanceOf(USER1), 5e11);
        assertEq(usdn.principalOf(USER1), 500055747208);
        assertEq(usdn.yield(USER1), 111506848);
        assertEq(usdn.balanceOf(USER2), 5e11);
        assertEq(usdn.principalOf(USER2), 499944252792);
        assertEq(usdn.yield(USER2), 0);

        // ACT: Claim yield for USER1.
        vm.prank(USER1);
        usdn.claim();

        // ASSERT: The yield was claimed.
        assertEq(usdn.index(), 1000111506849);
        assertEq(usdn.totalSupply(), 1000111506849);
        assertEq(usdn.totalPrincipal(), 1e12);
        assertEq(usdn.balanceOf(USER1), 500111506848);
        assertEq(usdn.principalOf(USER1), 500055747208);
        assertEq(usdn.yield(USER1), 0);
        assertEq(usdn.balanceOf(USER2), 5e11);
        assertEq(usdn.principalOf(USER2), 499944252792);
        assertEq(usdn.yield(USER2), 0);

        // ACT: Yield accrual of 111.506849 $USDN, 1 day's worth of 4.07% yield.
        bytes memory yieldPayload2 = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000024e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a0000000000000000000000000000000000000000000000000000000006a575a1"
        );
        (bool yieldSuccess2,) = MAILBOX.call(yieldPayload2);

        // ASSERT: The yield accrual was successful.
        assertEq(yieldSuccess2, true);

        assertEq(usdn.index(), 1000223013698);
        assertEq(usdn.totalSupply(), 1000223013698);
        assertEq(usdn.totalPrincipal(), 1e12);
        assertEq(usdn.balanceOf(USER1), 500111506848);
        assertEq(usdn.principalOf(USER1), 500055747208);
        assertEq(usdn.yield(USER1), 55759641);
        assertEq(usdn.balanceOf(USER2), 5e11);
        assertEq(usdn.principalOf(USER2), 499944252792);
        assertEq(usdn.yield(USER2), 55747208);
    }
}
