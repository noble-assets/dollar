package crosschain

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dollarv2 "dollar.noble.xyz/v2/types/v2"
	vaultsv2 "dollar.noble.xyz/v2/types/vaults/v2"
)

func TestNewHyperlaneProvider(t *testing.T) {
	// Create test mailbox ID
	mailboxId, err := util.DecodeHexAddress("0x1234567890123456789012345678901234567890123456789012345678901234")
	require.NoError(t, err)

	provider := NewHyperlaneProvider(
		nil, // coreKeeper - would be mocked in full tests
		nil, // warpKeeper - would be mocked in full tests
		4,   // Noble domain ID
		200000,
		time.Hour,
		mailboxId,
	)

	assert.NotNil(t, provider)
	assert.Equal(t, dollarv2.Provider_HYPERLANE, provider.GetProviderType())
	assert.Equal(t, uint32(4), provider.localDomain)
	assert.Equal(t, uint64(200000), provider.defaultGasLimit)
	assert.Equal(t, mailboxId, provider.mailboxId)
}

func TestValidateConfig(t *testing.T) {
	provider := &HyperlaneProvider{}

	tests := []struct {
		name        string
		config      *vaultsv2.CrossChainProviderConfig
		expectError bool
	}{
		{
			name: "valid config",
			config: &vaultsv2.CrossChainProviderConfig{
				Config: &vaultsv2.CrossChainProviderConfig_HyperlaneConfig{
					HyperlaneConfig: &vaultsv2.HyperlaneConfig{
						DomainId:       1,
						MailboxAddress: "0x1234567890123456789012345678901234567890123456789012345678901234",
						GasLimit:       200000,
						GasPrice:       math.NewInt(20000000000),
					},
				},
			},
			expectError: false,
		},
		{
			name: "missing hyperlane config",
			config: &vaultsv2.CrossChainProviderConfig{
				Config: &vaultsv2.CrossChainProviderConfig_IbcConfig{},
			},
			expectError: true,
		},
		{
			name: "zero domain ID",
			config: &vaultsv2.CrossChainProviderConfig{
				Config: &vaultsv2.CrossChainProviderConfig_HyperlaneConfig{
					HyperlaneConfig: &vaultsv2.HyperlaneConfig{
						DomainId:       0,
						MailboxAddress: "0x1234567890123456789012345678901234567890123456789012345678901234",
					},
				},
			},
			expectError: true,
		},
		{
			name: "empty mailbox address",
			config: &vaultsv2.CrossChainProviderConfig{
				Config: &vaultsv2.CrossChainProviderConfig_HyperlaneConfig{
					HyperlaneConfig: &vaultsv2.HyperlaneConfig{
						DomainId:       1,
						MailboxAddress: "",
					},
				},
			},
			expectError: true,
		},
		{
			name: "invalid mailbox address",
			config: &vaultsv2.CrossChainProviderConfig{
				Config: &vaultsv2.CrossChainProviderConfig_HyperlaneConfig{
					HyperlaneConfig: &vaultsv2.HyperlaneConfig{
						DomainId:       1,
						MailboxAddress: "invalid-address",
					},
				},
			},
			expectError: true,
		},
		{
			name: "gas limit too low",
			config: &vaultsv2.CrossChainProviderConfig{
				Config: &vaultsv2.CrossChainProviderConfig_HyperlaneConfig{
					HyperlaneConfig: &vaultsv2.HyperlaneConfig{
						DomainId:       1,
						MailboxAddress: "0x1234567890123456789012345678901234567890123456789012345678901234",
						GasLimit:       10000, // Too low
					},
				},
			},
			expectError: true,
		},
		{
			name: "negative gas price",
			config: &vaultsv2.CrossChainProviderConfig{
				Config: &vaultsv2.CrossChainProviderConfig_HyperlaneConfig{
					HyperlaneConfig: &vaultsv2.HyperlaneConfig{
						DomainId:       1,
						MailboxAddress: "0x1234567890123456789012345678901234567890123456789012345678901234",
						GasPrice:       math.NewInt(-1),
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := provider.ValidateConfig(tt.config)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEstimateGas(t *testing.T) {
	provider := &HyperlaneProvider{
		defaultGasLimit: 200000,
	}

	route := &vaultsv2.CrossChainRoute{
		ProviderConfig: &vaultsv2.CrossChainProviderConfig{
			Config: &vaultsv2.CrossChainProviderConfig_HyperlaneConfig{
				HyperlaneConfig: &vaultsv2.HyperlaneConfig{
					DomainId:       1,
					MailboxAddress: "0x1234567890123456789012345678901234567890123456789012345678901234",
					GasLimit:       300000,
					GasPrice:       math.NewInt(20000000000), // 20 gwei
				},
			},
		},
	}

	ctx := sdk.Context{} // Minimal context for testing

	tests := []struct {
		name        string
		messageType MessageType
		expectedGas uint64
		expectError bool
	}{
		{
			name:        "deposit message",
			messageType: MessageTypeDeposit,
			expectedGas: 300000, // Configured gas limit
		},
		{
			name:        "withdraw message",
			messageType: MessageTypeWithdraw,
			expectedGas: 300000, // Configured gas limit
		},
		{
			name:        "update message",
			messageType: MessageTypeUpdate,
			expectedGas: 300000, // Configured gas limit
		},
		{
			name:        "liquidate message",
			messageType: MessageTypeLiquidate,
			expectedGas: 300000, // Configured gas limit
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := CrossChainMessage{
				Type:      tt.messageType,
				Sender:    sdk.AccAddress("test"),
				Recipient: "0x1234567890123456789012345678901234567890",
				Amount:    math.NewInt(1000000),
			}

			gasLimit, totalCost, err := provider.EstimateGas(ctx, route, msg)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedGas, gasLimit)
				assert.True(t, totalCost.GT(math.ZeroInt()))
			}
		})
	}
}

func TestRequiredConfirmations(t *testing.T) {
	provider := &HyperlaneProvider{
		confirmationsMap: map[uint32]uint64{
			1:     12, // Ethereum
			137:   20, // Polygon
			42161: 1,  // Arbitrum
		},
	}

	tests := []struct {
		domain           uint32
		expectedConfirms uint64
	}{
		{1, 12},    // Ethereum
		{137, 20},  // Polygon
		{42161, 1}, // Arbitrum
		{999, 12},  // Unknown domain (default)
	}

	for _, tt := range tests {
		t.Run("domain_"+string(rune(tt.domain)), func(t *testing.T) {
			confirmations := provider.GetRequiredConfirmations(tt.domain)
			assert.Equal(t, tt.expectedConfirms, confirmations)
		})
	}
}

func TestSetRequiredConfirmations(t *testing.T) {
	provider := &HyperlaneProvider{
		confirmationsMap: make(map[uint32]uint64),
	}

	provider.SetRequiredConfirmations(1, 15)
	assert.Equal(t, uint64(15), provider.GetRequiredConfirmations(1))

	provider.SetRequiredConfirmations(1, 10)
	assert.Equal(t, uint64(10), provider.GetRequiredConfirmations(1))
}

func TestMailboxIdManagement(t *testing.T) {
	originalMailboxId, err := util.DecodeHexAddress("0x1234567890123456789012345678901234567890123456789012345678901234")
	require.NoError(t, err)

	provider := &HyperlaneProvider{
		mailboxId: originalMailboxId,
	}

	// Test getting mailbox ID
	assert.Equal(t, originalMailboxId, provider.GetMailboxId())

	// Test setting new mailbox ID
	newMailboxId, err := util.DecodeHexAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdabcdefabcdefabcdefabcdef")
	require.NoError(t, err)

	provider.SetMailboxId(newMailboxId)
	assert.Equal(t, newMailboxId, provider.GetMailboxId())
}

func TestMessageTypeToString(t *testing.T) {
	provider := &HyperlaneProvider{}

	tests := []struct {
		messageType MessageType
		expected    string
	}{
		{MessageTypeDeposit, "deposit"},
		{MessageTypeWithdraw, "withdraw"},
		{MessageTypeUpdate, "update"},
		{MessageTypeLiquidate, "liquidate"},
		{MessageType(999), "unknown"}, // Unknown type
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := provider.messageTypeToString(tt.messageType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateMessageBody(t *testing.T) {
	provider := &HyperlaneProvider{}

	tests := []struct {
		name     string
		msg      CrossChainMessage
		contains []string
	}{
		{
			name: "simple message",
			msg: CrossChainMessage{
				Type:      MessageTypeDeposit,
				Sender:    sdk.AccAddress("test"),
				Recipient: "0x1234567890123456789012345678901234567890",
				Amount:    math.NewInt(1000000),
			},
			contains: []string{"deposit", "cosmos1w3jhxaq8w2lx0", "0x1234567890123456789012345678901234567890", "1000000"},
		},
		{
			name: "message with data",
			msg: CrossChainMessage{
				Type:      MessageTypeWithdraw,
				Sender:    sdk.AccAddress("test"),
				Recipient: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
				Amount:    math.NewInt(500000),
				Data:      []byte("extra data"),
			},
			contains: []string{"withdraw", "cosmos1w3jhxaq8w2lx0", "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd", "500000"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := provider.createMessageBody(tt.msg)
			assert.NoError(t, err)
			assert.NotEmpty(t, body)

			bodyString := string(body)
			for _, expected := range tt.contains {
				assert.Contains(t, bodyString, expected)
			}
		})
	}
}

func TestUpdateMessageStatus(t *testing.T) {
	provider := &HyperlaneProvider{}

	tracking := &vaultsv2.ProviderTrackingInfo{
		TrackingInfo: &vaultsv2.ProviderTrackingInfo_HyperlaneTracking{
			HyperlaneTracking: &vaultsv2.HyperlaneTrackingInfo{
				MessageId:         []byte("test-message-id"),
				OriginDomain:      4,
				DestinationDomain: 1,
				Processed:         false,
			},
		},
	}

	ctx := sdk.Context{} // Minimal context for testing

	err := provider.UpdateMessageStatus(
		ctx,
		tracking,
		"0xdestinationTxHash",
		12345,
		67890,
		true,
	)

	assert.NoError(t, err)

	hyperlaneTracking := tracking.GetHyperlaneTracking()
	assert.Equal(t, "0xdestinationTxHash", hyperlaneTracking.DestinationTxHash)
	assert.Equal(t, uint64(12345), hyperlaneTracking.DestinationBlockNumber)
	assert.Equal(t, uint64(67890), hyperlaneTracking.GasUsed)
	assert.True(t, hyperlaneTracking.Processed)
}

// Benchmark tests for performance
func BenchmarkCreateMessageBody(b *testing.B) {
	provider := &HyperlaneProvider{}
	msg := CrossChainMessage{
		Type:      MessageTypeDeposit,
		Sender:    sdk.AccAddress("benchmark-test"),
		Recipient: "0x1234567890123456789012345678901234567890",
		Amount:    math.NewInt(1000000),
		Data:      make([]byte, 1024), // 1KB of data
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := provider.createMessageBody(msg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidateConfig(b *testing.B) {
	provider := &HyperlaneProvider{}
	config := &vaultsv2.CrossChainProviderConfig{
		Config: &vaultsv2.CrossChainProviderConfig_HyperlaneConfig{
			HyperlaneConfig: &vaultsv2.HyperlaneConfig{
				DomainId:       1,
				MailboxAddress: "0x1234567890123456789012345678901234567890123456789012345678901234",
				GasLimit:       200000,
				GasPrice:       math.NewInt(20000000000),
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := provider.ValidateConfig(config)
		if err != nil {
			b.Fatal(err)
		}
	}
}
