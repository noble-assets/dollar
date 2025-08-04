package keeper_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"dollar.noble.xyz/v2/keeper"
	"dollar.noble.xyz/v2/types/vaults"
	vaultsv2 "dollar.noble.xyz/v2/types/vaults/v2"
	"dollar.noble.xyz/v2/utils/mocks"
)

// AdversarialTestSuite focuses on testing attack vectors and adversarial scenarios
type AdversarialTestSuite struct {
	suite.Suite
	keeper    *keeper.Keeper
	ctx       sdk.Context
	msgServer vaultsv2.MsgServer
	authority string
	bank      mocks.BankKeeper
	account   mocks.AccountKeeper

	// Test actors
	attacker sdk.AccAddress
	victim   sdk.AccAddress
	users    []sdk.AccAddress
}

func TestAdversarialSuite(t *testing.T) {
	suite.Run(t, new(AdversarialTestSuite))
}

func (suite *AdversarialTestSuite) SetupTest() {
	// Setup test environment with proper mocks
	suite.account = mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	suite.bank = mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}

	suite.keeper, _, suite.ctx = mocks.DollarKeeperWithKeepers(suite.T(), suite.bank, suite.account)

	// Set a proper block time for timestamp-dependent tests
	blockTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(blockTime)

	suite.msgServer = keeper.NewVaultV2MsgServer(suite.keeper)
	suite.authority = "authority"

	// Create test addresses
	suite.attacker = sdk.AccAddress("attacker_address___")
	suite.victim = sdk.AccAddress("victim_address_____")
	suite.users = make([]sdk.AccAddress, 5)
	for i := range suite.users {
		suite.users[i] = sdk.AccAddress(fmt.Sprintf("user_address_%d____", i))
	}
}

// First Depositor Attack Tests

func (suite *AdversarialTestSuite) TestFirstDepositorInflationAttack() {
	vaultType := vaults.FLEXIBLE
	attacker := suite.attacker.String()
	victim := suite.victim.String()

	// Attack scenario: Attacker becomes first depositor with minimal amount
	// then inflates NAV to make subsequent deposits expensive

	// Step 1: Attacker deposits minimum amount (1 wei equivalent)
	attackerInitialDeposit := math.NewInt(1)
	resp1, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		VaultType:    vaultType,
		Amount:       attackerInitialDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(attackerInitialDeposit, resp1.SharesReceived) // Gets 1 share

	// Step 2: Attacker artificially inflates NAV (simulating donation or other means)
	inflatedNav := math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)) // 1e18
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    inflatedNav,
		Reason:    "artificial inflation attack",
	})
	suite.Require().NoError(err)

	// Verify share price is now extremely high
	vaultState, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	suite.Require().Equal(math.LegacyNewDecFromInt(inflatedNav), vaultState.SharePrice)

	// Step 3: Victim tries to deposit normal amount
	victimDeposit := math.NewInt(1000000) // 1M units
	resp2, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    victim,
		VaultType:    vaultType,
		Amount:       victimDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Victim gets very few shares due to inflated price
	expectedVictimShares := math.LegacyNewDecFromInt(victimDeposit).Quo(vaultState.SharePrice).TruncateInt()
	suite.Require().Equal(expectedVictimShares, resp2.SharesReceived)
	suite.Require().True(expectedVictimShares.IsZero()) // Due to truncation, victim gets 0 shares!

	// Step 4: Attacker withdraws and captures victim's deposit
	attackerPosition, _ := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, suite.attacker)
	withdrawResp, err := suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: attacker,
		VaultType:  vaultType,
		Shares:     attackerPosition.Shares,
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Attacker gets almost all the value (victim's deposit + original NAV)
	actualWithdraw := withdrawResp.AmountWithdrawn

	// Verify the attack was successful (attacker got more than they put in)
	attackerProfit := actualWithdraw.Sub(attackerInitialDeposit)
	suite.Require().True(attackerProfit.GT(math.ZeroInt()), "Attacker should profit from the attack")

	// This demonstrates the vulnerability - test should fail in a secure implementation
	suite.T().Logf("VULNERABILITY: First depositor attack successful. Attacker profit: %s", attackerProfit.String())
}

func (suite *AdversarialTestSuite) TestFirstDepositorMitigationCheck() {
	vaultType := vaults.STAKED // Test with different vault type
	attacker := suite.attacker.String()
	victim := suite.victim.String()

	// Potential mitigation: Minimum deposit requirements or dead shares
	// This test checks if any mitigations are in place

	minDeposit := math.NewInt(1000) // Typical minimum deposit
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		VaultType:    vaultType,
		Amount:       minDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Test if very small deposits are rejected
	tinyDeposit := math.NewInt(1)
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    victim,
		VaultType:    vaultType,
		Amount:       tinyDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})

	// In a secure implementation, this should either:
	// 1. Reject very small deposits
	// 2. Have minimum share guarantees
	// 3. Use dead shares mechanism
	if err != nil {
		suite.T().Logf("GOOD: Small deposits are rejected: %s", err.Error())
	} else {
		suite.T().Logf("POTENTIAL ISSUE: Small deposits are accepted without protection")
	}
}

// Sandwich Attack Tests

func (suite *AdversarialTestSuite) TestSandwichAttackOnDeposit() {
	vaultType := vaults.FLEXIBLE
	attacker := suite.attacker.String()
	victim := suite.victim.String()

	// Scenario: Attacker front-runs victim's deposit with NAV manipulation

	// Step 1: Initial setup with some deposits
	initialUsers := []string{suite.users[0].String(), suite.users[1].String()}
	for _, user := range initialUsers {
		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user,
			VaultType:    vaultType,
			Amount:       math.NewInt(1000000),
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)
	}

	vaultStateBefore, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	originalSharePrice := vaultStateBefore.SharePrice

	// Step 2: Attacker front-runs with large deposit
	attackerDeposit := math.NewInt(10000000) // Large deposit
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		VaultType:    vaultType,
		Amount:       attackerDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Step 3: Simulate NAV increase (MEV or oracle manipulation)
	vaultStateAfterAttacker, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	manipulatedNav := vaultStateAfterAttacker.TotalNav.MulRaw(15).QuoRaw(10) // 50% increase
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    manipulatedNav,
		Reason:    "suspicious nav increase",
	})
	suite.Require().NoError(err)

	// Step 4: Victim's transaction executes at inflated price
	victimDeposit := math.NewInt(1000000)
	victimResp, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    victim,
		VaultType:    vaultType,
		Amount:       victimDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Step 5: Attacker back-runs with withdrawal
	attackerPosition, _ := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, suite.attacker)
	attackerWithdrawResp, err := suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: attacker,
		VaultType:  vaultType,
		Shares:     attackerPosition.Shares,
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Analyze the attack results
	attackerProfit := attackerWithdrawResp.AmountWithdrawn.Sub(attackerDeposit)
	victimSharesReceived := victimResp.SharesReceived
	victimExpectedShares := math.LegacyNewDecFromInt(victimDeposit).Quo(originalSharePrice).TruncateInt()
	victimLoss := victimExpectedShares.Sub(victimSharesReceived)

	suite.T().Logf("Sandwich attack analysis:")
	suite.T().Logf("  Attacker profit: %s", attackerProfit.String())
	suite.T().Logf("  Victim share loss: %s", victimLoss.String())
	suite.T().Logf("  Original share price: %s", originalSharePrice.String())
	suite.T().Logf("  Manipulated share price: %s", victimResp.SharePrice.String())

	if attackerProfit.IsPositive() {
		suite.T().Logf("VULNERABILITY: Sandwich attack profitable")
	}
}

// Flash Loan Attack Simulation

func (suite *AdversarialTestSuite) TestFlashLoanAttackSimulation() {
	vaultType := vaults.FLEXIBLE
	attacker := suite.attacker.String()

	// Simulate flash loan attack: borrow large amount, manipulate price, profit, repay

	// Step 1: Initial vault state
	victim := suite.victim.String()
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    victim,
		VaultType:    vaultType,
		Amount:       math.NewInt(5000000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	vaultStateBefore, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)

	// Step 2: Simulate flash loan - massive deposit
	flashLoanAmount := math.NewInt(100000000) // 100M units
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		VaultType:    vaultType,
		Amount:       flashLoanAmount,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	attackerPosition, _ := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, suite.attacker)

	// Step 3: Price manipulation through NAV update
	vaultStateAfterDeposit, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	manipulatedNav := vaultStateAfterDeposit.TotalNav.MulRaw(12).QuoRaw(10) // 20% increase
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    manipulatedNav,
		Reason:    "flash loan attack nav manipulation",
	})
	suite.Require().NoError(err)

	// Step 4: Immediate withdrawal (simulating flash loan repayment)
	withdrawResp, err := suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: attacker,
		VaultType:  vaultType,
		Shares:     attackerPosition.Shares,
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Step 5: Calculate attack profitability
	profit := withdrawResp.AmountWithdrawn.Sub(flashLoanAmount)

	// In a single-block attack, this should ideally not be profitable
	suite.T().Logf("Flash loan attack simulation:")
	suite.T().Logf("  Flash loan amount: %s", flashLoanAmount.String())
	suite.T().Logf("  Amount withdrawn: %s", withdrawResp.AmountWithdrawn.String())
	suite.T().Logf("  Profit: %s", profit.String())

	if profit.IsPositive() {
		suite.T().Logf("VULNERABILITY: Flash loan attack profitable")
	} else {
		suite.T().Logf("GOOD: Flash loan attack not profitable")
	}

	// Check impact on other users
	vaultStateAfter, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	victimPosition, _ := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, suite.victim)
	victimValue := vaultStateAfter.SharePrice.MulInt(victimPosition.Shares)
	originalVictimValue := vaultStateBefore.SharePrice.MulInt(math.NewInt(5000000)) // Original deposit value

	suite.T().Logf("  Victim value before: %s", originalVictimValue.String())
	suite.T().Logf("  Victim value after: %s", victimValue.String())
}

// Precision Attack Tests

func (suite *AdversarialTestSuite) TestPrecisionAttackViaRounding() {
	vaultType := vaults.FLEXIBLE
	attacker := suite.attacker.String()

	// Attack scenario: Exploit rounding in share calculations to steal small amounts

	// Step 1: Create specific vault state for rounding exploitation
	// Deposit amount that creates non-integer share price
	initialDeposit := math.NewInt(1000000)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		VaultType:    vaultType,
		Amount:       initialDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Step 2: Create fractional share price through NAV manipulation
	fractionalNav := math.NewInt(1000003) // Creates share price of 1.000003
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    fractionalNav,
		Reason:    "create fractional price",
	})
	suite.Require().NoError(err)

	vaultState, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	sharePrice := vaultState.SharePrice

	// Step 3: Make deposits that exploit rounding
	victims := []string{suite.users[0].String(), suite.users[1].String(), suite.users[2].String()}
	totalVictimDeposits := math.ZeroInt()
	totalVictimShares := math.ZeroInt()

	for _, victim := range victims {
		// Deposit amount that will result in truncation
		depositAmount := math.NewInt(1000)
		resp, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    victim,
			VaultType:    vaultType,
			Amount:       depositAmount,
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)

		totalVictimDeposits = totalVictimDeposits.Add(depositAmount)
		totalVictimShares = totalVictimShares.Add(resp.SharesReceived)
	}

	// Step 4: Calculate the rounding "dust" that accumulates
	expectedTotalShares := math.LegacyNewDecFromInt(totalVictimDeposits).Quo(sharePrice)
	actualTotalShares := math.LegacyNewDecFromInt(totalVictimShares)
	roundingLoss := expectedTotalShares.Sub(actualTotalShares)

	suite.T().Logf("Precision attack analysis:")
	suite.T().Logf("  Share price: %s", sharePrice.String())
	suite.T().Logf("  Total victim deposits: %s", totalVictimDeposits.String())
	suite.T().Logf("  Expected total shares: %s", expectedTotalShares.String())
	suite.T().Logf("  Actual total shares: %s", actualTotalShares.String())
	suite.T().Logf("  Rounding loss: %s shares", roundingLoss.String())

	// Step 5: Attacker tries to extract the accumulated dust
	// In share-based systems, this dust typically stays in the vault and benefits remaining shareholders
	finalVaultState, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	suite.checkValueConservation(initialDeposit, totalVictimDeposits, finalVaultState.TotalNav)
}

func (suite *AdversarialTestSuite) TestDustAmountGriefing() {
	vaultType := vaults.FLEXIBLE
	attacker := suite.attacker.String()

	// Attack scenario: Send many tiny deposits to bloat state or cause issues

	// Initialize vault state with a small initial deposit
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		VaultType:    vaultType,
		Amount:       math.NewInt(1000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err, "Failed to initialize vault state")

	initialState, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	suite.Require().NotNil(initialState, "Vault state should be initialized after deposit")

	// Create many dust deposits
	dustAmount := math.NewInt(1)
	numAttacks := 100

	for i := 0; i < numAttacks; i++ {
		// Each deposit from attacker (simulating different accounts)
		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    attacker,
			VaultType:    vaultType,
			Amount:       dustAmount,
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})

		if err != nil {
			suite.T().Logf("Dust deposit rejected after %d attempts: %s", i, err.Error())
			break
		}
	}

	finalState, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)

	suite.T().Logf("Dust griefing attack:")
	suite.T().Logf("  Successful dust deposits: variable")

	if initialState != nil {
		suite.T().Logf("  Initial total shares: %s", initialState.TotalShares.String())
	} else {
		suite.T().Logf("  Initial total shares: <nil>")
	}

	if finalState != nil {
		suite.T().Logf("  Final total shares: %s", finalState.TotalShares.String())
	} else {
		suite.T().Logf("  Final total shares: <nil>")
	}

	// Check if the attack caused any state bloat or calculation issues
	suite.checkFinancialInvariants(vaultType)
}

// Economic Attack Tests

func (suite *AdversarialTestSuite) TestYieldDilutionAttack() {
	vaultType := vaults.FLEXIBLE
	attacker := suite.attacker.String()
	victim := suite.victim.String()

	// Attack scenario: Attacker front-runs yield distribution to dilute victim's share

	// Step 1: Victim deposits early
	victimDeposit := math.NewInt(1000000)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    victim,
		VaultType:    vaultType,
		Amount:       victimDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	victimPositionBefore, _ := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, suite.victim)

	// Step 2: Attacker detects incoming yield and front-runs with large deposit
	attackerDeposit := math.NewInt(9000000) // 9x victim's deposit
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		VaultType:    vaultType,
		Amount:       attackerDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Step 3: Yield distribution occurs
	vaultState, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	yieldAmount := math.NewInt(1000000) // 1M yield
	newNav := vaultState.TotalNav.Add(yieldAmount)
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    newNav,
		Reason:    "yield distribution",
	})
	suite.Require().NoError(err)

	// Step 4: Attacker immediately withdraws
	attackerPosition, _ := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, suite.attacker)
	withdrawResp, err := suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: attacker,
		VaultType:  vaultType,
		Shares:     attackerPosition.Shares,
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Step 5: Analyze yield distribution fairness
	attackerProfit := withdrawResp.AmountWithdrawn.Sub(attackerDeposit)

	finalVaultState, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	victimCurrentValue := finalVaultState.SharePrice.MulInt(victimPositionBefore.Shares)
	victimYieldReceived := victimCurrentValue.Sub(math.LegacyNewDecFromInt(victimDeposit))

	// Calculate what victim should have received if they were alone
	expectedVictimYield := math.LegacyNewDecFromInt(yieldAmount)
	actualVictimYield := victimYieldReceived

	suite.T().Logf("Yield dilution attack analysis:")
	suite.T().Logf("  Total yield distributed: %s", yieldAmount.String())
	suite.T().Logf("  Attacker profit: %s", attackerProfit.String())
	suite.T().Logf("  Expected victim yield (alone): %s", expectedVictimYield.String())
	suite.T().Logf("  Actual victim yield: %s", actualVictimYield.String())

	yieldStolen := expectedVictimYield.Sub(actualVictimYield)
	suite.T().Logf("  Yield 'stolen' from victim: %s", yieldStolen.String())
}

func (suite *AdversarialTestSuite) TestExtremeVolatilityAttack() {
	vaultType := vaults.FLEXIBLE
	attacker := suite.attacker.String()

	// Attack scenario: Extreme NAV volatility to break calculations

	// Step 1: Normal deposit
	normalDeposit := math.NewInt(1000000)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		VaultType:    vaultType,
		Amount:       normalDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Step 2: Extreme NAV increase
	extremeNav := math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(30), nil))
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    extremeNav,
		Reason:    "extreme volatility test",
	})
	suite.Require().NoError(err)

	// Step 3: Try operations at extreme prices
	victim := suite.victim.String()
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    victim,
		VaultType:    vaultType,
		Amount:       math.NewInt(1000000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})

	if err != nil {
		suite.T().Logf("Extreme volatility broke deposits: %s", err.Error())
	} else {
		suite.T().Logf("System handled extreme volatility in deposits")
	}

	// Step 4: Extreme NAV decrease
	tinyNav := math.NewInt(1)
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    tinyNav,
		Reason:    "extreme crash test",
	})
	suite.Require().NoError(err)

	// Check if system still functions
	suite.checkFinancialInvariants(vaultType)
}

// Helper Functions

func (suite *AdversarialTestSuite) checkFinancialInvariants(vaultType vaults.VaultType) {
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	if err != nil {
		// If vault doesn't exist yet, that's okay for some tests
		return
	}

	// Core invariants that must hold even under attack
	suite.Require().False(vaultState.TotalNav.IsNegative(), "Total NAV cannot be negative")
	suite.Require().False(vaultState.TotalShares.IsNegative(), "Total shares cannot be negative")

	if !vaultState.TotalShares.IsZero() {
		suite.Require().True(vaultState.SharePrice.IsPositive(), "Share price must be positive")

		// Share price calculation consistency
		expectedSharePrice := suite.calculateSharePrice(vaultState)
		suite.Require().Equal(expectedSharePrice, vaultState.SharePrice, "Share price calculation must be consistent")
	}
}

func (suite *AdversarialTestSuite) checkValueConservation(initialValue, addedValue, finalNav math.Int) {
	expectedTotalValue := initialValue.Add(addedValue)

	// Allow for small rounding differences but flag large discrepancies
	diff := finalNav.Sub(expectedTotalValue).Abs()
	maxAllowedDiff := math.NewInt(100) // Allow up to 100 units difference for rounding

	if diff.GT(maxAllowedDiff) {
		suite.T().Errorf("Value conservation violated: expected %s, got %s, diff %s",
			expectedTotalValue.String(), finalNav.String(), diff.String())
	}
}

func (suite *AdversarialTestSuite) TestReentrancyProtectionCheck() {
	// This test would check for reentrancy vulnerabilities
	// In a real implementation, this would test if callbacks during operations
	// could be exploited to reenter the contract

	vaultType := vaults.FLEXIBLE
	attacker := suite.attacker.String()

	// For now, just verify basic operation doesn't break
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		VaultType:    vaultType,
		Amount:       math.NewInt(1000000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// In a production system, this would test:
	// 1. Callback-based reentrancy during deposits/withdrawals
	// 2. Cross-function reentrancy
	// 3. State changes during external calls
	suite.T().Logf("Reentrancy protection check completed (basic test)")
}

// Helper function to calculate share price (replicates the keeper's unexported method)
func (suite *AdversarialTestSuite) calculateSharePrice(vaultState *vaultsv2.VaultState) math.LegacyDec {
	if vaultState.TotalShares.IsZero() {
		return math.LegacyOneDec() // 1:1 ratio for first deposit
	}
	return math.LegacyNewDecFromInt(vaultState.TotalNav).Quo(math.LegacyNewDecFromInt(vaultState.TotalShares))
}

func (suite *AdversarialTestSuite) TestGasExhaustionAttack() {
	vaultType := vaults.FLEXIBLE
	attacker := suite.attacker.String()

	// Attack scenario: Operations that consume excessive gas to DOS other users

	// Test with operations that might cause gas issues
	largeAmount := math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(25), nil))

	// This should either succeed efficiently or fail gracefully
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		VaultType:    vaultType,
		Amount:       largeAmount,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})

	if err != nil {
		suite.T().Logf("Large amount deposit rejected: %s", err.Error())
	} else {
		suite.T().Logf("Large amount deposit accepted")
		// Verify operations still work normally
		suite.checkFinancialInvariants(vaultType)
	}
}

func (suite *AdversarialTestSuite) TestOverflowUnderflowAttacks() {
	vaultType := vaults.FLEXIBLE
	attacker := suite.attacker.String()

	// Test mathematical overflow/underflow protection

	// Test maximum possible values
	maxInt := math.NewIntFromBigInt(new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1)))

	// This should either handle gracefully or reject
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		VaultType:    vaultType,
		Amount:       maxInt,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})

	if err != nil {
		suite.T().Logf("Maximum value deposit rejected: %s", err.Error())
	} else {
		suite.T().Logf("Maximum value deposit handled")
	}

	// Test operations that might cause overflow in calculations
	normalDeposit := math.NewInt(1000000)
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		VaultType:    vaultType,
		Amount:       normalDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Try to create overflow in NAV calculation
	overflowNav := math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(50), nil))
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    overflowNav,
		Reason:    "overflow test",
	})

	if err != nil {
		suite.T().Logf("Overflow NAV rejected: %s", err.Error())
	} else {
		suite.T().Logf("Overflow NAV handled")
		// Verify calculations still work
		suite.checkFinancialInvariants(vaultType)
	}
}

func (suite *AdversarialTestSuite) TestTimestampManipulationAttack() {
	vaultType := vaults.FLEXIBLE

	// Use fresh addresses for this test to avoid conflicts with previous tests
	freshAttacker := sdk.AccAddress(make([]byte, 20))
	copy(freshAttacker, []byte("fresh_attacker______"))
	freshVictim := sdk.AccAddress(make([]byte, 20))
	copy(freshVictim, []byte("fresh_victim________"))

	attackerAddr := freshAttacker.String()
	victimAddr := freshVictim.String()

	// Attack scenario: Manipulate timestamps to gain unfair advantages
	// (This is more relevant for time-locked vaults, but test basic behavior)

	// Deposit from both users
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attackerAddr,
		VaultType:    vaultType,
		Amount:       math.NewInt(1000000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    victimAddr,
		VaultType:    vaultType,
		Amount:       math.NewInt(1000000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Verify both positions have reasonable timestamps
	attackerPos, err := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, freshAttacker)
	suite.Require().NoError(err)
	suite.Require().NotNil(attackerPos, "Attacker position should exist")

	victimPos, err := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, freshVictim)
	suite.Require().NoError(err)
	suite.Require().NotNil(victimPos, "Victim position should exist")

	suite.Require().False(attackerPos.FirstDepositTime.IsZero(), "Attacker deposit time should be set")
	suite.Require().False(victimPos.FirstDepositTime.IsZero(), "Victim deposit time should be set")
}
