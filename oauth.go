package mixin

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"

	"github.com/gofrs/uuid"
)

// Token Auth
type (
	accessToken string

	EdOToken struct {
		*EdKey

		ClientID uuid.UUID `json:"client_id"`
		AuthID   uuid.UUID `json:"auth_id"`
		Scope    string    `json:"scope"`
	}
)

func NewEdOToken(clientID, authID uuid.UUID, privSeed, serverPub []byte, scope string) *EdOToken {
	ed := EdOToken{
		EdKey:    &EdKey{},
		ClientID: clientID,
		AuthID:   authID,
		Scope:    scope,
	}
	if len(privSeed) == ed25519.SeedSize {
		ed.EdPrivateKey = ed25519.NewKeyFromSeed(privSeed)
	}
	if len(serverPub) == ed25519.PublicKeySize {
		ed.EdServerPublicKey = ed25519.PublicKey(serverPub)
	}
	return &ed
}

func (ed *EdOToken) SetPrivateKeySeed(seed []byte) {
	ed.EdPrivateKey = ed25519.NewKeyFromSeed(seed)
}

// AuthorizeToken return access token and scope by authorizationCode
func AuthorizeToken(ctx context.Context, clientId, clientSecret string, authorizationCode string, codeVerifier string) (string, string, error) {
	params := map[string]interface{}{
		"client_id":     clientId,
		"client_secret": clientSecret,
		"code":          authorizationCode,
		"code_verifier": codeVerifier,
	}
	resp, err := Request(ctx).SetBody(params).Post("/oauth/token")
	if err != nil {
		return "", "", err
	}

	var body struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
	}

	err = UnmarshalResponse(resp, &body)
	return body.AccessToken, body.Scope, err
}

func AuthorizeTokenEd25519(ctx context.Context, clientID, secret string, code string, verifier string) (*EdOToken, error) {
	var seed = make([]byte, ed25519.SeedSize)
	rand.Read(seed)
	priv := ed25519.NewKeyFromSeed(seed)
	params := map[string]interface{}{
		"client_id":     clientID,
		"client_secret": secret,
		"code":          code,
		"code_verifier": verifier,
		"ed25519":       base64.RawURLEncoding.EncodeToString(priv[ed25519.SeedSize:]),
	}

	resp, err := Request(ctx).SetBody(params).Post("/oauth/token")
	if err != nil {
		return nil, err
	}

	var body struct {
		AuthID    uuid.UUID `json:"authorization_id"`
		PublicKey []byte    `json:"ed25519"`
		Scope     string    `json:"scope"`
	}

	if err = UnmarshalResponse(resp, &body); err != nil {
		return nil, err
	}

	cid, _ := uuid.FromString(clientID)
	ed := NewEdOToken(cid, body.AuthID, nil, body.PublicKey, body.Scope)
	ed.EdPrivateKey = priv
	return ed, err
}
