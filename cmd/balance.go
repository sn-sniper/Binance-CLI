package cmd

import (
	"fmt"
	"os"

	"example.com/binance/internal/display"
	"example.com/binance/internal/exchange"
	"github.com/spf13/cobra"
)

var balanceCmd = &cobra.Command{
	Use:     "balance",
	Aliases: []string{"wallet", "bal", "b"},
	Short:   "View non-zero spot wallet balances and live USD portfolio value",
	Long: `Fetch and display non-zero spot wallet balances from Binance with real-time market valuation.

Connects to Binance using your API credentials (BINANCE_API_KEY & BINANCE_SECRET_KEY),
filters out empty coins, looks up current ticker prices, and renders a clean table
showing free, locked, total balances, and estimated USD value.`,
	Example: `  # Fetch wallet balances with live portfolio valuation
  binance-cli balance

  # Using shorthand alias
  binance-cli bal`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Fetching wallet balances and market valuations...")

		valuedAssets, totalValue, err := exchange.GetValuedBalances()
		if err != nil {
			return fmt.Errorf("error fetching balance: %w", err)
		}

		if len(valuedAssets) == 0 {
			fmt.Println("Your wallet is completely empty (not even $0.80!).")
			return nil
		}

		display.RenderBalancesTable(os.Stdout, valuedAssets, totalValue)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)
}