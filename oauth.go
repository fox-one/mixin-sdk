package sdk

import (
	"context"
)

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
