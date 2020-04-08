package mixin

import (
	"net/url"
)

func PayURL(assetID, traceID, opponentID, amount, memo string) string {
	u := &url.URL{
		Scheme: "mixin",
		Host:   "pay",
	}

	q := u.Query()
	q.Set("asset", assetID)
	q.Set("trace", traceID)
	q.Set("amount", amount)
	q.Set("recipient", opponentID)
	q.Set("memo", memo)
	u.RawQuery = q.Encode()

	return u.String()
}
