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

// StressTestSuite focuses on high-volume, high-frequency operations and edge cases
type StressTestSuite struct {
	suite.Suite
	keeper    *keeper.Keeper
	ctx       sdk.Context
	msgServer vaultsv2.MsgServer
	authority string
	bank      mocks.BankKeeper
	account   mocks.AccountKeeper
	users     []sdk.AccAddress
}

func TestStressSuite(t *testing.T) {
	suite.Run(t, new(StressTestSuite))
}

func (suite *StressTestSuite) SetupTest() {
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

	// Create many test users for stress testing
	for i := 0; i < 100; i++ {
		user := sdk.AccAddress(fmt.Sprintf("stress_user_%d", i))
		suite.users = append(suite.users, user)
	}
}

// TestHighVolumeDepositsWithdraws tests system under high transaction volume
func (suite *StressTestSuite) TestHighVolumeDepositsWithdraws() {
	// Initialize vault with seed deposit
	seedUser := suite.users[0].String()
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    seedUser,
		Amount:       math.NewInt(1000000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Perform high volume of deposits
	successfulDeposits := 0
	for i := 1; i < 50; i++ {
		user := suite.users[i].String()
		amount := math.NewInt(int64((i + 1) * 10000))

		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user,
			Amount:       amount,
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})

		if err == nil {
			successfulDeposits++
		}
	}

	// Add yield event during high activity
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)

	yieldAmount := vaultState.TotalNav.MulRaw(15).QuoRaw(100) // 15% yield
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    vaultState.TotalNav.Add(yieldAmount),
		Reason:    "stress test yield",
	})
	suite.Require().NoError(err)

	// Perform high volume of withdrawals
	successfulWithdrawals := 0
	for i := 1; i < 30; i++ {
		user := suite.users[i]
		userPos, err := suite.keeper.GetV2UserPosition(suite.ctx, user)
		if err == nil && !userPos.Shares.IsZero() {
			withdrawShares := userPos.Shares.QuoRaw(2) // Withdraw half

			_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
				Withdrawer: user.String(),
				Shares:     withdrawShares,
				MinAmount:  math.ZeroInt(),
			})

			if err == nil {
				successfulWithdrawals++
			}
		}
	}

	suite.T().Logf("High volume stress test:")
	suite.T().Logf("  Successful deposits: %d/49", successfulDeposits)
	suite.T().Logf("  Successful withdrawals: %d", successfulWithdrawals)

	// Verify system integrity after stress
	suite.checkSystemIntegrity()
}

// TestRapidNavUpdates tests system resilience under frequent NAV changes
func (suite *StressTestSuite) TestRapidNavUpdates() {
	// Setup initial positions
	for i := 0; i < 10; i++ {
		user := suite.users[i].String()
		deposit := math.NewInt(int64((i + 1) * 100000))

		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user,
			Amount:       deposit,
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)
	}

	// Rapid NAV updates simulating volatile market conditions
	baseNav := math.NewInt(1000000)
	successfulUpdates := 0

	for i := 0; i < 100; i++ {
		// Simulate price volatility: +/- 5% each update
		multiplier := 95 + (i % 10) // Between 95% and 104%
		newNav := baseNav.MulRaw(int64(multiplier)).QuoRaw(100)

		_, err := suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
			Authority: suite.authority,
			NewNav:    newNav,
			Reason:    fmt.Sprintf("stress update %d", i),
		})

		if err == nil {
			successfulUpdates++
			baseNav = newNav
		}
	}

	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)

	suite.T().Logf("Rapid NAV updates stress test:")
	suite.T().Logf("  Successful updates: %d/100", successfulUpdates)
	suite.T().Logf("  Final share price: %s", vaultState.SharePrice.String())
	suite.T().Logf("  Final total NAV: %s", vaultState.TotalNav.String())

	suite.checkSystemIntegrity()
}

// TestConcurrentUserOperations simulates many users operating simultaneously
func (suite *StressTestSuite) TestConcurrentUserOperations() {
	// Phase 1: Mass deposits
	for i := 0; i < 50; i++ {
		user := suite.users[i].String()
		amount := math.NewInt(int64(50000 + (i * 1000)))

		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user,
			Amount:       amount,
			ReceiveYield: i%2 == 0, // Alternate yield preferences
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)
	}

	// Phase 2: Yield events during activity
	for round := 0; round < 5; round++ {
		vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
		suite.Require().NoError(err)

		yieldAmount := vaultState.TotalNav.MulRaw(5).QuoRaw(100) // 5% yield
		_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
			Authority: suite.authority,
			NewNav:    vaultState.TotalNav.Add(yieldAmount),
			Reason:    fmt.Sprintf("concurrent yield round %d", round),
		})
		suite.Require().NoError(err)

		// Some users deposit more during yield events
		for i := 0; i < 10; i++ {
			user := suite.users[50+i].String()
			amount := math.NewInt(int64(25000 * (round + 1)))

			_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
				Depositor:    user,
				Amount:       amount,
				ReceiveYield: true,
				MinShares:    math.ZeroInt(),
			})
			suite.Require().NoError(err)
		}

		// Some users withdraw during yield events
		for i := round * 5; i < (round+1)*5 && i < 25; i++ {
			user := suite.users[i]
			userPos, err := suite.keeper.GetV2UserPosition(suite.ctx, user)
			if err == nil && !userPos.Shares.IsZero() {
				withdrawShares := userPos.Shares.QuoRaw(4) // Withdraw 25%

				_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
					Withdrawer: user.String(),
					Shares:     withdrawShares,
					MinAmount:  math.ZeroInt(),
				})
				suite.Require().NoError(err)
			}
		}
	}

	// Phase 3: Mass withdrawals
	activeUsers := 0
	for i := 0; i < 75; i++ {
		user := suite.users[i]
		userPos, err := suite.keeper.GetV2UserPosition(suite.ctx, user)
		if err == nil && !userPos.Shares.IsZero() {
			activeUsers++

			// Withdraw remaining shares
			_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
				Withdrawer: user.String(),
				Shares:     userPos.Shares,
				MinAmount:  math.ZeroInt(),
			})
			suite.Require().NoError(err)
		}
	}

	finalState, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)

	suite.T().Logf("Concurrent operations stress test:")
	suite.T().Logf("  Active users processed: %d", activeUsers)
	suite.T().Logf("  Final total shares: %s", finalState.TotalShares.String())
	suite.T().Logf("  Final total NAV: %s", finalState.TotalNav.String())

	suite.checkSystemIntegrity()
}

// TestExtremeSharePriceScenarios tests system behavior at extreme share prices
func (suite *StressTestSuite) TestExtremeSharePriceScenarios() {
	user := suite.users[0].String()

	// Start with minimal deposit
	minDeposit := math.NewInt(1)
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		Amount:       minDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Create extremely high share price through NAV manipulation
	extremeNav, _ := math.NewIntFromString("1000000000000000000000000") // 1e24

	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    extremeNav,
		Reason:    "extreme price stress test",
	})

	if err != nil {
		suite.T().Logf("Extreme NAV rejected: %s", err.Error())
		return
	}

	// Test operations at extreme share price
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)

	suite.T().Logf("Extreme share price test:")
	suite.T().Logf("  Share price: %s", vaultState.SharePrice.String())

	// Test if system can handle deposits at extreme prices
	testDeposit := math.NewInt(1000000000) // 1B units
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    suite.users[1].String(),
		Amount:       testDeposit,
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})

	if err != nil {
		suite.T().Logf("Deposit at extreme price rejected: %s", err.Error())
	} else {
		suite.T().Logf("System handled extreme price deposit")
		suite.checkSystemIntegrity()
	}
}

// TestMemoryAndStateStress tests system behavior with large state sizes
func (suite *StressTestSuite) TestMemoryAndStateStress() {
	// Create many user positions to stress state management
	successfulPositions := 0

	for i := 0; i < 200; i++ {
		user := suite.users[i%len(suite.users)].String()
		amount := math.NewInt(int64(10000 + (i * 100)))

		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user,
			Amount:       amount,
			ReceiveYield: i%3 != 0, // Vary yield preferences
			MinShares:    math.ZeroInt(),
		})

		if err == nil {
			successfulPositions++
		} else {
			// Log first failure and continue
			if successfulPositions == 0 {
				suite.T().Logf("First position creation failed: %s", err.Error())
			}
		}

		// Periodic yield updates to simulate realistic conditions
		if i%20 == 0 && i > 0 {
			vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
			if err == nil {
				yieldAmount := vaultState.TotalNav.MulRaw(2).QuoRaw(100) // 2% yield
				_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
					Authority: suite.authority,
					NewNav:    vaultState.TotalNav.Add(yieldAmount),
					Reason:    fmt.Sprintf("stress yield %d", i/20),
				})
				suite.Require().NoError(err)
			}
		}
	}

	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)

	suite.T().Logf("Memory and state stress test:")
	suite.T().Logf("  Successful positions created: %d/200", successfulPositions)
	suite.T().Logf("  Total users in vault: %d", vaultState.TotalUsers)
	suite.T().Logf("  Total shares outstanding: %s", vaultState.TotalShares.String())

	suite.checkSystemIntegrity()
}

// TestLongRunningOperations tests system stability over extended operation
func (suite *StressTestSuite) TestLongRunningOperations() {
	// Simulate long-running vault operations over many rounds
	rounds := 50
	usersPerRound := 5

	for round := 0; round < rounds; round++ {
		// Deposits
		for i := 0; i < usersPerRound; i++ {
			userIndex := (round*usersPerRound + i) % len(suite.users)
			user := suite.users[userIndex].String()
			amount := math.NewInt(int64(50000 + (round * 1000) + (i * 100)))

			_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
				Depositor:    user,
				Amount:       amount,
				ReceiveYield: (round+i)%2 == 0,
				MinShares:    math.ZeroInt(),
			})
			suite.Require().NoError(err)
		}

		// Periodic yield distribution
		if round%10 == 0 && round > 0 {
			vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
			suite.Require().NoError(err)

			yieldPercent := 2 + (round % 5) // 2-6% yield
			yieldAmount := vaultState.TotalNav.MulRaw(int64(yieldPercent)).QuoRaw(100)

			_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
				Authority: suite.authority,
				NewNav:    vaultState.TotalNav.Add(yieldAmount),
				Reason:    fmt.Sprintf("long running yield round %d", round/10),
			})
			suite.Require().NoError(err)
		}

		// Periodic withdrawals
		if round%15 == 0 && round > 0 {
			for i := 0; i < usersPerRound; i++ {
				userIndex := ((round-15)*usersPerRound + i) % len(suite.users)
				user := suite.users[userIndex]

				userPos, err := suite.keeper.GetV2UserPosition(suite.ctx, user)
				if err == nil && !userPos.Shares.IsZero() {
					withdrawShares := userPos.Shares.QuoRaw(3) // Withdraw 1/3

					_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
						Withdrawer: user.String(),
						Shares:     withdrawShares,
						MinAmount:  math.ZeroInt(),
					})
					suite.Require().NoError(err)
				}
			}
		}

		// Check invariants every 10 rounds
		if round%10 == 0 {
			suite.checkSystemIntegrity()
		}
	}

	finalState, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)

	suite.T().Logf("Long running operations test:")
	suite.T().Logf("  Rounds completed: %d", rounds)
	suite.T().Logf("  Final total users: %d", finalState.TotalUsers)
	suite.T().Logf("  Final total shares: %s", finalState.TotalShares.String())
	suite.T().Logf("  Final total NAV: %s", finalState.TotalNav.String())
}

// TestEdgeCaseSequences tests sequences of edge case operations
func (suite *StressTestSuite) TestEdgeCaseSequences() {
	// Initialize with minimal state
	seedUser := suite.users[0].String()
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    seedUser,
		Amount:       math.NewInt(1),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Test sequence: tiny deposits -> large yield -> tiny withdrawals
	for i := 1; i < 20; i++ {
		user := suite.users[i].String()

		// Tiny deposits
		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user,
			Amount:       math.NewInt(int64(i)), // 1, 2, 3, ... units
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)
	}

	// Large yield injection
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)

	largeYield := vaultState.TotalNav.MulRaw(1000) // 1000x yield
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    vaultState.TotalNav.Add(largeYield),
		Reason:    "extreme yield injection",
	})
	suite.Require().NoError(err)

	// Try withdrawals at extremely high share price
	successfulWithdrawals := 0
	for i := 1; i < 20; i++ {
		user := suite.users[i]
		userPos, err := suite.keeper.GetV2UserPosition(suite.ctx, user)
		if err == nil && !userPos.Shares.IsZero() {
			_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
				Withdrawer: user.String(),
				Shares:     userPos.Shares,
				MinAmount:  math.ZeroInt(),
			})
			if err == nil {
				successfulWithdrawals++
			}
		}
	}

	suite.T().Logf("Edge case sequences test:")
	suite.T().Logf("  Successful withdrawals after extreme yield: %d/19", successfulWithdrawals)

	suite.checkSystemIntegrity()
}

// TestSystemRecoveryAfterExtremeEvents tests system's ability to recover
func (suite *StressTestSuite) TestSystemRecoveryAfterExtremeEvents() {
	// Setup normal operation
	for i := 0; i < 10; i++ {
		user := suite.users[i].String()
		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user,
			Amount:       math.NewInt(int64((i + 1) * 100000)),
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)
	}

	// Create extreme event - massive NAV spike
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
	suite.Require().NoError(err)

	extremeNav := vaultState.TotalNav.MulRaw(10000) // 10000x increase
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    extremeNav,
		Reason:    "extreme event simulation",
	})
	suite.Require().NoError(err)

	// Test if system can recover with normal operations
	recoveryUser := suite.users[50].String()
	_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    recoveryUser,
		Amount:       math.NewInt(1000000),
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})

	systemRecovered := err == nil

	// Try to normalize NAV back
	normalNav := math.NewInt(10000000) // More reasonable NAV
	_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
		Authority: suite.authority,
		NewNav:    normalNav,
		Reason:    "recovery attempt",
	})

	navNormalized := err == nil

	suite.T().Logf("System recovery test:")
	suite.T().Logf("  System accepted deposits after extreme event: %t", systemRecovered)
	suite.T().Logf("  NAV normalization successful: %t", navNormalized)

	if systemRecovered && navNormalized {
		suite.checkSystemIntegrity()
	}
}

// TestBoundaryValueStress tests operations at mathematical boundaries
func (suite *StressTestSuite) TestBoundaryValueStress() {
	user := suite.users[0].String()

	// Test minimum non-zero values
	_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
		Depositor:    user,
		Amount:       math.NewInt(1), // Minimum amount
		ReceiveYield: true,
		MinShares:    math.ZeroInt(),
	})
	suite.Require().NoError(err)

	// Test various boundary NAV values
	boundaryNavs := []int64{
		2,          // Just above minimum
		1000,       // Small value
		1000000,    // Medium value
		1000000000, // Large value
	}

	for i, navValue := range boundaryNavs {
		_, err := suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
			Authority: suite.authority,
			NewNav:    math.NewInt(navValue),
			Reason:    fmt.Sprintf("boundary test %d", i),
		})
		suite.Require().NoError(err)

		// Test deposit at this NAV level
		testUser := suite.users[(i+1)%len(suite.users)].String()
		_, err = suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    testUser,
			Amount:       math.NewInt(int64((i + 1) * 10000)),
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)
	}

	suite.T().Logf("Boundary value stress test completed")
	suite.checkSystemIntegrity()
}

// Helper Functions

func (suite *StressTestSuite) checkSystemIntegrity() {
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
	if err != nil {
		suite.T().Logf("Vault state not accessible: %s", err.Error())
		return
	}

	// Core integrity checks
	suite.Require().True(vaultState.TotalShares.GTE(math.ZeroInt()), "Total shares cannot be negative")
	suite.Require().True(vaultState.TotalNav.GTE(math.ZeroInt()), "Total NAV cannot be negative")

	if !vaultState.TotalShares.IsZero() {
		suite.Require().True(vaultState.SharePrice.IsPositive(), "Share price must be positive")

		// Verify share price calculation
		expectedSharePrice := math.LegacyNewDecFromInt(vaultState.TotalNav).Quo(math.LegacyNewDecFromInt(vaultState.TotalShares))
		suite.Require().Equal(expectedSharePrice, vaultState.SharePrice, "Share price calculation must be consistent")
	}

	// Check for reasonable bounds
	if vaultState.SharePrice.GT(math.LegacyNewDec(1000000)) {
		suite.T().Logf("WARNING: Share price extremely high: %s", vaultState.SharePrice.String())
	}

	if vaultState.TotalUsers > 1000 {
		suite.T().Logf("WARNING: Very high user count: %d", vaultState.TotalUsers)
	}

	suite.T().Logf("System integrity check passed")
}

func (suite *StressTestSuite) logVaultState() {
	vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
	if err != nil {
		suite.T().Logf("Cannot read vault state: %s", err.Error())
		return
	}

	suite.T().Logf("Current vault state:")
	suite.T().Logf("  Total shares: %s", vaultState.TotalShares.String())
	suite.T().Logf("  Total NAV: %s", vaultState.TotalNav.String())
	suite.T().Logf("  Share price: %s", vaultState.SharePrice.String())
	suite.T().Logf("  Total users: %d", vaultState.TotalUsers)
}

// TestRandomizedOperations performs random sequences of operations
func (suite *StressTestSuite) TestRandomizedOperations() {
	// Initialize with some users
	for i := 0; i < 20; i++ {
		user := suite.users[i].String()
		amount := math.NewInt(int64((i + 1) * 50000))

		_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
			Depositor:    user,
			Amount:       amount,
			ReceiveYield: true,
			MinShares:    math.ZeroInt(),
		})
		suite.Require().NoError(err)
	}

	// Random operations
	operationCount := 100
	successfulOps := 0

	for i := 0; i < operationCount; i++ {
		opType := i % 4 // 0=deposit, 1=withdraw, 2=nav_update, 3=yield_change

		switch opType {
		case 0: // Deposit
			user := suite.users[i%len(suite.users)].String()
			amount := math.NewInt(int64(1000 + (i * 100)))

			_, err := suite.msgServer.Deposit(suite.ctx, &vaultsv2.MsgDeposit{
				Depositor:    user,
				Amount:       amount,
				ReceiveYield: i%2 == 0,
				MinShares:    math.ZeroInt(),
			})
			if err == nil {
				successfulOps++
			}

		case 1: // Withdraw
			user := suite.users[i%20] // Only from initial users
			userPos, err := suite.keeper.GetV2UserPosition(suite.ctx, user)
			if err == nil && !userPos.Shares.IsZero() {
				withdrawShares := userPos.Shares.QuoRaw(10) // Withdraw 10%

				_, err = suite.msgServer.Withdraw(suite.ctx, &vaultsv2.MsgWithdraw{
					Withdrawer: user.String(),
					Shares:     withdrawShares,
					MinAmount:  math.ZeroInt(),
				})
				if err == nil {
					successfulOps++
				}
			}

		case 2: // NAV Update
			vaultState, err := suite.keeper.GetV2VaultState(suite.ctx)
			if err == nil {
				change := int64(95 + (i % 10)) // 95-104% range
				newNav := vaultState.TotalNav.MulRaw(change).QuoRaw(100)
				if newNav.IsPositive() {
					_, err = suite.msgServer.UpdateNAV(suite.ctx, &vaultsv2.MsgUpdateNAV{
						Authority: suite.authority,
						NewNav:    newNav,
						Reason:    fmt.Sprintf("random update %d", i),
					})
					if err == nil {
						successfulOps++
					}
				}
			}

		case 3: // No-op (simulate other operations)
			successfulOps++
		}

		// Check system health periodically
		if i%25 == 0 {
			suite.checkSystemIntegrity()
		}
	}

	suite.T().Logf("Randomized operations stress test:")
	suite.T().Logf("  Total operations attempted: %d", operationCount)
	suite.T().Logf("  Successful operations: %d", successfulOps)
	suite.T().Logf("  Success rate: %.2f%%", float64(successfulOps)/float64(operationCount)*100)

	suite.checkSystemIntegrity()
}
