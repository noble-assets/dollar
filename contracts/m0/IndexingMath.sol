// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.20 <0.9.0;

import { UIntMath } from "./UIntMath.sol";

/**
 * @title  Helper library for indexing math functions.
 * @author M^0 Labs
 */
library IndexingMath {
    /* ============ Variables ============ */

    /// @notice The scaling of indexes for exponent math.
    uint56 internal constant EXP_SCALED_ONE = 1e12;

    /* ============ Custom Errors ============ */

    /// @notice Emitted when a division by zero occurs.
    error DivisionByZero();

    /* ============ Exposed Functions ============ */

    /**
     * @notice Helper function to calculate `(x * EXP_SCALED_ONE) / y`, rounded down.
     * @dev    Inspired by USM (https://github.com/usmfum/USM/blob/master/contracts/WadMath.sol)
     */
    function divide240By128Down(uint240 x, uint128 y) internal pure returns (uint112) {
        if (y == 0) revert DivisionByZero();

        unchecked {
            // NOTE: While `uint256(x) * EXP_SCALED_ONE` can technically overflow, these divide/multiply functions are
            //       only used for the purpose of principal/present amount calculations for continuous indexing, and
            //       so for an `x` to be large enough to overflow this, it would have to be a possible result of
            //       `multiply112By128Down` or `multiply112By128Up`, which would already satisfy
            //       `uint256(x) * EXP_SCALED_ONE < type(uint240).max`.
            return UIntMath.safe112((uint256(x) * EXP_SCALED_ONE) / y);
        }
    }

    /**
     * @notice Helper function to calculate `(x * EXP_SCALED_ONE) / y`, rounded up.
     * @dev    Inspired by USM (https://github.com/usmfum/USM/blob/master/contracts/WadMath.sol)
     */
    function divide240By128Up(uint240 x, uint128 y) internal pure returns (uint112) {
        if (y == 0) revert DivisionByZero();

        unchecked {
            // NOTE: While `uint256(x) * EXP_SCALED_ONE` can technically overflow, these divide/multiply functions are
            //       only used for the purpose of principal/present amount calculations for continuous indexing, and
            //       so for an `x` to be large enough to overflow this, it would have to be a possible result of
            //       `multiply112By128Down` or `multiply112By128Up`, which would already satisfy
            //       `uint256(x) * EXP_SCALED_ONE < type(uint240).max`.
            return UIntMath.safe112(((uint256(x) * EXP_SCALED_ONE) + y - 1) / y);
        }
    }

    /**
     * @notice Helper function to calculate `(x * y) / EXP_SCALED_ONE`, rounded down.
     * @dev    Inspired by USM (https://github.com/usmfum/USM/blob/master/contracts/WadMath.sol)
     */
    function multiply112By128Down(uint112 x, uint128 y) internal pure returns (uint240) {
        unchecked {
            return uint240((uint256(x) * y) / EXP_SCALED_ONE);
        }
    }

    /**
     * @notice Helper function to calculate `(x * index) / EXP_SCALED_ONE`, rounded up.
     * @dev    Inspired by USM (https://github.com/usmfum/USM/blob/master/contracts/WadMath.sol)
     */
    function multiply112By128Up(uint112 x, uint128 index) internal pure returns (uint240 z) {
        unchecked {
            return uint240(((uint256(x) * index) + (EXP_SCALED_ONE - 1)) / EXP_SCALED_ONE);
        }
    }

    /**
     * @dev    Returns the present amount (rounded down) given the principal amount and an index.
     * @param  principalAmount The principal amount.
     * @param  index           An index.
     * @return The present amount rounded down.
     */
    function getPresentAmountRoundedDown(uint112 principalAmount, uint128 index) internal pure returns (uint240) {
        return multiply112By128Down(principalAmount, index);
    }

    /**
     * @dev    Returns the present amount (rounded up) given the principal amount and an index.
     * @param  principalAmount The principal amount.
     * @param  index           An index.
     * @return The present amount rounded up.
     */
    function getPresentAmountRoundedUp(uint112 principalAmount, uint128 index) internal pure returns (uint240) {
        return multiply112By128Up(principalAmount, index);
    }

    /**
     * @dev    Returns the principal amount given the present amount, using the current index.
     * @param  presentAmount The present amount.
     * @param  index         An index.
     * @return The principal amount rounded down.
     */
    function getPrincipalAmountRoundedDown(uint240 presentAmount, uint128 index) internal pure returns (uint112) {
        return divide240By128Down(presentAmount, index);
    }

    /**
     * @dev    Returns the principal amount given the present amount, using the current index.
     * @param  presentAmount The present amount.
     * @param  index         An index.
     * @return The principal amount rounded up.
     */
    function getPrincipalAmountRoundedUp(uint240 presentAmount, uint128 index) internal pure returns (uint112) {
        return divide240By128Up(presentAmount, index);
    }
}
