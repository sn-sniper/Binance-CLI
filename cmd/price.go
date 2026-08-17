package cmd

import (
	"fmt"
	"os"
	"strings"

	"example.com/binance/internal/display"
	"example.com/binance/internal/exchange"
	"github.com/spf13/cobra"
)

var priceCmd = &cobra.Command{
	Use:     "price [symbol...]",
	Aliases: []string{"ticker", "p"},
	Short:   "Check real-time spot price and 24h market stats for trading symbols",
	Long: `Check real-time spot price and 24-hour market statistics for one or more symbols.

You can specify full trading pairs (e.g. BTCUSDT, ETHBTC) or single coin names
(e.g. BTC, ETH, SOL, LINEA), which will automatically default to the USDT pair.`,
	Example: `  # Check price of Bitcoin (defaults to BTCUSDT)
  binance-cli price BTC

  # Check specific trading pair
  binance-cli price ETHBTC

  # Check multiple coins simultaneously
  binance-cli price BTC ETH SOL LINEA BNB

  # Using ticker alias
  binance-cli ticker SOLUSDT`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		symbolsToQuery := args
		if len(symbolsToQuery) == 0 {
			symbolsToQuery = []string{"BTC", "ETH", "SOL", "BNB"}
		}

		fmt.Printf("Fetching market price data for: %s...\n", strings.Join(symbolsToQuery, ", "))

		statsList, err := exchange.GetPriceStats(symbolsToQuery...)
		if err != nil {
			return fmt.Errorf("error fetching price: %w", err)
		}

		if len(statsList) == 0 {
			fmt.Println("No market data found for the requested symbol(s).")
			return nil
		}

		display.RenderPricesTable(os.Stdout, statsList)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(priceCmd)
}
