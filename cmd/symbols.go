package cmd

import (
	"fmt"
	"os"

	"example.com/binance/internal/display"
	"example.com/binance/internal/exchange"
	"github.com/spf13/cobra"
)

var (
	symbolsQuote  string
	symbolsSearch string
	symbolsLimit  int
)

var symbolsCmd = &cobra.Command{
	Use:     "symbols [search]",
	Aliases: []string{"list", "list-symbols", "pairs", "s"},
	Short:   "List active Binance trading symbols and latest prices (Name | Symbol | Price)",
	Long: `List active Binance trading symbols and their latest prices in a clean table.

Filter pairs by quote currency (e.g. USDT, BTC, EUR) or search for specific coins.`,
	Example: `  # List top 50 USDT trading pairs
  binance-cli symbols

  # Search for Bitcoin or Ethereum pairs
  binance-cli symbols BTC
  binance-cli symbols --search ETH

  # Filter by quote currency (e.g. BTC pairs)
  binance-cli symbols --quote BTC --limit 20

  # List all trading pairs across all quote assets
  binance-cli symbols --quote ALL --limit 100`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		search := symbolsSearch
		if len(args) > 0 && search == "" {
			search = args[0]
		}

		quoteDisplay := symbolsQuote
		if quoteDisplay == "" {
			quoteDisplay = "USDT"
		}

		fmt.Printf("Fetching symbol prices (quote: %s, search: %q, limit: %d)...\n", quoteDisplay, search, symbolsLimit)

		symbols, err := exchange.ListSymbols(symbolsQuote, search, symbolsLimit)
		if err != nil {
			return fmt.Errorf("error listing symbols: %w", err)
		}

		if len(symbols) == 0 {
			fmt.Println("No trading symbols found matching your criteria.")
			return nil
		}

		display.RenderSymbolsTable(os.Stdout, symbols, quoteDisplay)
		fmt.Printf("Showing %d symbol(s).\n", len(symbols))
		return nil
	},
}

func init() {
	symbolsCmd.Flags().StringVarP(&symbolsQuote, "quote", "q", "USDT", "Filter by quote currency (e.g. USDT, BTC, EUR, ALL)")
	symbolsCmd.Flags().StringVarP(&symbolsSearch, "search", "s", "", "Search query to filter symbol or coin name")
	symbolsCmd.Flags().IntVarP(&symbolsLimit, "limit", "l", 50, "Maximum number of symbols to display (0 for unlimited)")

	rootCmd.AddCommand(symbolsCmd)
}
