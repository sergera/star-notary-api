package domain

import (
	"testing"
)

func TestValidateWallet(t *testing.T) {
	tests := []struct {
		name        string
		wallet      WalletModel
		expectedErr bool
	}{
		{"Valid Wallet Address", WalletModel{Address: "0x1234567890ABCDEFFEDCBA098765432123456789"}, false},
		{"Invalid Wallet Address - Wrong Length", WalletModel{Address: "0x12345"}, true},
		{"Invalid Wallet Address - Missing Prefix", WalletModel{Address: "1234567890ABCDEFFEDCBA098765432123456789"}, true},
		{"Invalid Wallet Address - Invalid Characters", WalletModel{Address: "0xGHIJKLMNOPQRSTUVWXYZ1234567890ABCDEFFED"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.wallet.Validate()
			if (err != nil) != tt.expectedErr {
				t.Errorf("TestValidateWallet(%s): expected error: %v, got error: %v", tt.name, tt.expectedErr, err != nil)
			}
		})
	}
}
