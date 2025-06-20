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

import {FungibleTokenRouter} from "@hyperlane/token/libs/FungibleTokenRouter.sol";
import {TokenRouter} from "@hyperlane/token/libs/TokenRouter.sol";
import {ERC20Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";

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
contract NobleDollar is ERC20Upgradeable, FungibleTokenRouter {
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

    constructor(address mailbox_) FungibleTokenRouter(1, mailbox_) {}

    function initialize(address hook_, address ism_) public virtual initializer {
        __ERC20_init("Noble Dollar", "USDN");
        _MailboxClient_initialize(hook_, ism_, msg.sender);

        USDNStorage storage $ = _getUSDNStorage();
        $.index = 1e12;
    }

    /// TODO: Add NatSpec.
    function index() public view returns (uint256) {
        USDNStorage storage $ = _getUSDNStorage();
        return $.index;
    }

    /// TODO: Add NatSpec.
    function totalPrincipal() public view returns (uint256) {
        USDNStorage storage $ = _getUSDNStorage();
        return $.totalPrincipal;
    }

    /// TODO: Add NatSpec.
    function principalOf(address account) public view returns (uint256) {
        USDNStorage storage $ = _getUSDNStorage();
        return $.principal[account];
    }

    /// @inheritdoc ERC20Upgradeable
    function decimals() public pure override returns (uint8) {
        return 6;
    }

    /// @inheritdoc ERC20Upgradeable
    function balanceOf(address account) public view override(ERC20Upgradeable, TokenRouter) returns (uint256) {
        return ERC20Upgradeable.balanceOf(account);
    }

    /// @inheritdoc ERC20Upgradeable
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
                $.index = uint128(super.totalSupply() / $.totalPrincipal);

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

    /// @inheritdoc TokenRouter
    function _transferFromSender(uint256 _amount) internal virtual override returns (bytes memory) {
        _burn(msg.sender, _amount);
        return bytes(""); // no metadata
    }

    /// @inheritdoc TokenRouter
    function _transferTo(
        address _recipient,
        uint256 _amount,
        bytes calldata // no metadata
    ) internal virtual override {
        _mint(_recipient, _amount);
    }
}
