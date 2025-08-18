// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.20 <0.9.0;

/**
 * @title  Library to perform safe math operations on uint types
 * @author M0 Labs
 */
library UIntMath {
    /* ============ Custom Errors ============ */

    /// @notice Emitted when a passed value is greater than the maximum value of uint112.
    error InvalidUInt112();

    /// @notice Emitted when a passed value is greater than the maximum value of uint128.
    error InvalidUInt128();

    /* ============ Internal View/Pure Functions ============ */

    /**
     * @notice Casts a uint256 value to a uint112, ensuring that it is less than or equal to the maximum uint112 value.
     * @param  n The value to cast.
     * @return The value casted to uint112.
     */
    function safe112(uint256 n) internal pure returns (uint112) {
        if (n > type(uint112).max) revert InvalidUInt112();
        return uint112(n);
    }

    /**
     * @notice Casts a uint256 value to a uint128, ensuring that it is less than or equal to the maximum uint128 value.
     * @param  n The value to cast.
     * @return The value casted to uint128.
     */
    function safe128(uint256 n) internal pure returns (uint128) {
        if (n > type(uint128).max) revert InvalidUInt128();
        return uint128(n);
    }
}
