package exchange

import (
	"testing"

	"example.com/binance/internal/config"
	"github.com/adshao/go-binance/v2"
)

func TestNewService(t *testing.T) {
	cfg := &config.Config{
		APIKey:    "test_key",
		SecretKey: "test_secret",
	}
	svc := NewService(cfg)
	if svc == nil {
		t.Fatalf("expected NewService to return a non-nil Service instance")
	}
}

func TestNormalizeSymbol(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"btc", "BTCUSDT"},
		{"BTC", "BTCUSDT"},
		{"BTCUSDT", "BTCUSDT"},
		{"btc/usdt", "BTCUSDT"},
		{"ETH-USDT", "ETHUSDT"},
		{"eth_btc", "ETHBTC"},
		{"sol", "SOLUSDT"},
		{"linea", "LINEAUSDT"},
		{"bnb", "BNBUSDT"},
		{"btceur", "BTCEUR"},
		{"", ""},
	}

	for _, tt := range tests {
		result := NormalizeSymbol(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeSymbol(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFilterNonZeroBalances(t *testing.T) {
	balances := []binance.Balance{
		{Asset: "BTC", Free: "0.00000000", Locked: "0.00000000"},
		{Asset: "ETH", Free: "1.50000000", Locked: "0.00000000"},
		{Asset: "SOL", Free: "0.00000000", Locked: "10.00000000"},
		{Asset: "USDT", Free: "250.75000000", Locked: "50.00000000"},
		{Asset: "BNB", Free: "invalid", Locked: "0.00000000"},
	}

	filtered := FilterNonZeroBalances(balances)

	if len(filtered) != 3 {
		t.Fatalf("expected 3 non-zero assets, got %d", len(filtered))
	}

	expectedAssets := map[string]struct {
		Free   string
		Locked string
	}{
		"ETH":  {Free: "1.50000000", Locked: "0.00000000"},
		"SOL":  {Free: "0.00000000", Locked: "10.00000000"},
		"USDT": {Free: "250.75000000", Locked: "50.00000000"},
	}

	for _, asset := range filtered {
		expected, ok := expectedAssets[asset.Name]
		if !ok {
			t.Errorf("unexpected asset in filtered list: %s", asset.Name)
			continue
		}
		if asset.Free != expected.Free || asset.Locked != expected.Locked {
			t.Errorf("asset %s = (%s, %s), want (%s, %s)", asset.Name, asset.Free, asset.Locked, expected.Free, expected.Locked)
		}
	}
}
