package mixin

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"io"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/gofrs/uuid"
)

// User wallet entity
type User struct {
	UserID     string `json:"user_id"`
	SessionID  string `json:"session_id"`
	PINToken   string `json:"pin_token"`
	SessionKey string `json:"session_key"`

	FullName string `json:"full_name"`

	pinCipher  *cipher.Block
	privateKey *rsa.PrivateKey
	scopes     string
}

// NewUser new user
func NewUser(userID, sessionID, sessionKey string) (*User, error) {
	user := &User{
		UserID:    userID,
		SessionID: sessionID,
	}

	block, _ := pem.Decode([]byte(sessionKey))
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	user.SetPrivateKey(privateKey)
	return user, nil
}

// SetPrivateKey set private key
func (user *User) SetPrivateKey(privateKey *rsa.PrivateKey) {
	user.privateKey = privateKey
}

// SetScopes set scopes
func (user *User) SetScopes(scopes string) {
	user.scopes = scopes
}

// HasPrivateKey private key has been set
func (user *User) HasPrivateKey() bool {
	return user.privateKey != nil
}

func (user *User) loadPINCipher() error {
	if user.pinCipher != nil {
		return nil
	}

	token, err := base64.StdEncoding.DecodeString(user.PINToken)
	if err != nil {
		return err
	}

	keyBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, user.privateKey, token, []byte(user.SessionID))
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return err
	}

	user.pinCipher = &block

	return nil
}

func (user *User) signPIN(pin string) (string, error) {
	if len(pin) == 0 {
		return "", nil
	}

	if user.pinCipher == nil {
		if err := user.loadPINCipher(); err != nil {
			return "", err
		}
	}
	pinByte := []byte(pin)
	timeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(timeBytes, uint64(time.Now().Unix()))
	pinByte = append(pinByte, timeBytes...)
	iteratorBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(iteratorBytes, uint64(time.Now().UnixNano()))
	pinByte = append(pinByte, iteratorBytes...)
	padding := aes.BlockSize - len(pinByte)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	pinByte = append(pinByte, padtext...)
	ciphertext := make([]byte, aes.BlockSize+len(pinByte))
	iv := ciphertext[:aes.BlockSize]
	io.ReadFull(rand.Reader, iv)
	mode := cipher.NewCBCEncrypter(*user.pinCipher, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], pinByte)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// SignToken sign request
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
	}
	if user.scopes != "" {
		jwtMap["scp"] = user.scopes
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwtMap)

	return token.SignedString(user.privateKey)
}
