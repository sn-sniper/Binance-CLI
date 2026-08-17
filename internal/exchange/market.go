package exchange

import (
	"context"
	"fmt"
	"strings"
)

// Known quote currencies on Binance
var knownQuoteCurrencies = []string{
	"USDT", "USDC", "FDUSD", "TUSD", "BUSD",
	"BTC", "ETH", "BNB", "EUR", "TRY", "GBP", "AUD", "BRL", "RUB",
}

// NormalizeSymbol formats user input into standard Binance symbol format (e.g. "btc" -> "BTCUSDT", "ETH/USDT" -> "ETHUSDT")
func NormalizeSymbol(input string) string {
	cleaned := strings.ToUpper(strings.TrimSpace(input))
	cleaned = strings.ReplaceAll(cleaned, "/", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "_", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")

	if cleaned == "" {
		return ""
	}

	// Check if it already ends with a known quote currency
	for _, quote := range knownQuoteCurrencies {
		if strings.HasSuffix(cleaned, quote) && len(cleaned) > len(quote) {
			return cleaned
		}
	}

	// Default to USDT pair if no known quote currency is found
	return cleaned + "USDT"
}

// GetPriceStats fetches 24-hour ticker statistics for the specified symbols
func (s *Service) GetPriceStats(ctx context.Context, symbols ...string) ([]PriceStats, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("no symbols specified")
	}

	var results []PriceStats

	for _, rawSymbol := range symbols {
		normSymbol := NormalizeSymbol(rawSymbol)
		if normSymbol == "" {
			continue
		}

		stats, err := s.client.NewListPriceChangeStatsService().Symbol(normSymbol).Do(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch price for %s (%s): %w", rawSymbol, normSymbol, err)
		}

		for _, item := range stats {
			base := strings.TrimSuffix(item.Symbol, "USDT")
			results = append(results, PriceStats{
				Symbol:             item.Symbol,
				BaseAsset:          base,
				LastPrice:          item.LastPrice,
				PriceChange:        item.PriceChange,
				PriceChangePercent: item.PriceChangePercent,
				HighPrice:          item.HighPrice,
				LowPrice:           item.LowPrice,
				Volume:             item.Volume,
				QuoteVolume:        item.QuoteVolume,
			})
		}
	}

	return results, nil
}

// GetPriceStats is a package-level helper that uses DefaultService
func GetPriceStats(symbols ...string) ([]PriceStats, error) {
	return DefaultService().GetPriceStats(context.Background(), symbols...)
}
