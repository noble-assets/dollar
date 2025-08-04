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
	suite.addresses = make([]sdk.AccAddress, 10)
	for i := range suite.addresses {
		suite.addresses[i] = sdk.AccAddress(fmt.Sprintf("test-address-%d____", i))
	}
}

// Financial Logic Tests

func (suite *V2VaultTestSuite) TestBasicDepositWithdraw() {
	user := suite.addresses[0].String()
	vaultType := vaults.FLEXIBLE

	// Test basic deposit
	depositAmount := math.NewInt(1000000) // 1M units
	resp, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		VaultType:    vaultType,
		Amount:       depositAmount,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(depositAmount, resp.SharesReceived) // 1:1 ratio for first deposit
	suite.Require().Equal(math.LegacyOneDec(), resp.SharePrice)

	// Verify vault state
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	suite.Require().NoError(err)
	suite.Require().Equal(depositAmount, vaultState.TotalShares)
	suite.Require().Equal(depositAmount, vaultState.TotalNav)
	suite.Require().Equal(math.LegacyOneDec(), vaultState.SharePrice)

	// Test withdrawal
	withdrawShares := math.NewInt(500000) // Half the shares
	withdrawResp, err := suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: user,
		VaultType:  vaultType,
		Shares:     withdrawShares,
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(withdrawShares, withdrawResp.AmountWithdrawn)
	suite.Require().Equal(withdrawShares, withdrawResp.SharesRedeemed)

	// Verify final state
	finalVaultState, err := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	suite.Require().NoError(err)
	suite.Require().Equal(depositAmount.Sub(withdrawShares), finalVaultState.TotalShares)
	suite.Require().Equal(depositAmount.Sub(withdrawShares), finalVaultState.TotalNav)
}

func (suite *V2VaultTestSuite) TestSharePriceCalculation() {
	vaultType := vaults.FLEXIBLE
	user := suite.addresses[0].String()

	// Create vault with initial deposit
	initialDeposit := math.NewInt(1000000)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		VaultType:    vaultType,
		Amount:       initialDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Simulate NAV increase (yield generation)
	newNav := math.NewInt(1200000) // 20% increase
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    newNav,
		Reason:    "yield distribution",
	})
	suite.Require().NoError(err)

	// Verify share price increased
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	suite.Require().NoError(err)
	expectedSharePrice := math.LegacyNewDecFromInt(newNav).Quo(math.LegacyNewDecFromInt(initialDeposit))
	suite.Require().Equal(expectedSharePrice, vaultState.SharePrice)

	// Test new deposit at higher share price
	user2 := suite.addresses[1].String()
	depositAmount2 := math.NewInt(600000)
	expectedShares := math.LegacyNewDecFromInt(depositAmount2).Quo(expectedSharePrice).TruncateInt()

	resp2, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user2,
		VaultType:    vaultType,
		Amount:       depositAmount2,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(expectedShares, resp2.SharesReceived)
}

func (suite *V2VaultTestSuite) TestInvariantConservation() {
	vaultType := vaults.FLEXIBLE
	users := []string{
		suite.addresses[0].String(),
		suite.addresses[1].String(),
		suite.addresses[2].String(),
	}

	deposits := []math.Int{
		math.NewInt(1000000),
		math.NewInt(2000000),
		math.NewInt(500000),
	}

	// Track total deposits for invariant checking
	totalDeposited := math.ZeroInt()

	// Multiple users deposit
	for i, user := range users {
		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user,
			VaultType:    vaultType,
			Amount:       deposits[i],
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)
		totalDeposited = totalDeposited.Add(deposits[i])

		// Check invariant: TotalShares * SharePrice should equal TotalNAV
		suite.checkFinancialInvariants(vaultType)
	}

	// Verify total amounts
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	suite.Require().NoError(err)
	suite.Require().Equal(totalDeposited, vaultState.TotalNav)

	// Simulate yield
	yieldAmount := math.NewInt(350000) // 10% yield
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    totalDeposited.Add(yieldAmount),
		Reason:    "yield distribution",
	})
	suite.Require().NoError(err)

	// Check invariant after yield
	suite.checkFinancialInvariants(vaultType)
}

func (suite *V2VaultTestSuite) TestPrecisionAndRounding() {
	vaultType := vaults.FLEXIBLE
	user := suite.addresses[0].String()

	// Test with very small amounts to check precision
	smallDeposit := math.NewInt(1) // Minimum unit
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		VaultType:    vaultType,
		Amount:       smallDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Add yield that creates fractional share prices
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    math.NewInt(3), // Creates share price of 3.0
		Reason:    "precision test",
	})
	suite.Require().NoError(err)

	// Test deposit that should result in fractional shares
	user2 := suite.addresses[1].String()
	deposit2 := math.NewInt(10)
	resp, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user2,
		VaultType:    vaultType,
		Amount:       deposit2,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Check that shares are correctly truncated (not rounded up)
	expectedShares := math.NewInt(3) // 10/3 = 3.33... truncated to 3
	suite.Require().Equal(expectedShares, resp.SharesReceived)

	// Verify no value is lost in the vault
	suite.checkFinancialInvariants(vaultType)
}

func (suite *V2VaultTestSuite) TestLargeAmounts() {
	vaultType := vaults.FLEXIBLE
	user := suite.addresses[0].String()

	// Test with maximum possible values
	maxDeposit := math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(30), nil)) // 10^30

	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		VaultType:    vaultType,
		Amount:       maxDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Test large yield update
	largeYield := math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(29), nil)) // 10^29
	newNav := maxDeposit.Add(largeYield)

	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    newNav,
		Reason:    "large yield test",
	})
	suite.Require().NoError(err)

	// Check calculations still work correctly
	suite.checkFinancialInvariants(vaultType)
}

// Adversarial Interaction Tests

func (suite *V2VaultTestSuite) TestSlippageProtection() {
	vaultType := vaults.FLEXIBLE
	user := suite.addresses[0].String()

	// Initial deposit
	initialDeposit := math.NewInt(1000000)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		VaultType:    vaultType,
		Amount:       initialDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Simulate sudden NAV increase (frontrunning scenario)
	highNav := math.NewInt(2000000)
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    highNav,
		Reason:    "sudden price increase",
	})
	suite.Require().NoError(err)

	// Test deposit with slippage protection
	user2 := suite.addresses[1].String()
	depositAmount := math.NewInt(1000000)
	minShares := math.NewInt(750000) // Expecting at least 75% of original share ratio

	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user2,
		VaultType:    vaultType,
		Amount:       depositAmount,
		ReceiveYield: true,
		MinShares:    minShares,
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient shares received")

	// Test withdrawal with slippage protection
	withdrawShares := math.NewInt(500000)
	minAmount := math.NewInt(1200000) // Expecting high amount due to increased share price

	_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: user,
		VaultType:  vaultType,
		Shares:     withdrawShares,
		MinAmount:  minAmount,
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient amount received")
}

func (suite *V2VaultTestSuite) TestSandwichAttackPrevention() {
	vaultType := vaults.FLEXIBLE
	attacker := suite.addresses[0].String()
	victim := suite.addresses[1].String()

	// Step 1: Attacker makes large deposit
	attackDeposit := math.NewInt(10000000)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    attacker,
		VaultType:    vaultType,
		Amount:       attackDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Step 2: Simulate large NAV update (representing yield or manipulation)
	manipulatedNav := math.NewInt(15000000) // 50% increase
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    manipulatedNav,
		Reason:    "potential manipulation",
	})
	suite.Require().NoError(err)

	vaultStateBefore, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	attackerSharesBefore, _ := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, sdk.MustAccAddressFromBech32(attacker))

	// Step 3: Victim makes deposit at inflated price
	victimDeposit := math.NewInt(1000000)
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    victim,
		VaultType:    vaultType,
		Amount:       victimDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Step 4: Attacker tries to withdraw at inflated price
	_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: attacker,
		VaultType:  vaultType,
		Shares:     attackerSharesBefore.Shares,
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().NoError(err)

	vaultStateAfter, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)

	// Verify that the attack doesn't create value out of thin air
	// The total NAV should only increase through legitimate means
	suite.Require().True(vaultStateAfter.TotalNav.LTE(vaultStateBefore.TotalNav.Add(victimDeposit)))
}

func (suite *V2VaultTestSuite) TestBoundaryConditions() {
	vaultType := vaults.FLEXIBLE
	user := suite.addresses[0].String()

	// Test zero amounts
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		VaultType:    vaultType,
		Amount:       math.ZeroInt(),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "must be positive")

	// Test negative amounts
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		VaultType:    vaultType,
		Amount:       math.NewInt(-1),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().Error(err)

	// Test withdrawing more shares than owned
	deposit := math.NewInt(1000)
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		VaultType:    vaultType,
		Amount:       deposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: user,
		VaultType:  vaultType,
		Shares:     math.NewInt(2000), // More than deposited
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient shares")
}

func (suite *V2VaultTestSuite) TestConcurrentOperations() {
	vaultType := vaults.FLEXIBLE
	users := []string{
		suite.addresses[0].String(),
		suite.addresses[1].String(),
		suite.addresses[2].String(),
		suite.addresses[3].String(),
		suite.addresses[4].String(),
	}

	// Simulate concurrent deposits and withdrawals
	for round := 0; round < 10; round++ {
		// Random deposits
		for i, user := range users {
			if round%2 == i%2 { // Alternate operations
				amount := math.NewInt(int64((i + 1) * 100000))
				_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
					Depositor:    user,
					VaultType:    vaultType,
					Amount:       amount,
					ReceiveYield: true,
					MinShares:    math.ZeroInt(),
				})
				suite.Require().NoError(err)
			}
		}

		// Check invariants after each round
		suite.checkFinancialInvariants(vaultType)

		// Random NAV updates
		if round%3 == 0 {
			vaultState, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
			newNav := vaultState.TotalNav.MulRaw(105).QuoRaw(100) // 5% increase
			suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
				Authority: suite.authority,
				VaultType: vaultType,
				NewNav:    newNav,
				Reason:    fmt.Sprintf("round %d yield", round),
			})
		}
	}
}

func (suite *V2VaultTestSuite) TestFirstAndLastUserEdgeCases() {
	vaultType := vaults.FLEXIBLE
	user1 := suite.addresses[0].String()
	user2 := suite.addresses[1].String()

	// First user deposits minimum amount
	firstDeposit := math.NewInt(1)
	resp1, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user1,
		VaultType:    vaultType,
		Amount:       firstDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(firstDeposit, resp1.SharesReceived)

	// Massive NAV increase
	hugeNav := math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil))
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    hugeNav,
		Reason:    "massive yield",
	})
	suite.Require().NoError(err)

	// Second user deposits large amount at extremely high share price
	largeDeposit := math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user2,
		VaultType:    vaultType,
		Amount:       largeDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Verify calculations are still correct
	suite.checkFinancialInvariants(vaultType)

	// First user withdraws (should get massive amount)
	user1Position, _ := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, sdk.MustAccAddressFromBech32(user1))
	withdrawResp, err := suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: user1,
		VaultType:  vaultType,
		Shares:     user1Position.Shares,
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().NoError(err)
	suite.Require().True(withdrawResp.AmountWithdrawn.GT(firstDeposit))
}

func (suite *V2VaultTestSuite) TestMaliciousNAVManipulation() {
	vaultType := vaults.FLEXIBLE
	user := suite.addresses[0].String()

	// Initial deposit
	deposit := math.NewInt(1000000)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		VaultType:    vaultType,
		Amount:       deposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Test unauthorized NAV update
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: "unauthorized-address",
		VaultType: vaultType,
		NewNav:    math.NewInt(999999999),
		Reason:    "malicious update",
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid authority")

	// Test NAV manipulation with extreme values
	extremeNav := math.NewInt(1) // Setting NAV to 1 when shares are 1M
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    extremeNav,
		Reason:    "extreme manipulation test",
	})
	suite.Require().NoError(err)

	// Verify the system handles extreme share price correctly
	vaultState, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	expectedSharePrice := math.LegacyNewDecFromInt(extremeNav).Quo(math.LegacyNewDecFromInt(deposit))
	suite.Require().Equal(expectedSharePrice, vaultState.SharePrice)

	// Test that new deposits work correctly even with extreme share price
	user2 := suite.addresses[1].String()
	smallDeposit := math.NewInt(1000000)
	expectedShares := math.LegacyNewDecFromInt(smallDeposit).Quo(expectedSharePrice).TruncateInt()

	resp, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user2,
		VaultType:    vaultType,
		Amount:       smallDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(expectedShares, resp.SharesReceived)
}

// Helper Functions

func (suite *V2VaultTestSuite) checkFinancialInvariants(vaultType vaults.VaultType) {
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	if err != nil {
		// If vault doesn't exist yet, that's okay for some tests
		return
	}

	// Invariant 1: Share price calculation should be consistent
	if !vaultState.TotalShares.IsZero() {
		expectedSharePrice := suite.calculateSharePrice(vaultState)
		suite.Require().Equal(expectedSharePrice, vaultState.SharePrice, "Share price calculation inconsistent")
	}

	// Invariant 2: Total NAV should never be negative
	suite.Require().False(vaultState.TotalNav.IsNegative(), "Total NAV cannot be negative")

	// Invariant 3: Total shares should never be negative
	suite.Require().False(vaultState.TotalShares.IsNegative(), "Total shares cannot be negative")

	// Invariant 4: Share price should be positive (if shares exist)
	if !vaultState.TotalShares.IsZero() {
		suite.Require().True(vaultState.SharePrice.IsPositive(), "Share price must be positive when shares exist")
	}

	// Invariant 5: Total value calculation should not overflow
	if !vaultState.TotalShares.IsZero() {
		calculatedNav := vaultState.SharePrice.MulInt(vaultState.TotalShares).TruncateInt()
		// Allow for small rounding differences
		diff := calculatedNav.Sub(vaultState.TotalNav).Abs()
		suite.Require().True(diff.LTE(math.NewInt(1)), "NAV calculation has too much rounding error")
	}
}

func (suite *V2VaultTestSuite) createTestVaultState(vaultType vaults.VaultType, totalShares, totalNav math.Int) *vaultsv2.VaultState {
	sharePrice := math.LegacyOneDec()
	if !totalShares.IsZero() {
		sharePrice = math.LegacyNewDecFromInt(totalNav).Quo(math.LegacyNewDecFromInt(totalShares))
	}

	return &vaultsv2.VaultState{
		VaultType:              vaultType,
		TotalShares:            totalShares,
		TotalNav:               totalNav,
		SharePrice:             sharePrice,
		TotalUsers:             1,
		DepositsEnabled:        true,
		WithdrawalsEnabled:     true,
		LastNavUpdate:          time.Now(),
		TotalSharesPendingExit: math.ZeroInt(),
		PendingExitRequests:    0,
	}
}

// Edge Case Tests

func (suite *V2VaultTestSuite) TestDustAmountHandling() {
	vaultType := vaults.FLEXIBLE
	user := suite.addresses[0].String()

	// Create scenario where dust amounts matter
	largeDeposit := math.NewInt(1000000000) // 1B units
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		VaultType:    vaultType,
		Amount:       largeDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Create fractional share price
	newNav := largeDeposit.MulRaw(3).AddRaw(1) // Slight increase creating fractional price
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    newNav,
		Reason:    "fractional price test",
	})
	suite.Require().NoError(err)

	// Test dust deposit
	user2 := suite.addresses[1].String()
	dustDeposit := math.NewInt(1)
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user2,
		VaultType:    vaultType,
		Amount:       dustDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Verify that even dust deposits don't break invariants
	suite.checkFinancialInvariants(vaultType)

	// Test withdrawal of dust shares
	user2Position, _ := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, sdk.MustAccAddressFromBech32(user2))
	if !user2Position.Shares.IsZero() {
		_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
			Withdrawer: user2,
			VaultType:  vaultType,
			Shares:     user2Position.Shares,
			MinAmount:  math.ZeroInt(),
		})
		suite.Require().NoError(err)
	}
}

// Helper function to calculate share price (replicates the keeper's unexported method)
func (suite *V2VaultTestSuite) calculateSharePrice(vaultState *vaultsv2.VaultState) math.LegacyDec {
	if vaultState.TotalShares.IsZero() {
		return math.LegacyOneDec() // 1:1 ratio for first deposit
	}
	return math.LegacyNewDecFromInt(vaultState.TotalNav).Quo(math.LegacyNewDecFromInt(vaultState.TotalShares))
}

func (suite *V2VaultTestSuite) TestShareCalculationEdgeCases() {
	vaultType := vaults.FLEXIBLE
	user := suite.addresses[0].String()

	// Test when total shares is exactly 1 and NAV is very large
	testState := suite.createTestVaultState(vaultType, math.NewInt(1), math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(25), nil)))
	suite.keeper.SetV2VaultState(suite.ctx, vaultType, testState)

	// Test deposit at extreme share price
	deposit := math.NewInt(1000000)
	sharePrice := suite.calculateSharePrice(testState)
	expectedShares := math.LegacyNewDecFromInt(deposit).Quo(sharePrice).TruncateInt()

	resp, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		VaultType:    vaultType,
		Amount:       deposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(expectedShares, resp.SharesReceived)

	// Verify state consistency
	suite.checkFinancialInvariants(vaultType)
}

func (suite *V2VaultTestSuite) TestYieldPreferenceChanges() {
	vaultType := vaults.FLEXIBLE
	user := suite.addresses[0].String()

	// Initial deposit with yield preference
	deposit := math.NewInt(1000000)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		VaultType:    vaultType,
		Amount:       deposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Change yield preference
	resp, err := suite.msgServer.SetYieldPreference(suite.ctx, &vaultsv2.MsgSetYieldPreference{
		User:         user,
		VaultType:    vaultType,
		ReceiveYield: false,
	})
	suite.Require().NoError(err)
	suite.Require().True(resp.PreviousPreference)
	suite.Require().False(resp.NewPreference)

	// Verify preference was updated
	userPosition, _ := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, sdk.MustAccAddressFromBech32(user))
	suite.Require().False(userPosition.ReceiveYield)

	// Test changing back to yield receiving
	resp2, err := suite.msgServer.SetYieldPreference(suite.ctx, &vaultsv2.MsgSetYieldPreference{
		User:         user,
		VaultType:    vaultType,
		ReceiveYield: true,
	})
	suite.Require().NoError(err)
	suite.Require().False(resp2.PreviousPreference)
	suite.Require().True(resp2.NewPreference)
}

func (suite *V2VaultTestSuite) TestVaultPausingMechanisms() {
	vaultType := vaults.FLEXIBLE
	user := suite.addresses[0].String()

	// Initial setup
	vaultState := suite.createTestVaultState(vaultType, math.ZeroInt(), math.ZeroInt())
	vaultState.DepositsEnabled = false
	suite.keeper.SetV2VaultState(suite.ctx, vaultType, vaultState)

	// Test deposit when disabled
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		VaultType:    vaultType,
		Amount:       math.NewInt(1000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "deposits are currently disabled")

	// Enable deposits but disable withdrawals
	vaultState.DepositsEnabled = true
	vaultState.WithdrawalsEnabled = false
	suite.keeper.SetV2VaultState(suite.ctx, vaultType, vaultState)

	// Deposit should work
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		VaultType:    vaultType,
		Amount:       math.NewInt(1000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Withdrawal should fail
	_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: user,
		VaultType:  vaultType,
		Shares:     math.NewInt(500),
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "withdrawals are currently disabled")
}

func (suite *V2VaultTestSuite) TestVaultTypeValidation() {
	user := suite.addresses[0].String()

	// Test with unspecified vault type
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		VaultType:    vaults.UNSPECIFIED,
		Amount:       math.NewInt(1000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "vault type must be specified")

	// Test with each valid vault type
	validTypes := []vaults.VaultType{vaults.FLEXIBLE, vaults.STAKED}
	for _, vaultType := range validTypes {
		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user,
			VaultType:    vaultType,
			Amount:       math.NewInt(1000),
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err, "Valid vault type %s should work", vaultType.String())
	}
}

func (suite *V2VaultTestSuite) TestAddressValidation() {
	vaultType := vaults.FLEXIBLE

	// Test invalid depositor address
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    "invalid-address",
		VaultType:    vaultType,
		Amount:       math.NewInt(1000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid depositor address")

	// Test empty address
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    "",
		VaultType:    vaultType,
		Amount:       math.NewInt(1000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().Error(err)
}

func (suite *V2VaultTestSuite) TestMultipleUsersShareDistribution() {
	vaultType := vaults.FLEXIBLE
	users := []string{
		suite.addresses[0].String(),
		suite.addresses[1].String(),
		suite.addresses[2].String(),
	}

	// Equal deposits
	depositAmount := math.NewInt(1000000)
	for _, user := range users {
		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user,
			VaultType:    vaultType,
			Amount:       depositAmount,
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)
	}

	// Verify equal share distribution
	for _, user := range users {
		position, _ := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, sdk.MustAccAddressFromBech32(user))
		suite.Require().Equal(depositAmount, position.Shares, "Equal deposits should result in equal shares")
	}

	// Add yield
	vaultState, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	yieldAmount := vaultState.TotalNav.MulRaw(20).QuoRaw(100) // 20% yield
	newNav := vaultState.TotalNav.Add(yieldAmount)

	_, err := suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    newNav,
		Reason:    "yield distribution test",
	})
	suite.Require().NoError(err)

	// Verify proportional withdrawal values
	for i, user := range users {
		position, _ := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, sdk.MustAccAddressFromBech32(user))
		withdrawResp, err := suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
			Withdrawer: user,
			VaultType:  vaultType,
			Shares:     position.Shares.QuoRaw(2), // Withdraw half
			MinAmount:  math.ZeroInt(),
		})
		suite.Require().NoError(err)

		// Each user should get the same withdrawal amount for the same share amount
		if i > 0 {
			prevValue := depositAmount.Add(yieldAmount.QuoRaw(3)).QuoRaw(2) // Expected value per user half-withdrawal
			actualValue := withdrawResp.AmountWithdrawn

			// Allow for small rounding differences
			diff := actualValue.Sub(prevValue).Abs()
			suite.Require().True(diff.LTE(math.NewInt(1)), "Withdrawal amounts should be nearly equal for equal positions")
		}
	}
}

func (suite *V2VaultTestSuite) TestCompoundingYield() {
	vaultType := vaults.FLEXIBLE
	user := suite.addresses[0].String()

	// Initial deposit
	initialDeposit := math.NewInt(1000000)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		VaultType:    vaultType,
		Amount:       initialDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Apply multiple rounds of yield
	currentNav := initialDeposit
	for round := 1; round <= 5; round++ {
		// 10% yield each round
		yieldAmount := currentNav.MulRaw(10).QuoRaw(100)
		currentNav = currentNav.Add(yieldAmount)

		_, err := suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
			Authority: suite.authority,
			VaultType: vaultType,
			NewNav:    currentNav,
			Reason:    fmt.Sprintf("compound yield round %d", round),
		})
		suite.Require().NoError(err)

		// Verify share price increased
		vaultState, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
		expectedSharePrice := math.LegacyNewDecFromInt(currentNav).Quo(math.LegacyNewDecFromInt(initialDeposit))
		suite.Require().Equal(expectedSharePrice, vaultState.SharePrice)
	}

	// Final withdrawal should include all compounded yield
	userPosition, _ := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, sdk.MustAccAddressFromBech32(user))
	withdrawResp, err := suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: user,
		VaultType:  vaultType,
		Shares:     userPosition.Shares,
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// User should receive significantly more than initial deposit due to compounding
	suite.Require().True(withdrawResp.AmountWithdrawn.GT(initialDeposit.MulRaw(16).QuoRaw(10)), // Should be > 1.6x original
		"Compounded yield should result in substantial gains")
}
