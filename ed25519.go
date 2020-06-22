package mixin

import "crypto/ed25519"

type EdKey struct {
	EdPrivateKey      ed25519.PrivateKey `json:"ed_priv_key"`
	EdServerPublicKey ed25519.PublicKey  `json:"ed_server_pub_key"`
}

func NewEdKey(privSeed, serverPub []byte) *EdKey {
	return &EdKey{
		EdPrivateKey:      ed25519.NewKeyFromSeed(privSeed),
		EdServerPublicKey: ed25519.PublicKey(serverPub),
	}
}
