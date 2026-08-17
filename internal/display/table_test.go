package display

import (
	"bytes"
	"testing"

	"example.com/binance/internal/exchange"
)

func TestFormatPrice(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{64300.50, "$64300.50"},
		{1.00, "$1.00"},
		{0.002136, "$0.002136"},
		{0.00, "$0.00"},
	}

	for _, tt := range tests {
		res := FormatPrice(tt.input)
		if res != tt.expected {
			t.Errorf("FormatPrice(%f) = %q, want %q", tt.input, res, tt.expected)
		}
	}
}

func TestFormatChangePercent(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"2.50", "+2.50%"},
		{"-1.25", "-1.25%"},
		{"0.00", "0.00%"},
	}

	for _, tt := range tests {
		res := FormatChangePercent(tt.input)
		if res != tt.expected {
			t.Errorf("FormatChangePercent(%q) = %q, want %q", tt.input, res, tt.expected)
		}
	}
}

func TestRenderTables(t *testing.T) {
	buf := new(bytes.Buffer)

	// Test RenderPricesTable
	stats := []exchange.PriceStats{
		{
			Symbol:             "BTCUSDT",
			BaseAsset:          "BTC",
			LastPrice:          "64000.00",
			PriceChangePercent: "2.5",
			HighPrice:          "65000.00",
			LowPrice:           "63000.00",
			Volume:             "1000.00",
		},
	}
	RenderPricesTable(buf, stats)
	if !bytes.Contains(buf.Bytes(), []byte("BTCUSDT")) {
		t.Errorf("expected prices table to contain BTCUSDT")
	}

	// Test RenderSymbolsTable
	buf.Reset()
	symbols := []exchange.SymbolOverview{
		{Name: "BTC", Symbol: "BTCUSDT", Price: "64000.00"},
	}
	RenderSymbolsTable(buf, symbols, "USDT")
	if !bytes.Contains(buf.Bytes(), []byte("BTCUSDT")) {
		t.Errorf("expected symbols table to contain BTCUSDT")
	}
}
