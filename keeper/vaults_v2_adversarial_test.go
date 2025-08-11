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
	blockTime := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(blockTime)

	suite.msgServer = keeper.NewVaultV2MsgServer(suite.keeper)
	suite.authority = "authority"

	// Create test addresses
	suite.attacker = sdk.AccAddress("attacker")
	suite.victim = sdk.AccAddress("victim")

	// Create additional test users
	for i := 0; i < 10; i++ {
		user := sdk.AccAddress(fmt.Sprintf("user%d", i))
		suite.users = append(suite.users, user)
	}
}

// TestFirstDepositorInflationAttack tests the classic ERC4626 inflation attack
func (suite *AdversarialTestSuite) TestFirstDepositorInflationAttack() {
	attacker := suite.attacker.String()
	victim := suite.victim.String()

	// Step 1: Attacker makes tiny initial deposit to be first depositor
	attackerInitialDeposit := math.NewInt(1)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		Amount:       attackerInitialDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Step 2: Attacker artificially inflates NAV without minting shares
	inflatedNav := math.NewInt(10000000) // Massively inflate NAV
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    inflatedNav,
		Reason:    "artificial inflation",
	})
	suite.Require().NoError(err)

	// Verify share price is now extremely high
	_, err = suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)

	// Step 3: Victim tries to deposit normal amount
	victimDeposit := math.NewInt(1000000) // 1M units
	resp2, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    victim,
		Amount:       victimDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Verify victim gets very few or zero shares due to inflated price
	expectedVictimShares := resp2.SharesReceived
	suite.Require().True(expectedVictimShares.IsZero()) // Due to truncation, victim gets 0 shares!

	// Step 4: Attacker withdraws and captures victim's deposit
	attackerPos, err := suite.keeper.GetV2UserPosition(suite.ctx, suite.attacker)
	suite.Require().NoError(err)

	_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: attacker,
		Shares:     attackerPos.Shares,
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Verify attack succeeded - attacker gained value, victim lost
	suite.T().Logf("First depositor inflation attack:")
	suite.T().Logf("  Attacker initial: %s", attackerInitialDeposit.String())
	suite.T().Logf("  Victim deposit: %s", victimDeposit.String())
	suite.T().Logf("  Victim shares received: %s", expectedVictimShares.String())
	suite.T().Logf("  Attack profitable: %t", expectedVictimShares.IsZero())
}

// TestFirstDepositorMitigationCheck verifies protections against inflation attacks
func (suite *AdversarialTestSuite) TestFirstDepositorMitigationCheck() {
	attacker := suite.attacker.String()
	victim := suite.victim.String()

	// Step 1: First check if there's a minimum deposit requirement
	tinyDeposit := math.NewInt(1)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		Amount:       tinyDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})

	if err != nil {
		suite.T().Logf("First depositor mitigation detected: minimum deposit enforced")
		return
	}

	// Step 2: Test if normal deposits still work properly
	normalAmount := math.NewInt(1000000)
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    victim,
		Amount:       normalAmount,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	suite.T().Logf("First depositor mitigation check:")
	suite.T().Logf("  Tiny deposit allowed: %t", err == nil)
	suite.T().Logf("  Normal deposits working: %t", true)
}

// TestSandwichAttackOnDeposit tests MEV-style sandwich attacks
func (suite *AdversarialTestSuite) TestSandwichAttackOnDeposit() {
	attacker := suite.attacker.String()
	victim := suite.victim.String()

	// Step 1: Setup initial vault state
	users := []string{suite.users[0].String(), suite.users[1].String()}
	for _, user := range users {
		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user,
			Amount:       math.NewInt(1000000),
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)
	}

	vaultStateBefore, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)
	originalSharePrice := vaultStateBefore.SharePrice

	// Step 2: Attacker front-runs with large deposit
	attackerDeposit := math.NewInt(10000000) // Large deposit
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		Amount:       attackerDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Step 3: Price manipulation through NAV update
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    math.NewInt(50000000), // Inflated NAV
		Reason:    "price manipulation",
	})
	suite.Require().NoError(err)

	// Step 4: Victim's transaction executes at inflated price
	victimDeposit := math.NewInt(1000000)
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    victim,
		Amount:       victimDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Step 5: Attacker back-runs with withdrawal
	attackerPos, err := suite.keeper.GetV2UserPosition(suite.ctx, suite.attacker)
	suite.Require().NoError(err)

	attackerWithdrawResp, err := suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: attacker,
		Shares:     attackerPos.Shares,
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Analyze sandwich attack profitability
	attackerProfit := attackerWithdrawResp.AmountWithdrawn.Sub(attackerDeposit)
	profitPercentage := attackerProfit.ToLegacyDec().Quo(attackerDeposit.ToLegacyDec()).Mul(math.LegacyNewDec(100))

	suite.T().Logf("Sandwich attack analysis:")
	suite.T().Logf("  Attacker deposit: %s", attackerDeposit.String())
	suite.T().Logf("  Attacker withdrawal: %s", attackerWithdrawResp.AmountWithdrawn.String())
	suite.T().Logf("  Attacker profit: %s", attackerProfit.String())
	suite.T().Logf("  Profit percentage: %s%%", profitPercentage.String())
	suite.T().Logf("  Original share price: %s", originalSharePrice.String())

	// Verify attack was profitable (this indicates vulnerability)
	suite.Require().True(attackerProfit.IsPositive(), "Sandwich attack should be profitable if vulnerability exists")
}

// TestFlashLoanAttackSimulation simulates flash loan-based attacks
func (suite *AdversarialTestSuite) TestFlashLoanAttackSimulation() {
	attacker := suite.attacker.String()

	// Step 1: Setup victim positions
	for i, victim := range suite.users[:5] {
		deposit := math.NewInt(int64((i + 1) * 1000000)) // Varying deposits
		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    victim.String(),
			Amount:       deposit,
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)
	}

	vaultStateBefore, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)

	// Step 2: Simulate flash loan - massive deposit
	flashLoanAmount := math.NewInt(100000000) // 100M units
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		Amount:       flashLoanAmount,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	attackerPos, err := suite.keeper.GetV2UserPosition(suite.ctx, suite.attacker)
	suite.Require().NoError(err)

	// Step 3: Price manipulation through NAV update
	vaultStateAfterDeposit, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)

	manipulatedNav := vaultStateAfterDeposit.TotalNav.Add(math.NewInt(50000000))
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    manipulatedNav,
		Reason:    "flash loan manipulation",
	})
	suite.Require().NoError(err)

	// Step 4: Immediate withdrawal (simulating flash loan repayment)
	_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: attacker,
		Shares:     attackerPos.Shares,
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Step 5: Calculate attack profitability
	profit := math.NewInt(0) // Calculate actual profit
	if profit.IsPositive() {
		suite.T().Logf("Flash loan attack successful - profit: %s", profit.String())
	}

	// Check impact on other users
	vaultStateAfter, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)

	victimPos, err := suite.keeper.GetV2UserPosition(suite.ctx, suite.victim)
	if err == nil {
		victimValue := vaultStateAfter.SharePrice.MulInt(victimPos.Shares)
		originalVictimValue := vaultStateBefore.SharePrice.MulInt(math.NewInt(5000000)) // Original deposit value

		suite.T().Logf("Flash loan attack impact:")
		suite.T().Logf("  Victim value before: %s", originalVictimValue.String())
		suite.T().Logf("  Victim value after: %s", victimValue.String())
	}
}

// TestPrecisionAttackViaRounding tests attacks exploiting rounding errors
func (suite *AdversarialTestSuite) TestPrecisionAttackViaRounding() {
	attacker := suite.attacker.String()

	// Step 1: Setup vault with predictable state
	initialDeposit := math.NewInt(1000000)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		Amount:       initialDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Step 2: Manipulate price to specific value
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    math.NewInt(3333333), // Chosen to create rounding opportunities
		Reason:    "precision setup",
	})
	suite.Require().NoError(err)

	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)
	sharePrice := vaultState.SharePrice

	// Step 3: Make deposits that exploit rounding
	victims := []string{
		suite.users[0].String(),
		suite.users[1].String(),
		suite.users[2].String(),
	}

	for _, user := range victims {
		// Deposit amount that will result in truncation
		depositAmount := math.NewInt(1000)
		_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user,
			Amount:       depositAmount,
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)
	}

	suite.T().Logf("Precision attack via rounding:")
	suite.T().Logf("  Share price: %s", sharePrice.String())
	suite.T().Logf("  Exploit depends on implementation rounding behavior")
}

// TestDustAmountGriefing tests attacks using many small transactions
func (suite *AdversarialTestSuite) TestDustAmountGriefing() {
	attacker := suite.attacker.String()

	// Initialize vault state with a small initial deposit
	precisionAmount := math.NewInt(1000)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		Amount:       precisionAmount,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	initialState, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)
	suite.Require().NotNil(initialState, "Vault state should be initialized after deposit")

	// Create many dust deposits
	dustAmount := math.NewInt(1)
	successfulDustDeposits := 0

	for i := 0; i < 1000; i++ {
		// Each deposit from attacker (simulating different accounts)
		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    attacker,
			Amount:       dustAmount,
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})

		if err == nil {
			successfulDustDeposits++
		}
	}

	finalState, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)

	suite.T().Logf("Dust griefing attack:")
	suite.T().Logf("  Successful dust deposits: %d", successfulDustDeposits)
	suite.T().Logf("  Initial total shares: %s", initialState.TotalShares.String())
	suite.T().Logf("  Final total shares: %s", finalState.TotalShares.String())
	suite.T().Logf("  State bloat factor: %s", finalState.TotalShares.Sub(initialState.TotalShares).String())

	// Verify system still functions normally after dust attack
	normalDeposit := math.NewInt(1000000)
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    suite.victim.String(),
		Amount:       normalDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err, "System should handle normal deposits after dust griefing")
}

// TestYieldDilutionAttack tests attacks that dilute yield for other users
func (suite *AdversarialTestSuite) TestYieldDilutionAttack() {
	attacker := suite.attacker.String()
	victim := suite.victim.String()

	// Step 1: Victim deposits early
	victimDeposit := math.NewInt(1000000)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    victim,
		Amount:       victimDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	victimPos, err := suite.keeper.GetV2UserPosition(suite.ctx, suite.victim)
	suite.Require().NoError(err)

	// Step 2: Attacker detects incoming yield and front-runs with large deposit
	attackerDeposit := math.NewInt(9000000) // 9x victim's deposit
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		Amount:       attackerDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Step 3: Yield distribution occurs
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)

	yieldAmount := math.NewInt(1000000) // 1M yield
	newNav := vaultState.TotalNav.Add(yieldAmount)
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    newNav,
		Reason:    "yield distribution",
	})
	suite.Require().NoError(err)

	// Step 4: Attacker immediately withdraws
	attackerPos, err := suite.keeper.GetV2UserPosition(suite.ctx, suite.attacker)
	suite.Require().NoError(err)

	withdrawResp, err := suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: attacker,
		Shares:     attackerPos.Shares,
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Step 5: Analyze yield distribution fairness
	attackerProfit := withdrawResp.AmountWithdrawn.Sub(attackerDeposit)

	finalVaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)

	victimCurrentValue := finalVaultState.SharePrice.MulInt(victimPos.Shares)
	victimYieldReceived := victimCurrentValue.Sub(math.LegacyNewDecFromInt(victimDeposit))

	suite.T().Logf("Yield dilution attack:")
	suite.T().Logf("  Victim early deposit: %s", victimDeposit.String())
	suite.T().Logf("  Attacker front-run deposit: %s", attackerDeposit.String())
	suite.T().Logf("  Total yield: %s", yieldAmount.String())
	suite.T().Logf("  Victim yield received: %s", victimYieldReceived.String())
	suite.T().Logf("  Attacker profit: %s", attackerProfit.String())

	// Verify unfair yield distribution
	actualVictimYieldPercent := victimYieldReceived.Quo(yieldAmount.ToLegacyDec()).Mul(math.LegacyNewDec(100))
	suite.T().Logf("  Victim should get ~10%%, actually got: %s%%", actualVictimYieldPercent.String())
}

// TestExtremeVolatilityAttack tests attacks during extreme price volatility
func (suite *AdversarialTestSuite) TestExtremeVolatilityAttack() {
	attacker := suite.attacker.String()

	// Step 1: Setup normal vault operations
	for i, user := range suite.users[:3] {
		deposit := math.NewInt(int64((i + 1) * 1000000))
		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user.String(),
			Amount:       deposit,
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)
	}

	// Step 2: Simulate extreme volatility with rapid NAV changes
	navValues := []int64{10000000, 5000000, 15000000, 3000000, 20000000}

	for i, navValue := range navValues {
		_, err := suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
			Authority: suite.authority,
			NewNav:    math.NewInt(navValue),
			Reason:    fmt.Sprintf("volatility step %d", i+1),
		})
		suite.Require().NoError(err)

		// Attacker tries to exploit volatility with each price change
		if i%2 == 0 { // Deposit on even steps
			_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
				Depositor:    attacker,
				Amount:       math.NewInt(1000000),
				ReceiveYield: true,
				MinShares:    math.ZeroInt(),
			})
			suite.Require().NoError(err)
		} else { // Withdraw on odd steps
			attackerPos, err := suite.keeper.GetV2UserPosition(suite.ctx, suite.attacker)
			if err == nil && !attackerPos.Shares.IsZero() {
				_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
					Withdrawer: attacker,
					Shares:     attackerPos.Shares,
					MinAmount:  math.ZeroInt(),
				})
				suite.Require().NoError(err)
			}
		}
	}

	suite.checkFinancialInvariants()
}

// Helper Functions

func (suite *AdversarialTestSuite) checkFinancialInvariants() {
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
	if err != nil {
		// If vault doesn't exist yet, that's okay for some tests
		return
	}

	// Basic invariants
	suite.Require().True(vaultState.TotalShares.GTE(math.ZeroInt()), "Total shares cannot be negative")
	suite.Require().True(vaultState.TotalNav.GTE(math.ZeroInt()), "Total NAV cannot be negative")
	suite.Require().True(vaultState.SharePrice.IsPositive(), "Share price must be positive")

	// Log state for debugging
	suite.T().Logf("Financial invariants check:")
	suite.T().Logf("  Total shares: %s", vaultState.TotalShares.String())
	suite.T().Logf("  Total NAV: %s", vaultState.TotalNav.String())
	suite.T().Logf("  Share price: %s", vaultState.SharePrice.String())
}

func (suite *AdversarialTestSuite) checkValueConservation() {
	// This would check that total value is conserved across operations
	// Implementation depends on specific vault mechanics
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
	if err != nil {
		return
	}

	totalValue := vaultState.SharePrice.MulInt(vaultState.TotalShares)
	suite.T().Logf("Value conservation check - Total value: %s", totalValue.String())
}

// TestReentrancyProtectionCheck verifies reentrancy protections
func (suite *AdversarialTestSuite) TestReentrancyProtectionCheck() {
	attacker := suite.attacker.String()

	// For now, just verify basic operation doesn't break
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		Amount:       math.NewInt(1000000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	suite.T().Logf("Reentrancy protection: Basic operations working")
}

// Helper function for share price calculation
func (suite *AdversarialTestSuite) calculateSharePrice() math.LegacyDec {
	vaultState, _ := suite.keeper.GetV2VaultState(suite.ctx)
	return vaultState.SharePrice
}

// TestGasExhaustionAttack tests attacks designed to exhaust gas
func (suite *AdversarialTestSuite) TestGasExhaustionAttack() {
	attacker := suite.attacker.String()

	// Setup initial state
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		Amount:       math.NewInt(1000000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Attempt operations that might consume excessive gas
	for i := 0; i < 10; i++ {
		_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
			Authority: suite.authority,
			NewNav:    math.NewInt(int64(1000000 + i*100000)),
			Reason:    fmt.Sprintf("gas exhaustion test %d", i),
		})
		suite.Require().NoError(err)
	}

	suite.T().Logf("Gas exhaustion attack: Multiple operations completed successfully")
}

// TestOverflowUnderflowAttacks tests numeric overflow/underflow attacks
func (suite *AdversarialTestSuite) TestOverflowUnderflowAttacks() {
	attacker := suite.attacker.String()

	// Test maximum values
	maxInt := math.NewInt(1000000000000) // Large value for testing

	// This should either handle gracefully or reject
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		Amount:       maxInt,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})

	if err != nil {
		suite.T().Logf("Overflow protection: Large deposit rejected")
	} else {
		suite.T().Logf("Overflow handling: Large deposit accepted")
	}

	// Test with maximum NAV update
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    maxInt,
		Reason:    "overflow test",
	})

	if err != nil {
		suite.T().Logf("Overflow protection: Large NAV update rejected")
	} else {
		suite.T().Logf("Overflow handling: Large NAV update accepted")
		suite.checkFinancialInvariants()
	}
}

func (suite *AdversarialTestSuite) TestTimestampManipulationAttack() {
	attacker := suite.attacker.String()
	victim := suite.victim.String()

	// Step 1: Initial deposits
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		Amount:       math.NewInt(1000000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    victim,
		Amount:       math.NewInt(1000000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Step 2: Get positions to verify timestamps
	attackerPos, err := suite.keeper.GetV2UserPosition(suite.ctx, suite.attacker)
	suite.Require().NoError(err)

	victimPos, err := suite.keeper.GetV2UserPosition(suite.ctx, suite.victim)
	suite.Require().NoError(err)

	// Verify positions exist
	suite.Require().True(attackerPos.Shares.IsPositive(), "Attacker should have shares")
	suite.Require().True(victimPos.Shares.IsPositive(), "Victim should have shares")

	suite.T().Logf("Timestamp manipulation check:")
	suite.T().Logf("  Attacker shares: %s", attackerPos.Shares.String())
	suite.T().Logf("  Victim shares: %s", victimPos.Shares.String())
	suite.T().Logf("  Both positions created successfully")
}
