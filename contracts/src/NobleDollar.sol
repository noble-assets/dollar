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
    /**
     * @notice Constructs the Noble Dollar contract.
     * @param  mailbox_ The address of the default Hyperlane Mailbox.
     * @param  hook_    The address of the default Hyperlane Hook.
     * @param  ism_     The address of the default Hyperlane Interchain Security Module.
     */
    constructor(address mailbox_, address hook_, address ism_) HypERC20(6, 1, mailbox_) {
        super.initialize(0, "Noble Dollar", "USDN", hook_, ism_, msg.sender);
    }
}
