package util

const (
	USDT = "USDT"
	BTC  = "BTC"
	ETH  = "ETH"
)

func IsValidCurrency(currency string) bool {
	switch currency {
	case USDT, BTC, ETH:
		return true
	default:
		return false
	}
}
