package currency

type Currency string

var (
	EUR Currency = "EUR"
	USD Currency = "USD"
	MXN Currency = "MXN"
)

var supportedCurrencies = map[Currency]struct{}{
	EUR: {},
	USD: {},
	MXN: {},
}

func IsValid(c Currency) bool {
	_, ok := supportedCurrencies[c]
	return ok
}
