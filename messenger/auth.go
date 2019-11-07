package messenger

import (
	"context"

	mixin_sdk "github.com/fox-one/mixin-sdk"
)

// AuthorizeToken return access token and scope by authorizationCode
func AuthorizeToken(ctx context.Context, clientId, clientSecret string, authorizationCode string, codeVerifier string) (string, string, error) {
	params := map[string]interface{}{
		"client_id":     clientId,
		"client_secret": clientSecret,
		"code":          authorizationCode,
		"code_verifier": codeVerifier,
	}
	resp,err := mixin_sdk.Request(ctx).SetBody(params).Post("/oauth/token")
	if err != nil {
		return "","",err
	}

	var body struct{
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
	}

	err = mixin_sdk.UnmarshalResponse(resp,&body)
	return body.AccessToken,body.Scope,err
}
