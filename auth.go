package mixin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-resty/resty/v2"
)

type Authentication interface {
	Auth(r *http.Request) (string, error)
	VerifyResponse(r *resty.Response) error
}

func RequestSig(method, uri string, body []byte) string {
	sum := sha256.Sum256(append([]byte(method+uri), body...))
	return hex.EncodeToString(sum[:])
}

func (token accessToken) Auth(r *http.Request) (string, error) {
	return string(token), nil
}

func (token accessToken) VerifyResponse(r *resty.Response) error {
	return nil
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
	return user.SignToken(RequestSig(r.Method, uri, body), r.Header.Get(requestIDHeaderKey), time.Minute)
}

func (user *User) SignToken(sig, reqID string, exp time.Duration) (string, error) {
	jwtMap := jwt.MapClaims{
		"uid": user.UserID,
		"sid": user.SessionID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(exp).Unix(),
		"jti": reqID,
		"sig": sig,
		"scp": "FULL",
	}
	if user.scopes != "" {
		jwtMap["scp"] = user.scopes
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwtMap)
	return token.SignedString(user.privateKey)
}

func (user *User) VerifyResponse(r *resty.Response) error {
	return nil
}

func (ed *EdOToken) Auth(r *http.Request) (string, error) {
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
	return ed.SignToken(RequestSig(r.Method, uri, body), r.Header.Get(requestIDHeaderKey), time.Minute)
}

func (ed *EdOToken) SignToken(sig, reqID string, exp time.Duration) (string, error) {
	jwtMap := jwt.MapClaims{
		"iss": ed.ClientID,
		"aid": ed.AuthID,
		"scp": ed.Scope,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(exp).Unix(),
		"sig": sig,
		"jti": reqID,
	}
	return jwt.NewWithClaims(Ed25519SigningMethod, jwtMap).
		SignedString(ed.EdPrivateKey)
}

func (ed *EdOToken) VerifyResponse(r *resty.Response) error {
	sign := r.Header().Get(integrityTokenHeaderKey)
	if sign == "" && IsErrorCodes(UnmarshalResponse(r, nil), 401) {
		return nil
	}

	token, err := jwt.Parse(sign, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*EdDSASigningMethod); !ok {
			return nil, jwt.ErrInvalidKeyType
		}
		return ed.EdServerPublicKey, nil
	})
	if err != nil {
		return err
	}

	claim := token.Claims.(jwt.MapClaims)

	{
		expect := r.Header().Get(requestIDHeaderKey)
		got, ok := claim["jti"].(string)
		if !ok || got != expect {
			return fmt.Errorf("token.jti mismatch, expect %q but got %q", expect, got)
		}
	}

	{
		method := r.Request.Method
		url := r.Request.RawRequest.URL
		idx := strings.Index(url.String(), url.Path)
		uri := url.String()[idx:]

		sum := sha256.Sum256(append([]byte(method+uri), r.Body()...))
		expect := hex.EncodeToString(sum[:])
		got, ok := claim["sig"].(string)
		if !ok || got != expect {
			return fmt.Errorf("token.sig mismatch, expect %q but got %q", expect, got)
		}
	}
	return nil
}
