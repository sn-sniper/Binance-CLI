package exchange

// Asset represents an exchange asset with raw balance strings
type Asset struct {
	Name   string
	Free   string
	Locked string
}

// ValuedAsset represents a wallet holding with real-time market valuation
type ValuedAsset struct {
	Name      string
	Free      string
	Locked    string
	Total     float64
	PriceUSDT float64
	ValueUSDT float64
}

// PriceStats holds 24-hour ticker statistics for a symbol
type PriceStats struct {
	Symbol             string
	BaseAsset          string
	LastPrice          string
	PriceChange        string
	PriceChangePercent string
	HighPrice          string
	LowPrice           string
	Volume             string
	QuoteVolume        string
}

// SymbolOverview represents a symbol with its base asset name and latest price
type SymbolOverview struct {
	Name   string
	Symbol string
	Price  string
}
