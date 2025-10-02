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
        assertTrue(mintSuccess, "Initial mint should succeed");

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
        assertTrue(yieldSuccess, "Yield accrual should succeed");

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
        assertTrue(yieldSuccess2, "Second yield accrual should succeed");

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

    function test_transferToUSDNFromNonZeroAccountReverts() public {
        // ACT: Transfer of 1M $USDN from Noble Core to USER1.
        bytes memory mintPayload = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000004e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000d8da6bf26964af9d7eed9e03e53415d37aa96045000000000000000000000000000000000000000000000000000000e8d4a51000"
        );
        (bool mintSuccess,) = MAILBOX.call(mintPayload);
        assertTrue(mintSuccess, "Initial mint should succeed");

        uint256 _user1Balance = usdn.balanceOf(USER1);

        assertEq(_user1Balance, 1000000e6, "user 1 should have 1 million usdn");

        vm.expectRevert(abi.encodeWithSelector(NobleDollar.InvalidTransfer.selector));

        vm.prank(USER1);
        usdn.transfer(address(usdn), 1000e6);
    }

    function test_transferFromToUSDNFromNonZeroAccountReverts() public {
        // ACT: Transfer of 1M $USDN from Noble Core to USER1.
        bytes memory mintPayload = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000004e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000d8da6bf26964af9d7eed9e03e53415d37aa96045000000000000000000000000000000000000000000000000000000e8d4a51000"
        );
        (bool mintSuccess,) = MAILBOX.call(mintPayload);
        assertTrue(mintSuccess, "Initial mint should succeed");

        uint256 _user1Balance = usdn.balanceOf(USER1);

        assertEq(_user1Balance, 1000000e6, "user 1 should have 1 million usdn");

        vm.prank(USER1);
        usdn.approve(USER2, type(uint256).max);

        vm.expectRevert(abi.encodeWithSelector(NobleDollar.InvalidTransfer.selector));

        vm.prank(USER2);
        usdn.transferFrom(USER1, address(usdn), 1000e6);
    }

    function test_noClaimableYield() public {
        // Test when timestamp has not progressed from mint so claimable yield should revert
        bytes memory mintPayload = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000004e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000d8da6bf26964af9d7eed9e03e53415d37aa96045000000000000000000000000000000000000000000000000000000e8d4a51000"
        );
        (bool mintSuccess,) = MAILBOX.call(mintPayload);
        assertTrue(mintSuccess, "Initial mint should succeed");

        vm.expectRevert(NobleDollar.NoClaimableYield.selector);

        vm.prank(USER1);
        usdn.claim();

        // Test when account with zero balance calls claim()
        vm.expectRevert(NobleDollar.NoClaimableYield.selector);

        vm.prank(USER2);
        usdn.claim();
    }

    function test_uint112MaxMint() public {
        // Test when timestamp has not progressed from mint so claimable yield should revert
        bytes memory mintPayload = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000004e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000d8da6bf26964af9d7eed9e03e53415d37aa96045000000000000000000000000000000000000ffffffffffffffffffffffffffff"
        );
        (bool mintSuccess,) = MAILBOX.call(mintPayload);
        assertTrue(mintSuccess, "Initial mint should succeed");

        uint256 _user1Balance = usdn.balanceOf(USER1);

        assertEq(_user1Balance, type(uint112).max, "user 1 should have max usdn");

        // Test when timestamp has not progressed from mint so claimable yield should revert
        mintPayload = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000004e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000d8da6bf26964af9d7eed9e03e53415d37aa960450000000000000000000000000000000000000000000000000000000000000001"
        );
        (mintSuccess,) = MAILBOX.call(mintPayload);

        assertEq(mintSuccess, false, "minting more than uint112 max should fail");
    }

    function test_secondDepositPostYieldReceivesCorrectIndex() public {
        // Mint 1 million to USER1
        bytes memory mintPayload = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000004e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000d8da6bf26964af9d7eed9e03e53415d37aa96045000000000000000000000000000000000000000000000000000000e8d4a51000"
        );
        (bool mintSuccess,) = MAILBOX.call(mintPayload);
        assertTrue(mintSuccess, "Initial mint should succeed");

        // Accrue 1 million in yield
        bytes memory yieldPayload = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000014e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000000000000000000000000000000000e8d4a51000"
        );
        (bool yieldSuccess,) = MAILBOX.call(yieldPayload);
        assertTrue(yieldSuccess, "Yield accrual should succeed");

        // Mint 1 million to USER2
        mintPayload = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000004e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000f2f1acbe0ba726fee8d75f3e32900526874740bb000000000000000000000000000000000000000000000000000000e8d4a51000"
        );
        (mintSuccess,) = MAILBOX.call(mintPayload);
        assertEq(mintSuccess, true);

        uint256 _principalUSER1 = usdn.principalOf(USER1);
        uint256 _principalUSER2 = usdn.principalOf(USER2);

        assertEq(_principalUSER1, 1000000e6, "user 1 should have 1 million principal");
        assertEq(_principalUSER2, 500000e6, "user 2 should have 500 thousand principal");
    }

    function test_claimYield() public {
        // Mint 1M USDN to USER1
        bytes memory mintPayload = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000004e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000d8da6bf26964af9d7eed9e03e53415d37aa96045000000000000000000000000000000000000000000000000000000e8d4a51000"
        );
        (bool mintSuccess,) = MAILBOX.call(mintPayload);
        assertTrue(mintSuccess, "Initial mint should succeed");

        // Verify initial state
        assertEq(usdn.balanceOf(USER1), 1e12, "USER1 should have 1M USDN");
        assertEq(usdn.yield(USER1), 0, "USER1 should have no yield initially");

        // Accrue yield equal to deposit (1M USDN - 100% yield)
        bytes memory yieldPayload = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000014e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000000000000000000000000000000000e8d4a51000"
        );
        (bool yieldSuccess,) = MAILBOX.call(yieldPayload);
        assertTrue(yieldSuccess, "Yield accrual should succeed");

        // Verify yield is available
        uint256 claimableYield = usdn.yield(USER1);
        assertEq(claimableYield, 1e12, "USER1 should have 1M USDN in claimable yield");

        // Claim yield
        vm.expectEmit(true, true, true, true);
        emit NobleDollar.YieldClaimed(USER1, 1e12);

        vm.prank(USER1);
        usdn.claim();

        // Verify post-claim state
        assertEq(usdn.balanceOf(USER1), 2e12, "USER1 balance should be doubled (original + yield)");
        assertEq(usdn.balanceOf(address(usdn)), 0, "Contract balance should be zero after full claim");
        assertEq(usdn.yield(USER1), 0, "USER1 should have no claimable yield after claiming");
        assertEq(usdn.principalOf(USER1), 1e12, "USER1 principal should remain unchanged");
        assertEq(usdn.totalSupply(), 2e12, "Total supply should reflect the claimed yield");
    }

    function test_claimYieldMultipleUsers() public {
        // Mint 1M to USER1 and 2M to USER2
        bytes memory mintPayload1 = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000004e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000d8da6bf26964af9d7eed9e03e53415d37aa96045000000000000000000000000000000000000000000000000000000e8d4a51000"
        );
        (bool mintSuccess1,) = MAILBOX.call(mintPayload1);
        assertTrue(mintSuccess1);

        bytes memory mintPayload2 = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000004e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000f2f1acbe0ba726fee8d75f3e32900526874740bb000000000000000000000000000000000000000000000000000001d1a94a2000"
        );
        (bool mintSuccess2,) = MAILBOX.call(mintPayload2);
        assertTrue(mintSuccess2);

        // Accrue yield equal to total deposits (3M USDN total yield)
        bytes memory yieldPayload = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000014e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000000000000000000000000000000002ba7def3000"
        );
        (bool yieldSuccess,) = MAILBOX.call(yieldPayload);
        assertTrue(yieldSuccess);

        // Both users should have yield equal to their deposits (100% yield)
        uint256 user1Yield = usdn.yield(USER1);
        uint256 user2Yield = usdn.yield(USER2);

        assertEq(user1Yield, 1e12, "USER1 should have 1M yield");
        assertEq(user2Yield, 2e12, "USER2 should have 2M yield");

        // USER1 claims
        vm.prank(USER1);
        usdn.claim();
        assertEq(usdn.balanceOf(USER1), 2e12, "USER1 balance should be doubled");
        assertEq(usdn.yield(USER1), 0, "USER1 should have no yield after claiming");

        // USER2 claims
        vm.prank(USER2);
        usdn.claim();
        assertEq(usdn.balanceOf(USER2), 4e12, "USER2 balance should be doubled");
        assertEq(usdn.yield(USER2), 0, "USER2 should have no yield after claiming");
    }

    function test_claimYieldAfterTransfer() public {
        // Mint 1M to USER1
        bytes memory mintPayload = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000004e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000d8da6bf26964af9d7eed9e03e53415d37aa96045000000000000000000000000000000000000000000000000000000e8d4a51000"
        );
        (bool success,) = MAILBOX.call(mintPayload);

        // Accrue 1M yield (100% yield)
        bytes memory yieldPayload = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000014e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000000000000000000000000000000000e8d4a51000"
        );
        (success,) = MAILBOX.call(yieldPayload);

        // USER1 should have 1M yield
        assertEq(usdn.yield(USER1), 1e12, "USER1 should have 1M yield before transfer");

        // Transfer half (500k) to USER2
        vm.prank(USER1);
        usdn.transfer(USER2, 5e11);

        // Call yield
        uint256 user1Yield = usdn.yield(USER1);
        uint256 user2Yield = usdn.yield(USER2);

        assertEq(user1Yield, 1e12, "USER1 should have 1m of yield");
        assertEq(user2Yield, 0, "USER2 should have 0 yield");

        // USER1 claims their yield
        vm.prank(USER1);
        usdn.claim();
        assertEq(usdn.balanceOf(USER1), 15e11, "USER1 should have 1.5M balance after claiming yield");

        // USER2 yield claim should fail
        vm.expectRevert(abi.encode(NobleDollar.NoClaimableYield.selector));
        vm.prank(USER2);
        usdn.claim();

        // Accrue another 1M yield (now distributed proportionally)
        yieldPayload = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000024e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000000000000000000000000000000000e8d4a51000"
        );
        (success,) = MAILBOX.call(yieldPayload);

        // Both users should have yield proportional to their principal
        // USER1 has ~500k principal, USER2 has ~500k principal (from transfer)
        // So each should get approximately 500k yield
        uint256 user1NewYield = usdn.yield(USER1);
        uint256 user2NewYield = usdn.yield(USER2);

        // Due to rounding, yields might not be exactly 500k each
        assertApproxEqAbs(user1NewYield, 75e10, 2, "USER1 should have ~750k new yield");
        assertApproxEqAbs(user2NewYield, 25e10, 2, "USER1 should have ~250k new yield");

        // USER1 claims their yield
        vm.prank(USER1);
        usdn.claim();
        assertEq(usdn.balanceOf(USER1), 225e10, "USER1 should have 2.25M balance after claiming yield");

        // USER2 claims their yield
        vm.prank(USER2);
        usdn.claim();
        assertEq(usdn.balanceOf(USER2), 75e10, "USER2 should have 750k balance after claiming yield");
    }

    function test_indexUpdateWithZeroTotalPrincipal() public {
        // Edge case: yield accrual when totalPrincipal is 0 (no deposits)
        // The contract should handle this gracefully without reverting

        assertEq(usdn.totalSupply(), 0, "Total supply should be 0 initially");
        assertEq(usdn.totalPrincipal(), 0, "Total principal should be 0 initially");
        assertEq(usdn.index(), 1e12, "Index should be 1.0 initially");

        // Try to accrue yield when no principal exists
        bytes memory yieldPayload = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000014e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000000000000000000000000000000000e8d4a51000"
        );
        (bool success,) = MAILBOX.call(yieldPayload);
        assertTrue(success, "Yield accrual should succeed even with zero principal");

        // Index should remain unchanged when totalPrincipal is 0
        assertEq(usdn.index(), 1e12, "Index should remain 1.0 when no principal exists");
        assertEq(usdn.totalSupply(), 1e12, "Supply should increase by yield amount");
        assertEq(usdn.totalPrincipal(), 0, "Total principal should still be 0");
    }

    function test_indexUpdateAfterBurn() public {
        // Setup: Mint 1M to USER1
        bytes memory mintPayload = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000004e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000d8da6bf26964af9d7eed9e03e53415d37aa96045000000000000000000000000000000000000000000000000000000e8d4a51000"
        );
        (bool success,) = MAILBOX.call(mintPayload);
        assertTrue(success, "Initial mint should succeed");

        // Accrue 100% yield
        bytes memory yieldPayload1 = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000014e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000000000000000000000000000000000e8d4a51000"
        );
        (success,) = MAILBOX.call(yieldPayload1);
        assertTrue(success, "Yield accrual should succeed");
        assertEq(usdn.index(), 2e12, "Index should be 2.0 after 100% yield");

        // Burn entire pre yield claimed balance of USER1
        vm.prank(USER1);
        usdn.transferRemote{value: 1 ether}(
            1313817164, // destination domain (you used this in setUp)
            bytes32(uint256(uint160(USER2))),
            1e12 // amount to transfer (1M USDN)
        );

        uint256 _balanceUSER1 = usdn.balanceOf(USER1);
        uint256 _principalUSER1 = usdn.principalOf(USER1);

        assertEq(_balanceUSER1, 0, "USER1 balance should be 0 after transfer");
        assertEq(_principalUSER1, 5e11, "USER1 principal should be halved after burn");

        assertEq(usdn.totalPrincipal(), 5e11, "Total principal should be 500k after burn");
        assertEq(usdn.totalSupply(), 1e12, "Contract should still hold unclaimed yield");

        // Claim entire yield with remaining USER1 principal
        vm.prank(USER1);
        usdn.claim();

        _balanceUSER1 = usdn.balanceOf(USER1);
        _principalUSER1 = usdn.principalOf(USER1);

        assertEq(_balanceUSER1, 1e12, "USER1 should have all yield as balance after claiming");
        assertEq(_principalUSER1, 5e11, "USER1 principal should remain the same");

        // Realize 100% yield again, 1m tokens
        bytes memory yieldPayload2 = abi.encodeWithSignature(
            "process(bytes,bytes)",
            0x0,
            hex"03000000024e4f424c726f757465725f6170700000000000000000000000000001000000000000000000000001000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000f62849f9a0b5bf2913b396098f7c7019b51a820a000000000000000000000000000000000000000000000000000000e8d4a51000"
        );
        (success,) = MAILBOX.call(yieldPayload2);
        assertTrue(success, "Second yield accrual should succeed");

        _balanceUSER1 = usdn.balanceOf(USER1);
        _principalUSER1 = usdn.principalOf(USER1);

        assertEq(usdn.index(), 4e12, "index should be doubled again");
    }
}
