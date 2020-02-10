package mixin

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"io"
	"time"
)

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

	user.pinCipher = block
	return nil
}

func (user *User) EncryptPIN(pin string) (string, error) {
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
	mode := cipher.NewCBCEncrypter(user.pinCipher, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], pinByte)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// API

// ModifyPIN modify pin
func (user *User) ModifyPIN(ctx context.Context, oldPIN, pin string) error {
	if pin == oldPIN {
		return nil
	}

	pinEncrypted, err := user.EncryptPIN(oldPIN)
	if err != nil {
		return err
	}
	return user.RequestWithPIN(ctx, "POST", "/pin/update", map[string]interface{}{"old_pin": pinEncrypted}, pin, nil)
}

// VerifyPIN verify user pin
func (user *User) VerifyPIN(ctx context.Context, pin string) error {
	return user.RequestWithPIN(ctx, "POST", "/pin/verify", nil, pin, nil)
}
