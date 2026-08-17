package exchange

import (
	"context"
	"fmt"
	"strings"
)

// ListSymbols fetches active symbols with their current prices, supporting quote filtering and search
func (s *Service) ListSymbols(ctx context.Context, quoteAsset, searchQuery string, limit int) ([]SymbolOverview, error) {
	// 1. Fetch exchange info for symbols and base/quote assets
	info, err := s.client.NewExchangeInfoService().Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange information: %w", err)
	}

	// 2. Fetch latest prices
	prices, err := s.client.NewListPricesService().Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ticker prices: %w", err)
	}

	priceMap := make(map[string]string, len(prices))
	for _, p := range prices {
		priceMap[p.Symbol] = p.Price
	}

	quoteAsset = strings.ToUpper(strings.TrimSpace(quoteAsset))
	searchQuery = strings.ToUpper(strings.TrimSpace(searchQuery))

	var list []SymbolOverview
	for _, sym := range info.Symbols {
		if sym.Status != "TRADING" {
			continue
		}

		if quoteAsset != "" && quoteAsset != "ALL" && sym.QuoteAsset != quoteAsset {
			continue
		}

		if searchQuery != "" {
			if !strings.Contains(sym.Symbol, searchQuery) && !strings.Contains(sym.BaseAsset, searchQuery) {
				continue
			}
		}

		price := priceMap[sym.Symbol]
		if price == "" {
			price = "0.00"
		}

		list = append(list, SymbolOverview{
			Name:   sym.BaseAsset,
			Symbol: sym.Symbol,
			Price:  price,
		})

		if limit > 0 && len(list) >= limit {
			break
		}
	}

	return list, nil
}

// ListSymbols is a package-level helper that uses DefaultService
func ListSymbols(quoteAsset, searchQuery string, limit int) ([]SymbolOverview, error) {
	return DefaultService().ListSymbols(context.Background(), quoteAsset, searchQuery, limit)
}
