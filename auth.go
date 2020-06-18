package mixin

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
)

type Authentication interface {
	Auth(r *http.Request) (string, error)
}

// Token Auth

type accessToken string

func (token accessToken) Auth(r *http.Request) (string, error) {
	return string(token), nil
}

// User Auth

// SignToken sign request
func (user *User) Auth(r *http.Request) (string, error) {
	url := r.URL.String()
	idx := strings.Index(url, r.URL.Path)
	uri := url[idx:]

	var body []byte
	if r.GetBody != nil {
		if rc, _ := r.GetBody(); rc != nil {
			defer rc.Close()
			body, _ = ioutil.ReadAll(rc)
		}
	}
	return user.SignToken(r.Method, uri, body, time.Minute)
}

func (user *User) SignToken(method, uri string, body []byte, expire ...time.Duration) (string, error) {
	e := time.Hour * 24 * 30 * 3
	if len(expire) > 0 && expire[0] > 5 {
		e = expire[0]
	}
	expireAt := time.Now().UTC().Add(e)
	sum := sha256.Sum256(append([]byte(method+uri), body...))

	jwtMap := jwt.MapClaims{
		"uid": user.UserID,
		"sid": user.SessionID,
		"iat": time.Now().UTC().Unix(),
		"exp": expireAt.Unix(),
		"jti": uuid.Must(uuid.NewV4()).String(),
		"sig": hex.EncodeToString(sum[:]),
		"scp": "FULL",
	}
	if user.scopes != "" {
		jwtMap["scp"] = user.scopes
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwtMap)
	return token.SignedString(user.privateKey)
}
