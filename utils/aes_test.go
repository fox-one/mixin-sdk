package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpending(t *testing.T) {
	text := "text"

	{
		expended := expendKey([]byte(text), 16)
		assert.Equal(t, "texttexttexttext", string(expended))
	}

	{
		expended := expendKey([]byte(text), 32)
		assert.Equal(t, "texttexttexttexttexttexttexttext", string(expended))
	}
}

func TestPadding(t *testing.T) {
	text := "text"

	padded := PKCS7Padding([]byte(text), 6)
	assert.Equal(t, "text\x02\x02", string(padded))

	text1 := UnPKCS7Padding(padded)
	assert.Equal(t, text, string(text1))
}

func TestAES(t *testing.T) {
	text := "hello world"

	key := "test"
	iv := "test"

	// AES-128
	{
		encrypted, err := Encrypt([]byte(text), []byte(key), []byte(iv))
		assert.Nil(t, err)
		assert.Equal(t, "Oj9NLVqWPnhs7sIgrwifeg==", encrypted)

		plain, err := Decrypt(encrypted, []byte(key), []byte(iv))
		assert.Nil(t, err)
		assert.Equal(t, text, string(plain))
	}

	// AES-192
	{
		encrypted, err := Encrypt([]byte(text), []byte(key), []byte(iv), 24)
		assert.Nil(t, err)
		assert.Equal(t, "hxoOCb8inPlPsJL9/3MTag==", encrypted)

		plain, err := Decrypt(encrypted, []byte(key), []byte(iv), 24)
		assert.Nil(t, err)
		assert.Equal(t, text, string(plain))
	}

	// AES-192
	{
		encrypted, err := Encrypt([]byte(text), []byte(key), []byte(iv), 32)
		assert.Nil(t, err)
		assert.Equal(t, "ko43YhBBei89Mc3tja0g9g==", encrypted)

		plain, err := Decrypt(encrypted, []byte(key), []byte(iv), 32)
		assert.Nil(t, err)
		assert.Equal(t, text, string(plain))
	}
}
