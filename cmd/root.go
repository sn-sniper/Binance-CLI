package cmd

import (
	"os"

	"example.com/binance/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "binance-cli [command]",
	Short: "A fast terminal dashboard and market analysis tool for Binance",
	Long: `binance-cli is a modular command-line dashboard and market tool for Binance.

Available Commands:
  balance     View non-zero spot wallet balances and live USD portfolio value
  price       Check real-time spot price and 24h market stats (24h high/low/vol)
  symbols     List active Binance trading symbols and prices (Name | Symbol | Price)
  completion  Generate shell autocompletion scripts (bash, zsh, fish, powershell)
  help        Help about any command

Configuration:
  Set your credentials in a .env file or system environment variables:
    BINANCE_API_KEY=your_api_key
    BINANCE_SECRET_KEY=your_secret_key`,
	Example: `  # 1. Check wallet balances and portfolio valuation:
  binance-cli balance

  # 2. Check spot price and 24h ticker:
  binance-cli price BTC
  binance-cli price BTC ETH SOL LINEA

  # 3. List symbols in a clean Name | Symbol | Price table:
  binance-cli symbols
  binance-cli symbols BTC --limit 10
  binance-cli symbols --quote BTC

  # 4. View help for any subcommand:
  binance-cli balance --help
  binance-cli price --help
  binance-cli symbols --help`,
}

// Execute adds all child commands to the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Initialize configuration globally
	_ = config.Load()
}