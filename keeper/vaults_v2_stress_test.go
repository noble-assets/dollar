package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"dollar.noble.xyz/v2/keeper"
	"dollar.noble.xyz/v2/types/vaults"
	vaultsv2 "dollar.noble.xyz/v2/types/vaults/v2"
	"dollar.noble.xyz/v2/utils/mocks"
)

// StressTestSuite focuses on extreme conditions, boundary cases, and high-stress scenarios
type StressTestSuite struct {
	suite.Suite
	keeper    *keeper.Keeper
	ctx       sdk.Context
	msgServer vaultsv2.MsgServer
	authority string
	users     []sdk.AccAddress
	bank      mocks.BankKeeper
	account   mocks.AccountKeeper
}

func TestStressSuite(t *testing.T) {
	suite.Run(t, new(StressTestSuite))
}

func (suite *StressTestSuite) SetupTest() {
	// Setup test environment with proper mocks
	suite.account = mocks.AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	suite.bank = mocks.BankKeeper{
		Balances: make(map[string]sdk.Coins),
	}

	suite.keeper, _, suite.ctx = mocks.DollarKeeperWithKeepers(suite.T(), suite.bank, suite.account)
	suite.msgServer = keeper.NewVaultV2MsgServer(suite.keeper)
	suite.authority = "authority"

	// Create test addresses
	suite.users = make([]sdk.AccAddress, 20) // More users for stress testing
	for i := range suite.users {
		suite.users[i] = sdk.AccAddress(fmt.Sprintf("stress_user_%02d_____", i))
	}
}

// Precision and Rounding Stress Tests

func (suite *StressTestSuite) TestExtremePrecisionStress() {
	vaultType := vaults.FLEXIBLE
	user1 := suite.users[0].String()
	user2 := suite.users[1].String()

	// Create scenario with maximum precision stress
	// Start with 1 wei deposit
	microDeposit := math.NewInt(1)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user1,
		VaultType:    vaultType,
		Amount:       microDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Apply extreme NAV increase
	extremeNav := math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(35), nil))
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    extremeNav,
		Reason:    "extreme precision test",
	})
	suite.Require().NoError(err)

	// Test multiple deposits at extreme share price
	for i := 0; i < 10; i++ {
		depositAmount := math.NewInt(int64(1 + i)) // 1, 2, 3, ... 10
		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user2,
			VaultType:    vaultType,
			Amount:       depositAmount,
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)

		// Verify invariants after each deposit
		suite.checkInvariants(vaultType)
	}

	// Test withdrawal precision
	user2Position, _ := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, suite.users[1])
	if !user2Position.Shares.IsZero() {
		_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
			Withdrawer: user2,
			VaultType:  vaultType,
			Shares:     user2Position.Shares.QuoRaw(2), // Withdraw half
			MinAmount:  math.ZeroInt(),
		})
		suite.Require().NoError(err)
		suite.checkInvariants(vaultType)
	}
}

func (suite *StressTestSuite) TestRoundingAccumulationStress() {
	vaultType := vaults.FLEXIBLE

	// Create fractional share price to maximize rounding effects
	baseUser := suite.users[0].String()
	baseDeposit := math.NewInt(1000000)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    baseUser,
		VaultType:    vaultType,
		Amount:       baseDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Set NAV to create fractional share price
	fractionalNav := math.NewInt(1000003) // Creates 1.000003 share price
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    fractionalNav,
		Reason:    "fractional price setup",
	})
	suite.Require().NoError(err)

	// Get initial vault state after setup (this includes the base deposit + yield)
	initialVaultState, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	initialVaultValue := initialVaultState.TotalNav

	// Track only incremental deposits and withdrawals during the test rounds
	incrementalDeposited := math.ZeroInt()
	incrementalWithdrawn := math.ZeroInt()

	for round := 0; round < 100; round++ {
		if round%2 == 0 {
			// Deposit round - user deposits
			userIndex := (round / 2) % len(suite.users)
			user := suite.users[userIndex].String()
			amount := math.NewInt(int64(101 + round%7)) // Prime-ish numbers to force rounding
			_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
				Depositor:    user,
				VaultType:    vaultType,
				Amount:       amount,
				ReceiveYield: true,
				MinShares:    math.ZeroInt(),
			})
			suite.Require().NoError(err)
			incrementalDeposited = incrementalDeposited.Add(amount)
		} else {
			// Withdrawal round - try to withdraw from user who deposited in previous even round
			userIndex := ((round - 1) / 2) % len(suite.users)
			user := suite.users[userIndex].String()
			userAddr := suite.users[userIndex]
			position, err := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, userAddr)
			if err == nil && !position.Shares.IsZero() {
				withdrawShares := position.Shares.QuoRaw(3) // Withdraw 1/3
				if !withdrawShares.IsZero() {
					resp, err := suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
						Withdrawer: user,
						VaultType:  vaultType,
						Shares:     withdrawShares,
						MinAmount:  math.ZeroInt(),
					})
					if err == nil {
						incrementalWithdrawn = incrementalWithdrawn.Add(resp.AmountWithdrawn)
					}
				}
			}
		}

		// Check invariants every 10 rounds
		if round%10 == 9 {
			suite.checkInvariants(vaultType)
		}
	}

	// Final invariant check
	suite.checkInvariants(vaultType)

	// Verify no major value leakage due to rounding
	finalVaultState, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	finalVaultValue := finalVaultState.TotalNav

	// Expected final value = initial value + incremental deposits - incremental withdrawals
	expectedFinalValue := initialVaultValue.Add(incrementalDeposited).Sub(incrementalWithdrawn)

	diff := finalVaultValue.Sub(expectedFinalValue).Abs()
	maxExpectedRounding := math.NewInt(10000) // Allow for accumulated rounding from many operations

	suite.Require().True(diff.LTE(maxExpectedRounding),
		"Rounding accumulation too large: expected ~%s, got %s, diff %s",
		expectedFinalValue.String(), finalVaultValue.String(), diff.String())
}

// Large Number Stress Tests

func (suite *StressTestSuite) TestMassiveScaleOperations() {
	vaultType := vaults.FLEXIBLE

	// Test with very large numbers that approach overflow limits
	massiveAmount := math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(28), nil)) // 10^28

	user := suite.users[0].String()
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		VaultType:    vaultType,
		Amount:       massiveAmount,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Test NAV updates with massive numbers
	evenLargerNav := math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(30), nil)) // 10^30
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    evenLargerNav,
		Reason:    "massive scale test",
	})
	suite.Require().NoError(err)

	// Verify calculations still work
	suite.checkInvariants(vaultType)

	// Test withdrawal of massive amounts
	userPosition, _ := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, suite.users[0])
	_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: user,
		VaultType:  vaultType,
		Shares:     userPosition.Shares.QuoRaw(2),
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().NoError(err)

	suite.checkInvariants(vaultType)
}

func (suite *StressTestSuite) TestMixedScaleOperations() {
	vaultType := vaults.FLEXIBLE

	// Mix extremely large and extremely small operations
	scales := []int64{
		1,          // Minimum
		1000,       // Small
		1000000,    // Medium
		1000000000, // Large
	}

	// Add massive scale
	massiveScale := math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(25), nil))

	// Perform operations at all scales
	for i, scale := range scales {
		user := suite.users[i].String()
		amount := math.NewInt(scale)

		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user,
			VaultType:    vaultType,
			Amount:       amount,
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)
	}

	// Add massive deposit
	massiveUser := suite.users[4].String()
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    massiveUser,
		VaultType:    vaultType,
		Amount:       massiveScale,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Apply yield and verify all scales work together
	vaultState, _ := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	yieldNav := vaultState.TotalNav.MulRaw(11).QuoRaw(10) // 10% yield
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    yieldNav,
		Reason:    "mixed scale yield test",
	})
	suite.Require().NoError(err)

	// Test withdrawals at all scales
	for i := range scales {
		user := suite.users[i].String()
		userAddr := suite.users[i]
		position, err := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, userAddr)
		if err == nil && !position.Shares.IsZero() {
			withdrawShares := position.Shares.QuoRaw(3) // Withdraw 1/3
			if !withdrawShares.IsZero() {
				_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
					Withdrawer: user,
					VaultType:  vaultType,
					Shares:     withdrawShares,
					MinAmount:  math.ZeroInt(),
				})
				suite.Require().NoError(err)
			}
		}
	}

	suite.checkInvariants(vaultType)
}

// High-Frequency Operation Stress Tests

func (suite *StressTestSuite) TestHighFrequencyOperations() {
	vaultType := vaults.FLEXIBLE

	// Rapid-fire operations to stress state management
	numOperations := 1000
	userCount := len(suite.users)

	for i := 0; i < numOperations; i++ {
		user := suite.users[i%userCount].String()
		userAddr := suite.users[i%userCount]

		if i%3 == 0 {
			// Deposit
			amount := math.NewInt(int64(1000 + i%10000))
			_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
				Depositor:    user,
				VaultType:    vaultType,
				Amount:       amount,
				ReceiveYield: true,
				MinShares:    math.ZeroInt(),
			})
			suite.Require().NoError(err)
		} else if i%3 == 1 {
			// Withdraw (if user has shares)
			position, err := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, userAddr)
			if err == nil && !position.Shares.IsZero() {
				withdrawShares := position.Shares.QuoRaw(int64(2 + i%5)) // Withdraw fraction
				if !withdrawShares.IsZero() {
					_, err := suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
						Withdrawer: user,
						VaultType:  vaultType,
						Shares:     withdrawShares,
						MinAmount:  math.ZeroInt(),
					})
					suite.Require().NoError(err)
				}
			}
		} else {
			// NAV update (every third operation)
			vaultState, err := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
			if err == nil && !vaultState.TotalNav.IsZero() {
				// Small random NAV changes
				change := int64(95 + i%10) // 95-104, so -5% to +4%
				newNav := vaultState.TotalNav.MulRaw(change).QuoRaw(100)
				if newNav.IsPositive() {
					suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
						Authority: suite.authority,
						VaultType: vaultType,
						NewNav:    newNav,
						Reason:    fmt.Sprintf("high freq update %d", i),
					})
				}
			}
		}

		// Check invariants every 100 operations
		if i%100 == 99 {
			suite.checkInvariants(vaultType)
		}
	}

	// Final check
	suite.checkInvariants(vaultType)
}

func (suite *StressTestSuite) TestConcurrentUserSimulation() {
	vaultType := vaults.FLEXIBLE

	// Simulate many users performing operations simultaneously
	// (In a real test environment, this would use goroutines)

	rounds := 50
	usersPerRound := len(suite.users)

	for round := 0; round < rounds; round++ {
		// Each round simulates concurrent operations
		for userIdx := 0; userIdx < usersPerRound; userIdx++ {
			user := suite.users[userIdx].String()
			userAddr := suite.users[userIdx]

			operationType := (round + userIdx) % 4

			switch operationType {
			case 0, 1: // Deposit (50% of operations)
				amount := math.NewInt(int64(1000 * (1 + userIdx + round)))
				suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
					Depositor:    user,
					VaultType:    vaultType,
					Amount:       amount,
					ReceiveYield: true,
					MinShares:    math.ZeroInt(),
				})

			case 2: // Withdraw (25% of operations)
				position, err := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, userAddr)
				if err == nil && !position.Shares.IsZero() {
					withdrawShares := position.Shares.QuoRaw(int64(2 + round%3))
					if !withdrawShares.IsZero() {
						suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
							Withdrawer: user,
							VaultType:  vaultType,
							Shares:     withdrawShares,
							MinAmount:  math.ZeroInt(),
						})
					}
				}

			case 3: // Yield preference change (25% of operations)
				position, err := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, userAddr)
				if err == nil {
					newPref := !position.ReceiveYield
					suite.msgServer.SetYieldPreference(suite.ctx, &vaultsv2.MsgSetYieldPreference{
						User:         user,
						VaultType:    vaultType,
						ReceiveYield: newPref,
					})
				}
			}
		}

		// Random NAV update each round
		if round%3 == 0 {
			vaultState, err := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
			if err == nil && !vaultState.TotalNav.IsZero() {
				// Random yield between -2% and +5%
				yieldBps := int64(9800 + (round*7)%700) // 98.00% to 105.00%
				newNav := vaultState.TotalNav.MulRaw(yieldBps).QuoRaw(10000)
				if newNav.IsPositive() {
					suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
						Authority: suite.authority,
						VaultType: vaultType,
						NewNav:    newNav,
						Reason:    fmt.Sprintf("concurrent round %d", round),
					})
				}
			}
		}

		// Check invariants every 10 rounds
		if round%10 == 9 {
			suite.checkInvariants(vaultType)
		}
	}

	suite.checkInvariants(vaultType)
}

// Edge Case Boundary Tests

func (suite *StressTestSuite) TestZeroAndNearZeroStates() {
	vaultType := vaults.FLEXIBLE

	// Test operations when vault is in various near-zero states
	user := suite.users[0].String()

	// Start with minimum deposit
	minDeposit := math.NewInt(1)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		VaultType:    vaultType,
		Amount:       minDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Test NAV going to near-zero
	nearZeroNav := math.NewInt(1)
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		VaultType: vaultType,
		NewNav:    nearZeroNav,
		Reason:    "near zero test",
	})
	suite.Require().NoError(err)

	// Test operations at near-zero NAV
	user2 := suite.users[1].String()
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user2,
		VaultType:    vaultType,
		Amount:       math.NewInt(1000000), // Large deposit at near-zero share price
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	suite.checkInvariants(vaultType)

	// Test withdrawal when NAV is near zero
	user2Position, _ := suite.keeper.GetV2UserPosition(suite.ctx, vaultType, suite.users[1])
	_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
		Withdrawer: user2,
		VaultType:  vaultType,
		Shares:     user2Position.Shares.QuoRaw(2),
		MinAmount:  math.ZeroInt(),
	})
	suite.Require().NoError(err)

	suite.checkInvariants(vaultType)
}

func (suite *StressTestSuite) TestExtremeVolatilityStress() {
	vaultType := vaults.FLEXIBLE
	user := suite.users[0].String()

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

	// Apply extreme volatility sequence
	volatilitySequence := []struct {
		multiplier int64
		reason     string
	}{
		{10000, "10000x pump"},   // 10,000x increase
		{1, "crash to 1"},        // Crash to 1
		{1000000, "1M recovery"}, // 1M recovery
		{500000, "50% dump"},     // 50% dump
		{2000000, "4x pump"},     // 4x pump
		{100000, "95% crash"},    // 95% crash
	}

	baseNav := initialDeposit
	for i, vol := range volatilitySequence {
		newNav := baseNav.MulRaw(vol.multiplier).QuoRaw(1000000) // Normalize to reasonable range
		if newNav.IsZero() {
			newNav = math.NewInt(1) // Prevent zero NAV
		}

		_, err := suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
			Authority: suite.authority,
			VaultType: vaultType,
			NewNav:    newNav,
			Reason:    fmt.Sprintf("volatility step %d: %s", i, vol.reason),
		})
		suite.Require().NoError(err)

		// Test operations during extreme volatility
		if i%2 == 0 {
			// Test deposit during volatility
			testUser := suite.users[(i+1)%len(suite.users)].String()
			_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
				Depositor:    testUser,
				VaultType:    vaultType,
				Amount:       math.NewInt(10000),
				ReceiveYield: true,
				MinShares:    math.ZeroInt(),
			})
			suite.Require().NoError(err)
		}

		suite.checkInvariants(vaultType)
	}
}

// Mathematical Edge Cases

func (suite *StressTestSuite) TestMathematicalEdgeCases() {
	vaultType := vaults.FLEXIBLE

	// Test sequence of operations that could cause mathematical issues
	testCases := []struct {
		name           string
		shares         math.Int
		nav            math.Int
		expectedIssues string
	}{
		{
			name:   "High precision shares",
			shares: math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
			nav:    math.NewInt(1),
		},
		{
			name:   "High precision NAV",
			shares: math.NewInt(1),
			nav:    math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
		},
		{
			name:   "Both large",
			shares: math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(15), nil)),
			nav:    math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(15), nil)),
		},
		{
			name:   "Ratio stress",
			shares: math.NewInt(999999999),
			nav:    math.NewInt(1000000001),
		},
	}

	for i, tc := range testCases {
		// Reset vault state
		suite.SetupTest()

		user := suite.users[0].String()

		// Set up test case by manipulating vault state
		// First create minimal vault
		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user,
			VaultType:    vaultType,
			Amount:       math.NewInt(1),
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)

		// Create test state through NAV updates
		_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
			Authority: suite.authority,
			VaultType: vaultType,
			NewNav:    tc.nav,
			Reason:    fmt.Sprintf("test case %s setup", tc.name),
		})
		suite.Require().NoError(err)

		// Test operations
		testUser := suite.users[1].String()
		_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    testUser,
			VaultType:    vaultType,
			Amount:       math.NewInt(1000000),
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})

		if err != nil {
			suite.T().Logf("Test case %d (%s) failed deposit: %s", i, tc.name, err.Error())
		} else {
			suite.T().Logf("Test case %d (%s) passed deposit", i, tc.name)
			suite.checkInvariants(vaultType)
		}
	}
}

// Helper Functions

func (suite *StressTestSuite) checkInvariants(vaultType vaults.VaultType) {
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	if err != nil {
		// If vault doesn't exist yet, that's okay for some tests
		return
	}

	// Core invariants
	suite.Require().False(vaultState.TotalNav.IsNegative(), "Total NAV cannot be negative")
	suite.Require().False(vaultState.TotalShares.IsNegative(), "Total shares cannot be negative")

	if !vaultState.TotalShares.IsZero() {
		suite.Require().True(vaultState.SharePrice.IsPositive(), "Share price must be positive when shares exist")

		// Share price calculation consistency (with tolerance for precision)
		expectedSharePrice := suite.calculateSharePrice(vaultState)
		priceDiff := expectedSharePrice.Sub(vaultState.SharePrice).Abs()
		tolerance := math.LegacyNewDecWithPrec(1, 12) // 1e-12 tolerance

		suite.Require().True(priceDiff.LTE(tolerance),
			"Share price calculation inconsistent: expected %s, got %s, diff %s",
			expectedSharePrice.String(), vaultState.SharePrice.String(), priceDiff.String())

		// Verify no overflow in value calculations
		totalValue := vaultState.SharePrice.MulInt(vaultState.TotalShares)
		navDec := math.LegacyNewDecFromInt(vaultState.TotalNav)
		valueDiff := totalValue.Sub(navDec).Abs()

		suite.Require().True(valueDiff.LTE(math.LegacyOneDec()),
			"Total value calculation overflow or major precision loss")
	}

	// Check for reasonable bounds (catch extreme states)
	maxReasonableNav := math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(40), nil))
	suite.Require().True(vaultState.TotalNav.LTE(maxReasonableNav),
		"NAV exceeds reasonable bounds: %s", vaultState.TotalNav.String())

	maxReasonableShares := math.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(40), nil))
	suite.Require().True(vaultState.TotalShares.LTE(maxReasonableShares),
		"Shares exceed reasonable bounds: %s", vaultState.TotalShares.String())
}

// Helper function to calculate share price (replicates the keeper's unexported method)
func (suite *StressTestSuite) calculateSharePrice(vaultState *vaultsv2.VaultState) math.LegacyDec {
	if vaultState.TotalShares.IsZero() {
		return math.LegacyOneDec() // 1:1 ratio for first deposit
	}
	return math.LegacyNewDecFromInt(vaultState.TotalNav).Quo(math.LegacyNewDecFromInt(vaultState.TotalShares))
}

func (suite *StressTestSuite) logVaultState(vaultType vaults.VaultType, context string) {
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx, vaultType)
	if err != nil {
		suite.T().Logf("%s: Failed to get vault state: %s", context, err.Error())
		return
	}

	suite.T().Logf("%s: Vault State", context)
	suite.T().Logf("  Total Shares: %s", vaultState.TotalShares.String())
	suite.T().Logf("  Total NAV: %s", vaultState.TotalNav.String())
	suite.T().Logf("  Share Price: %s", vaultState.SharePrice.String())
	suite.T().Logf("  Total Users: %d", vaultState.TotalUsers)
}
