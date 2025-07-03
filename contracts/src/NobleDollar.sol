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

import {HypERC20} from "@hyperlane/token/HypERC20.sol";

/*

███╗   ██╗ ██████╗ ██████╗ ██╗     ███████╗      
████╗  ██║██╔═══██╗██╔══██╗██║     ██╔════╝      
██╔██╗ ██║██║   ██║██████╔╝██║     █████╗        
██║╚██╗██║██║   ██║██╔══██╗██║     ██╔══╝        
██║ ╚████║╚██████╔╝██████╔╝███████╗███████╗      
╚═╝  ╚═══╝ ╚═════╝ ╚═════╝ ╚══════╝╚══════╝      
                                                 
██████╗  ██████╗ ██╗     ██╗      █████╗ ██████╗ 
██╔══██╗██╔═══██╗██║     ██║     ██╔══██╗██╔══██╗
██║  ██║██║   ██║██║     ██║     ███████║██████╔╝
██║  ██║██║   ██║██║     ██║     ██╔══██║██╔══██╗
██████╔╝╚██████╔╝███████╗███████╗██║  ██║██║  ██║
╚═════╝  ╚═════╝ ╚══════╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝

*/

/**
 * @title  NobleDollar
 * @author NASD Inc.
 * @notice ERC20 Noble Dollar.
 */
contract NobleDollar is HypERC20 {
    /// TODO: Add NatSpec.
    event IndexUpdated(uint128 oldIndex, uint128 newIndex, uint256 totalPrincipal, uint256 yieldAccrued);

    /// @custom:storage-location erc7201:noble.storage.USDN
    struct USDNStorage {
        uint128 index;
        mapping(address account => uint256) principal;
        uint256 totalPrincipal;
    }

    // keccak256(abi.encode(uint256(keccak256("noble.storage.USDN")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 private constant USDNStorageLocation = 0xccec1a0a356b34ea3899fbc248aeaeba5687659563a3acddccc6f1e8a5d84200;

    function _getUSDNStorage() private pure returns (USDNStorage storage $) {
        assembly {
            $.slot := USDNStorageLocation
        }
    }

    constructor(address mailbox_) HypERC20(6, 1, mailbox_) {}

    function initialize(address hook_, address ism_) public virtual initializer {
        super.initialize("Noble Dollar", "USDN", hook_, ism_, msg.sender);

        USDNStorage storage $ = _getUSDNStorage();
        $.index = 1e12;
    }

    /**
     * @dev Returns the current index used for yield calculations.
     */
    function index() public view returns (uint256) {
        USDNStorage storage $ = _getUSDNStorage();
        return $.index;
    }

    /**
     * @dev Returns the amount of principal in existence.
     */
    function totalPrincipal() public view returns (uint256) {
        USDNStorage storage $ = _getUSDNStorage();
        return $.totalPrincipal;
    }

    /**
     * @dev Returns the amount of principal owned by `account`.
     */
    function principalOf(address account) public view returns (uint256) {
        USDNStorage storage $ = _getUSDNStorage();
        return $.principal[account];
    }

    /// @notice Returns the amount of yield claimable for a given account.
    function yield(address account) public view returns (uint256) {
        USDNStorage storage $ = _getUSDNStorage();
        uint256 expectedBalance = $.principal[account] * $.index / 1e12;
        uint256 currentBalance = balanceOf(account);

        return expectedBalance > currentBalance ? expectedBalance - currentBalance : 0;
    }

    /**
     * @dev Claims the yield for the caller.
     */
    function claim() public {
        uint256 amount = yield(msg.sender);
        if (amount == 0) {
            revert("No yield to claim");
        }

        _transfer(address(this), msg.sender, amount);
    }

    /// TODO: Add NatSpec.
    function _update(address from, address to, uint256 value) internal virtual override {
        USDNStorage storage $ = _getUSDNStorage();

        super._update(from, to, value);

        if (from == address(this)) {
            // We don't want to perform any principal updates in the case of yield payout.
            return;
        }
        if (to == address(this)) {
            if (from == address(0)) {
                // We don't want to perform any principal updates in the case of yield accrual.
                uint128 oldIndex = $.index;

                $.index = uint128(super.totalSupply() * 1e12 / $.totalPrincipal);

                emit IndexUpdated(oldIndex, $.index, $.totalPrincipal, value);

                return;
            }

            // TODO: Should we block a user transfer to the contract like on Noble?
        }

        uint256 principal = ((value * 1e12) + $.index - 1) / $.index;

        // We don't want to update the sender's principal in the case of issuance.
        if (from != address(0)) {
            $.principal[from] -= principal;
        } else {
            $.totalPrincipal += principal;
        }

        // We don't want to update the recipient's principal in the case of withdrawal.
        if (to != address(0)) {
            if (from == address(0)) {
                principal = (value * 1e12) / $.index;
            }

            $.principal[to] += principal;
        } else {
            $.totalPrincipal -= principal;
        }
    }
}
