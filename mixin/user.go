package mixin

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"io"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
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
}

// SetPrivateKey set private key
func (user *User) SetPrivateKey(privateKey *rsa.PrivateKey) {
	user.privateKey = privateKey
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
func (user *User) SignToken(method, uri string, body []byte) (string, error) {
	expire := time.Now().UTC().Add(time.Hour * 24 * 30 * 3)
	sum := sha256.Sum256(append([]byte(method+uri), body...))

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"uid": user.UserID,
		"sid": user.SessionID,
		"iat": time.Now().UTC().Unix(),
		"exp": expire.Unix(),
		"jti": uuid.Must(uuid.NewV4()).String(),
		"sig": hex.EncodeToString(sum[:]),
	})

	return token.SignedString(user.privateKey)
}
