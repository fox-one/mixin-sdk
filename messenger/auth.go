package messenger

import (
	"context"
	"encoding/json"

	"github.com/fox-one/mixin-sdk/mixin"
	"github.com/fox-one/mixin-sdk/utils"
)

// AuthorizeToken return access token and scope by authorizationCode
func AuthorizeToken(ctx context.Context, clientId, clientSecret string, authorizationCode string, codeVerifier string) (string, string, error) {
	params, _ := json.Marshal(map[string]string{
		"client_id":     clientId,
		"client_secret": clientSecret,
		"code":          authorizationCode,
		"code_verifier": codeVerifier,
	})

	result := utils.SendRequest(ctx, "/oauth/token", "POST", string(params), "Content-Type", "application/json")
	data,err := result.Bytes()
	if err != nil {
		return "","",err
	}

	var resp struct {
		Data struct {
			AccessToken string `json:"access_token"`
			Scope       string `json:"scope"`
		} `json:"data"`
		Error *mixin.Error `json:"error,omitempty"`
	}

	if err = json.Unmarshal(data, &resp); err != nil {
		return "","", requestError(err)
	} else if resp.Error != nil {
		return "","", resp.Error
	}

	return resp.Data.AccessToken,resp.Data.Scope,nil
}
