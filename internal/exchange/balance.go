package exchange

import (
	"context"
	"strconv"

	"github.com/adshao/go-binance/v2"
)

// FilterNonZeroBalances extracts assets with free or locked balances > 0
func FilterNonZeroBalances(balances []binance.Balance) []Asset {
	var activeAssets []Asset
	for _, balance := range balances {
		free, errFree := strconv.ParseFloat(balance.Free, 64)
		locked, errLocked := strconv.ParseFloat(balance.Locked, 64)

		if (errFree == nil && free > 0) || (errLocked == nil && locked > 0) {
			activeAssets = append(activeAssets, Asset{
				Name:   balance.Asset,
				Free:   balance.Free,
				Locked: balance.Locked,
			})
		}
	}
	return activeAssets
}

// GetNonZeroBalances fetches the wallet and filters out empty coins
func (s *Service) GetNonZeroBalances(ctx context.Context) ([]Asset, error) {
	if err := s.cfg.Validate(); err != nil {
		return nil, err
	}

	res, err := s.client.NewGetAccountService().Do(ctx)
	if err != nil {
		return nil, err
	}

	return FilterNonZeroBalances(res.Balances), nil
}

// GetValuedBalances fetches wallet balances and calculates real-time USD/USDT valuation
func (s *Service) GetValuedBalances(ctx context.Context) ([]ValuedAsset, float64, error) {
	assets, err := s.GetNonZeroBalances(ctx)
	if err != nil {
		return nil, 0, err
	}

	if len(assets) == 0 {
		return nil, 0, nil
	}

	prices, err := s.client.NewListPricesService().Do(ctx)
	priceMap := make(map[string]float64)
	if err == nil {
		for _, p := range prices {
			if val, parseErr := strconv.ParseFloat(p.Price, 64); parseErr == nil {
				priceMap[p.Symbol] = val
			}
		}
	}

	var valued []ValuedAsset
	var totalPortfolioValue float64

	for _, a := range assets {
		free, _ := strconv.ParseFloat(a.Free, 64)
		locked, _ := strconv.ParseFloat(a.Locked, 64)
		totalAmount := free + locked

		var unitPrice float64
		switch a.Name {
		case "USDT", "USDC", "BUSD", "FDUSD", "TUSD", "USD":
			unitPrice = 1.0
		default:
			if p, ok := priceMap[a.Name+"USDT"]; ok {
				unitPrice = p
			} else if p, ok := priceMap[a.Name+"BUSD"]; ok {
				unitPrice = p
			} else if p, ok := priceMap[a.Name+"USDC"]; ok {
				unitPrice = p
			}
		}

		totalVal := totalAmount * unitPrice
		totalPortfolioValue += totalVal

		valued = append(valued, ValuedAsset{
			Name:      a.Name,
			Free:      a.Free,
			Locked:    a.Locked,
			Total:     totalAmount,
			PriceUSDT: unitPrice,
			ValueUSDT: totalVal,
		})
	}

	return valued, totalPortfolioValue, nil
}

// GetNonZeroBalances is a package-level helper that uses DefaultService
func GetNonZeroBalances() ([]Asset, error) {
	return DefaultService().GetNonZeroBalances(context.Background())
}

// GetValuedBalances is a package-level helper that uses DefaultService
func GetValuedBalances() ([]ValuedAsset, float64, error) {
	return DefaultService().GetValuedBalances(context.Background())
}
