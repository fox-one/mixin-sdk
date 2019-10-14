package mixin

// WithdrawAddress withdraw address
type WithdrawAddress struct {
	AddressID string `json:"address_id,omitempty"`
	AssetID   string `json:"asset_id"`

	Destination string `json:"destination,omitempty"`
	Tag         string `json:"tag,omitempty"`

	Fee  string `json:"fee,omitempty"`
	Dust string `json:"dust,omitempty"`

	// TODO Deprecated
	PublicKey   string `json:"public_key,omitempty"`
	Label       string `json:"label,omitempty"`
	AccountName string `json:"account_name,omitempty"`
	AccountTag  string `json:"account_tag,omitempty"`
}
