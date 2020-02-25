package mixin

import (
	"context"

	"github.com/shopspring/decimal"
)

type ExchangeRate struct {
	Code string          `json:"code,omitempty"`
	Rate decimal.Decimal `json:"rate,omitempty"`
}

func (user *User) ReadExchangeRates(ctx context.Context) ([]ExchangeRate, error) {
	var rates []ExchangeRate
	if err := user.Request(ctx, "GET", "/fiats", nil, &rates); err != nil {
		return nil, err
	}

	return rates, nil
}
