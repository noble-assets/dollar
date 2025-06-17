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
    uint256 public index = 1e12;

    constructor(address mailbox_) FungibleTokenRouter(1, mailbox_) {}

    function initialize(address hook_, address ism_) public virtual initializer {
        __ERC20_init("Noble Dollar", "USDN");
        _MailboxClient_initialize(hook_, ism_, msg.sender);
    }

    /// @inheritdoc ERC20Upgradeable
    function decimals() public pure override returns (uint8) {
        return 6;
    }

    /// @inheritdoc ERC20Upgradeable
    function totalSupply() public view override returns (uint256) {
        ERC20Storage storage $ = super._getERC20Storage();
        return $._totalSupply * index;
    }

    // TODO: Add NatSpec.
    function totalPrincipal() public view returns (uint256) {
        ERC20Storage storage $ = super._getERC20Storage();
        return $._totalSupply;
    }

    /// @inheritdoc ERC20Upgradeable
    function balanceOf(address account) public view override(ERC20Upgradeable, TokenRouter) returns (uint256) {
        ERC20Storage storage $ = super._getERC20Storage();
        return $._balances[account] * index;
    }

    /// TODO: Add NatSpec.
    function principalOf(address account) public view returns (uint256) {
        ERC20Storage storage $ = super._getERC20Storage();
        return $._balances[account];
    }

    /// @inheritdoc ERC20Upgradeable
    function _update(address from, address to, uint256 value) internal virtual override {
        uint256 principal = value / index;

        ERC20Storage storage $ = super._getERC20Storage();
        if (from == address(0)) {
            if (to == address(this)) {
                index = (totalSupply() + value) / totalSupply();

                // TODO(@john): Emit a rebasing event.

                return;
            } else {
                // Overflow check required: The rest of the code assumes that totalPrincipal never overflows
                $._totalSupply += principal;
            }
        } else {
            uint256 fromPrincipal = $._balances[from];
            if (fromPrincipal < principal) {
                revert ERC20InsufficientBalance(from, balanceOf(from), value);
            }
            unchecked {
                // Overflow not possible: principal <= fromPrincipal <= totalPrincipal.
                $._balances[from] = fromPrincipal - principal;
            }
        }

        if (to == address(0)) {
            unchecked {
                // Overflow not possible: principal <= totalPrincipal or principal <= fromPrincipal <= totalPrincipal.
                $._totalSupply -= principal;
            }
        } else {
            unchecked {
                // Overflow not possible: principal + toPrincipal is at most totalPrincipal, which we know fits into a uint256.
                $._balances[to] += principal;
            }
        }

        if (from != address(0) && to != address(this)) {
            emit Transfer(from, to, value);
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
