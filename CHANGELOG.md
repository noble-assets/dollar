# CHANGELOG

## v2.0.0

*Jun 6, 2025*

Second major release of the `x/dollar` module, enabling composable yield for the Noble Dollar ($USDN).

### BUG FIXES

- Correctly emit values inside `RewardClaimed` event. ([#33](https://github.com/noble-assets/dollar/pull/33))

### FEATURES

- Enable transfers and yield distribution across specific IBC channels. ([#28](https://github.com/noble-assets/dollar/pull/28))
- Enable transfers and yield distribution across specific Hyperlane routes. ([#31](https://github.com/noble-assets/dollar/pull/31))

### IMPROVEMENTS

- Update module path for v2 release line. ([#29](https://github.com/noble-assets/dollar/pull/29))
- Migrate custom ICS4 wrapper from main Noble codebase. ([#30](https://github.com/noble-assets/dollar/pull/30))
- Enforce lowercase user address when generating their flexible vault account. ([#36](https://github.com/noble-assets/dollar/pull/36))

## v1.0.2

*May 8, 2025*

This is a non-consensus breaking patch to the `v1` release line.

### FEATURES

- Include position details when querying a vault user's pending rewards. ([#39](https://github.com/noble-assets/dollar/pull/39))

## v1.0.1

*Mar 4, 2025*

This is a non-consensus breaking patch to the `v1` release line.

### BUG FIXES

- Correctly encode the recipient address in the `MTokenReceived` event. ([#27](https://github.com/noble-assets/dollar/pull/27))

## v1.0.0

*Feb 28, 2025*

Initial release of the `x/dollar` module, enabling the issuance of the Noble Dollar ($USDN).

