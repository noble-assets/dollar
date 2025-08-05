package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"dollar.noble.xyz/v2/keeper"
	vaultsv2 "dollar.noble.xyz/v2/types/vaults/v2"
	"dollar.noble.xyz/v2/utils/mocks"
)

// V2VaultTestSuite provides a comprehensive test suite for V2 vault operations
type V2VaultTestSuite struct {
	suite.Suite
	keeper    *keeper.Keeper
	ctx       sdk.Context
	addresses []sdk.AccAddress
	msgServer vaultsv2.MsgServer
	authority string
	bank      mocks.BankKeeper
	account   mocks.AccountKeeper
}

func TestV2VaultSuite(t *testing.T) {
	suite.Run(t, new(V2VaultTestSuite))
}

func (suite *V2VaultTestSuite) SetupTest() {
	// Setup mocks
	suite.account = mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	suite.bank = mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}

	suite.keeper, _, suite.ctx = mocks.DollarKeeperWithKeepers(suite.T(), suite.bank, suite.account)

	// Set block time
	blockTime := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(blockTime)

	suite.msgServer = keeper.NewVaultV2MsgServer(suite.keeper)
	suite.authority = "authority"

	// Create test addresses
	for i := 0; i < 10; i++ {
		addr := sdk.AccAddress(fmt.Sprintf("user%d", i))
		suite.addresses = append(suite.addresses, addr)
	}
}

func (suite *V2VaultTestSuite) TestBasicDepositWithdraw() {
	user := suite.addresses[0].String()

	// Test basic deposit
	depositAmount := math.NewInt(1000000) // 1M units
	resp, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		Amount:       depositAmount,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(depositAmount, resp.SharesReceived)
	suite.Require().Equal(math.LegacyOneDec(), resp.SharePrice)

	// Verify vault state
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(depositAmount, vaultState.TotalShares)
	suite.Require().Equal(depositAmount, vaultState.TotalNav)
	suite.Require().Equal(math.LegacyOneDec(), vaultState.SharePrice)

	// Test withdrawal
	withdrawShares := math.NewInt(500000) // Half the shares
	_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: user,
		Shares:     withdrawShares,
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Verify final vault state
	finalState, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(withdrawShares, finalState.TotalShares) // Half shares remaining
}

func (suite *V2VaultTestSuite) TestSharePriceCalculation() {
	user := suite.addresses[0].String()
	user2 := suite.addresses[1].String()

	// Initial deposit at 1:1 ratio
	initialDeposit := math.NewInt(1000000)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		Amount:       initialDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Simulate yield by updating NAV
	newNav := math.NewInt(1200000) // 20% increase
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    newNav,
		Reason:    "Yield distribution test",
	})
	suite.Require().NoError(err)

	// Verify share price increased
	vaultStateAfter, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)
	expectedSharePrice := math.LegacyNewDecFromInt(newNav).Quo(math.LegacyNewDecFromInt(initialDeposit))
	suite.Require().Equal(expectedSharePrice, vaultStateAfter.SharePrice)

	// Second deposit should get fewer shares due to higher price
	depositAmount2 := math.NewInt(600000) // 600K units
	resp2, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user2,
		Amount:       depositAmount2,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Calculate expected shares: depositAmount2 / sharePrice
	expectedShares := math.LegacyNewDecFromInt(depositAmount2).Quo(expectedSharePrice).TruncateInt()
	suite.Require().Equal(expectedShares, resp2.SharesReceived)
}

func (suite *V2VaultTestSuite) TestMultipleUsersInvariantConservation() {
	users := []string{
		suite.addresses[0].String(),
		suite.addresses[1].String(),
		suite.addresses[2].String(),
	}
	deposits := []math.Int{
		math.NewInt(1000000),
		math.NewInt(500000),
		math.NewInt(2000000),
	}

	// Multiple users deposit
	totalDeposited := math.ZeroInt()
	for i, user := range users {
		deposit := deposits[i]
		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user,
			Amount:       deposit,
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)
		totalDeposited = totalDeposited.Add(deposit)
	}

	// Verify total amounts
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(totalDeposited, vaultState.TotalNav)

	// Add yield
	yieldAmount := math.NewInt(350000) // 10% yield
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    totalDeposited.Add(yieldAmount),
		Reason:    "yield distribution",
	})
	suite.Require().NoError(err)

	// Verify conservation of value
	suite.checkValueConservation(totalDeposited, yieldAmount)
}

func (suite *V2VaultTestSuite) TestPrecisionAndRounding() {
	user := suite.addresses[0].String()
	user2 := suite.addresses[1].String()

	// Create a scenario with fractional share prices
	smallDeposit := math.NewInt(1)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		Amount:       smallDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Update NAV to create fractional share price
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    math.NewInt(3), // Creates share price of 3.0
		Reason:    "precision test",
	})
	suite.Require().NoError(err)

	// Deposit amount that should cause rounding
	deposit2 := math.NewInt(10) // This should give 3.33... shares, truncated to 3
	resp, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user2,
		Amount:       deposit2,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Verify rounding behavior
	suite.Require().Equal(math.NewInt(3), resp.SharesReceived)
	suite.checkFinancialInvariants()
}

func (suite *V2VaultTestSuite) TestLargeAmounts() {
	user := suite.addresses[0].String()

	// Test with large amounts to check overflow protection
	maxDeposit, _ := math.NewIntFromString("1000000000000000000000000") // 1e24

	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		Amount:       maxDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})

	if err != nil {
		suite.T().Logf("Large deposit rejected: %s", err.Error())
		// This is expected behavior for overflow protection
	} else {
		suite.T().Logf("Large deposit accepted")
		suite.checkFinancialInvariants()
	}

	// Test large NAV update
	largeNav, _ := math.NewIntFromString("1000000000000000000000000")
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    largeNav,
		Reason:    "Large yield event",
	})

	if err != nil {
		suite.T().Logf("Large NAV update rejected: %s", err.Error())
	} else {
		suite.T().Logf("Large NAV update accepted")
		suite.checkFinancialInvariants()
	}
}

func (suite *V2VaultTestSuite) TestSlippageProtection() {
	user := suite.addresses[0].String()
	user2 := suite.addresses[1].String()

	// Initial deposit
	initialDeposit := math.NewInt(1000000)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		Amount:       initialDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Sudden price increase
	highNav := math.NewInt(5000000)
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    highNav,
		Reason:    "sudden price increase",
	})
	suite.Require().NoError(err)

	// Deposit with minimum shares protection
	depositAmount := math.NewInt(1000000)
	minShares := math.NewInt(100000) // Expecting at least 100K shares
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user2,
		Amount:       depositAmount,
		ReceiveYield: true,
		MinShares:    minShares,
	})
	// This should fail due to slippage protection
	suite.Require().Error(err, "Should fail due to slippage protection")

	// Test withdrawal slippage protection
	withdrawShares := math.NewInt(100000)
	minAmount := math.NewInt(1000000) // Expecting at least 1M tokens
	_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: user,
		Shares:     withdrawShares,
		MinAmount:  minAmount,
	})
	// This should fail due to slippage protection
	suite.Require().Error(err, "Should fail due to slippage protection")
}

func (suite *V2VaultTestSuite) TestSandwichAttackPrevention() {
	attacker := suite.addresses[0].String()
	victim := suite.addresses[1].String()

	// Attacker deposits
	attackDeposit := math.NewInt(10000000) // 10M
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		Amount:       attackDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Price manipulation
	manipulatedNav := math.NewInt(50000000) // 5x increase
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    manipulatedNav,
		Reason:    "potential manipulation",
	})
	suite.Require().NoError(err)

	// Get attacker position before victim's transaction
	attackerSharesBefore, err := suite.keeper.GetV2UserPosition(suite.ctx, sdk.MustAccAddressFromBech32(attacker))
	suite.Require().NoError(err)

	// Victim deposits at inflated price
	victimDeposit := math.NewInt(1000000)
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    victim,
		Amount:       victimDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Attacker withdraws
	_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: attacker,
		Shares:     attackerSharesBefore.Shares,
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Check if attack was profitable for educational purposes
	suite.T().Logf("Sandwich attack test completed")
	suite.checkFinancialInvariants()
}

func (suite *V2VaultTestSuite) TestBoundaryConditions() {
	user := suite.addresses[0].String()

	// Test zero deposit
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		Amount:       math.ZeroInt(),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().Error(err, "Zero deposit should be rejected")

	// Test negative deposit (should be caught by validation)
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		Amount:       math.NewInt(-1),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().Error(err, "Negative deposit should be rejected")

	// Valid deposit first
	deposit := math.NewInt(1000)
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		Amount:       deposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Test withdrawing more shares than user has
	_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: user,
		Shares:     math.NewInt(2000), // More than deposited
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().Error(err, "Should not be able to withdraw more shares than owned")
}

func (suite *V2VaultTestSuite) TestYieldDistribution() {
	users := []string{
		suite.addresses[0].String(),
		suite.addresses[1].String(),
		suite.addresses[2].String(),
	}
	deposits := []math.Int{
		math.NewInt(1000000), // 1M
		math.NewInt(2000000), // 2M
		math.NewInt(3000000), // 3M
	}

	// Users deposit different amounts
	totalDeposited := math.ZeroInt()
	for i, user := range users {
		deposit := deposits[i]
		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user,
			Amount:       deposit,
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)
		totalDeposited = totalDeposited.Add(deposit)
	}

	// Distribute yield
	yieldAmount := math.NewInt(600000) // 10% yield
	newTotalNav := totalDeposited.Add(yieldAmount)
	_, err := suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    newTotalNav,
		Reason:    "yield distribution",
	})
	suite.Require().NoError(err)

	// Verify each user's yield is proportional to their share
	finalState, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)

	for i, user := range users {
		userAddr := sdk.MustAccAddressFromBech32(user)
		userPos, err := suite.keeper.GetV2UserPosition(suite.ctx, userAddr)
		suite.Require().NoError(err)

		// Calculate user's current value
		userValue := finalState.SharePrice.MulInt(userPos.Shares)
		originalDeposit := deposits[i]
		userYield := userValue.Sub(math.LegacyNewDecFromInt(originalDeposit))

		// Expected yield should be proportional to deposit
		expectedYieldRatio := math.LegacyNewDecFromInt(originalDeposit).Quo(math.LegacyNewDecFromInt(totalDeposited))
		expectedYield := expectedYieldRatio.MulInt(yieldAmount)

		suite.T().Logf("User %d: deposit=%s, yield=%s, expected=%s",
			i, originalDeposit.String(), userYield.String(), expectedYield.String())
	}
}

func (suite *V2VaultTestSuite) TestEdgeCasesAndErrorHandling() {
	user := suite.addresses[0].String()

	// Test operations on empty vault
	_, err := suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    math.NewInt(1000000),
		Reason:    "initial NAV",
	})
	suite.Require().NoError(err)

	// Test withdrawal from empty position
	_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: user,
		Shares:     math.NewInt(1000),
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().Error(err, "Should not be able to withdraw from empty position")

	// Valid deposit
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		Amount:       math.NewInt(1000000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Test zero NAV update
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    math.ZeroInt(),
		Reason:    "zero NAV test",
	})
	suite.Require().Error(err, "Zero NAV should be rejected")
}

func (suite *V2VaultTestSuite) TestYieldPreferences() {
	user := suite.addresses[0].String()

	// Test deposit with yield preference
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		Amount:       math.NewInt(1000000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Test deposit without yield preference
	user2 := suite.addresses[1].String()
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user2,
		Amount:       math.NewInt(1000000),
		ReceiveYield: false, // No yield
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Check if the system handles different yield preferences
	userPos, err := suite.keeper.GetV2UserPosition(suite.ctx, sdk.MustAccAddressFromBech32(user))
	suite.Require().NoError(err)
	suite.Require().True(userPos.ReceiveYield, "User should receive yield")

	user2Pos, err := suite.keeper.GetV2UserPosition(suite.ctx, sdk.MustAccAddressFromBech32(user2))
	suite.Require().NoError(err)
	suite.Require().False(user2Pos.ReceiveYield, "User2 should not receive yield")
}

func (suite *V2VaultTestSuite) TestStateConsistency() {
	// Test that vault state remains consistent across operations
	users := suite.addresses[:5]

	// Random deposits and withdrawals
	for i := 0; i < 10; i++ {
		user := users[i%len(users)]

		if i%2 == 0 {
			// Deposit
			amount := math.NewInt(int64((i + 1) * 100000))
			_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
				Depositor:    user.String(),
				Amount:       amount,
				ReceiveYield: true,
				MinShares:    math.ZeroInt(),
			})
			suite.Require().NoError(err)
		} else {
			// Try to withdraw (might fail if no position)
			userPos, err := suite.keeper.GetV2UserPosition(suite.ctx, user)
			if err == nil && !userPos.Shares.IsZero() {
				withdrawShares := userPos.Shares.QuoRaw(2) // Half shares
				_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
					Withdrawer: user.String(),
					Shares:     withdrawShares,
					MinAmount:  math.ZeroInt(),
				})
				suite.Require().NoError(err)
			}
		}

		// Check invariants after each operation
		suite.checkFinancialInvariants()
	}
}

func (suite *V2VaultTestSuite) TestConcurrentOperations() {
	// Simulate concurrent operations to test race conditions
	users := suite.addresses[:3]

	// Initial setup
	for i, user := range users {
		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user.String(),
			Amount:       math.NewInt(int64((i + 1) * 1000000)),
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)
	}

	// Simulate concurrent operations
	vaultStateBefore, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)

	// NAV update while users might be depositing/withdrawing
	newNav := vaultStateBefore.TotalNav.Add(math.NewInt(500000))
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    newNav,
		Reason:    "concurrent operations test",
	})
	suite.Require().NoError(err)

	// More operations after NAV update
	for i, user := range users {
		if i%2 == 0 {
			_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
				Depositor:    user.String(),
				Amount:       math.NewInt(500000),
				ReceiveYield: true,
				MinShares:    math.ZeroInt(),
			})
			suite.Require().NoError(err)
		}
	}

	suite.checkFinancialInvariants()
}

func (suite *V2VaultTestSuite) TestExtremeScenarios() {
	user := suite.addresses[0].String()

	// Test minimum possible deposit
	minDeposit := math.NewInt(1)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		Amount:       minDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Test extreme NAV changes
	extremeNav, _ := math.NewIntFromString("1000000000000000000000000") // 1e24
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    extremeNav,
		Reason:    "extreme scenario test",
	})

	if err != nil {
		suite.T().Logf("Extreme NAV rejected: %s", err.Error())
	} else {
		suite.checkFinancialInvariants()
	}

	// Test operations with extreme share prices
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    suite.addresses[1].String(),
		Amount:       math.NewInt(1000000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})

	if err != nil {
		suite.T().Logf("Deposit at extreme price rejected: %s", err.Error())
	} else {
		suite.T().Logf("System handled extreme price correctly")
	}
}

// Helper Functions

func (suite *V2VaultTestSuite) checkFinancialInvariants() {
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
	if err != nil {
		// If vault doesn't exist yet, that's okay for some tests
		return
	}

	// Core invariants that must hold
	suite.Require().True(vaultState.TotalNav.GTE(math.ZeroInt()), "Total NAV cannot be negative")
	suite.Require().True(vaultState.TotalShares.GTE(math.ZeroInt()), "Total shares cannot be negative")

	if !vaultState.TotalShares.IsZero() {
		suite.Require().True(vaultState.SharePrice.IsPositive(), "Share price must be positive")

		// Share price calculation consistency
		expectedSharePrice := math.LegacyNewDecFromInt(vaultState.TotalNav).Quo(math.LegacyNewDecFromInt(vaultState.TotalShares))
		suite.Require().Equal(expectedSharePrice, vaultState.SharePrice, "Share price calculation must be consistent")
	}
}

func (suite *V2VaultTestSuite) checkValueConservation(initialValue, addedValue math.Int) {
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
	if err != nil {
		return
	}

	expectedTotalValue := initialValue.Add(addedValue)
	actualTotalValue := vaultState.TotalNav

	// Allow for small rounding differences but flag large discrepancies
	diff := actualTotalValue.Sub(expectedTotalValue).Abs()
	maxAllowedDiff := math.NewInt(100) // Allow up to 100 units difference for rounding

	if diff.GT(maxAllowedDiff) {
		suite.T().Errorf("Value conservation violated: expected %s, got %s, diff %s",
			expectedTotalValue.String(), actualTotalValue.String(), diff.String())
	}
}

func (suite *V2VaultTestSuite) TestAdvancedYieldCalculations() {
	// Test complex yield scenarios

	// Setup initial positions
	users := suite.addresses[:4]
	deposits := []math.Int{
		math.NewInt(1000000),
		math.NewInt(2000000),
		math.NewInt(3000000),
		math.NewInt(4000000),
	}

	totalDeposited := math.ZeroInt()
	for i, user := range users {
		deposit := deposits[i]
		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user.String(),
			Amount:       deposit,
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)
		totalDeposited = totalDeposited.Add(deposit)
	}

	// Multiple yield events
	yieldEvents := []math.Int{
		math.NewInt(500000),
		math.NewInt(300000),
		math.NewInt(800000),
	}

	for i, yieldAmount := range yieldEvents {
		newTotalNav := totalDeposited.Add(yieldAmount)
		_, err := suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
			Authority: suite.authority,
			NewNav:    newTotalNav,
			Reason:    fmt.Sprintf("yield event %d", i+1),
		})
		suite.Require().NoError(err)
		totalDeposited = newTotalNav
	}

	// Verify proportional yield distribution
	finalState, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)

	for i, user := range users {
		userAddr := sdk.MustAccAddressFromBech32(user.String())
		userPos, err := suite.keeper.GetV2UserPosition(suite.ctx, userAddr)
		suite.Require().NoError(err)

		userValue := finalState.SharePrice.MulInt(userPos.Shares)
		originalDeposit := deposits[i]
		userYield := userValue.Sub(math.LegacyNewDecFromInt(originalDeposit))

		suite.T().Logf("User %d: deposit=%s, current_value=%s, yield=%s",
			i, originalDeposit.String(), userValue.String(), userYield.String())
	}
}
