package display

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"example.com/binance/internal/exchange"
	"github.com/olekukonko/tablewriter"
)

// NewTable initializes a tablewriter Table with standardized styling
func NewTable(writer io.Writer, headers []string) *tablewriter.Table {
	table := tablewriter.NewWriter(writer)
	table.SetHeader(headers)
	table.SetBorder(true)
	table.SetRowLine(true)

	// Set cyan bold header colors for all columns
	headerColors := make([]tablewriter.Colors, len(headers))
	for i := range headerColors {
		headerColors[i] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor}
	}
	table.SetHeaderColor(headerColors...)

	return table
}

// FormatPrice formats a floating point price nicely with appropriate decimals and dollar prefix
func FormatPrice(price float64) string {
	if price >= 1.0 {
		return fmt.Sprintf("$%.2f", price)
	} else if price > 0 {
		return fmt.Sprintf("$%.6f", price)
	}
	return "$0.00"
}

// FormatChangePercent formats a price change percentage with a plus sign for positive numbers
func FormatChangePercent(changeStr string) string {
	val, err := strconv.ParseFloat(changeStr, 64)
	if err != nil {
		return changeStr
	}

	if val > 0 {
		return fmt.Sprintf("+%.2f%%", val)
	}
	return fmt.Sprintf("%.2f%%", val)
}

// RenderBalancesTable renders a formatted table of valued wallet balances
func RenderBalancesTable(writer io.Writer, assets []exchange.ValuedAsset, totalValue float64) {
	table := NewTable(writer, []string{
		"Asset", "Free Balance", "Locked Balance", "Total", "Price (USDT)", "Est. Value (USDT)",
	})

	for _, asset := range assets {
		var formattedPrice string
		if asset.PriceUSDT > 0 {
			formattedPrice = FormatPrice(asset.PriceUSDT)
		} else {
			formattedPrice = "N/A"
		}

		var formattedValue string
		if asset.ValueUSDT > 0 {
			formattedValue = fmt.Sprintf("$%.2f", asset.ValueUSDT)
		} else {
			formattedValue = "$0.00"
		}

		table.Append([]string{
			asset.Name,
			asset.Free,
			asset.Locked,
			fmt.Sprintf("%.8f", asset.Total),
			formattedPrice,
			formattedValue,
		})
	}

	table.SetFooter([]string{"", "", "", "", "Total Value:", fmt.Sprintf("$%.2f USDT", totalValue)})
	table.SetFooterColor(
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
	)

	table.Render()
}

// RenderPricesTable renders 24h ticker prices and market statistics
func RenderPricesTable(writer io.Writer, statsList []exchange.PriceStats) {
	table := NewTable(writer, []string{
		"Asset", "Symbol", "Price (USDT)", "24h Change", "24h High", "24h Low", "24h Volume",
	})

	for _, s := range statsList {
		priceVal, _ := strconv.ParseFloat(s.LastPrice, 64)
		highVal, _ := strconv.ParseFloat(s.HighPrice, 64)
		lowVal, _ := strconv.ParseFloat(s.LowPrice, 64)
		volVal, _ := strconv.ParseFloat(s.Volume, 64)

		table.Append([]string{
			s.BaseAsset,
			s.Symbol,
			FormatPrice(priceVal),
			FormatChangePercent(s.PriceChangePercent),
			fmt.Sprintf("$%.4f", highVal),
			fmt.Sprintf("$%.4f", lowVal),
			fmt.Sprintf("%.2f", volVal),
		})
	}

	table.Render()
}

// RenderSymbolsTable renders a list of trading symbols in Name | Symbol | Price format
func RenderSymbolsTable(writer io.Writer, symbols []exchange.SymbolOverview, quoteCurrency string) {
	table := NewTable(writer, []string{"Name", "Symbol", "Price"})

	isUSDT := strings.ToUpper(quoteCurrency) == "USDT" || quoteCurrency == ""

	for _, s := range symbols {
		priceVal, parseErr := strconv.ParseFloat(s.Price, 64)
		var formattedPrice string
		if parseErr == nil {
			if isUSDT || strings.HasSuffix(s.Symbol, "USDT") || strings.HasSuffix(s.Symbol, "USDC") {
				formattedPrice = FormatPrice(priceVal)
			} else {
				formattedPrice = s.Price
			}
		} else {
			formattedPrice = s.Price
		}

		table.Append([]string{s.Name, s.Symbol, formattedPrice})
	}

	table.Render()
}
