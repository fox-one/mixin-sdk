package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

var ckey cipher.Block

func expendKey(key []byte) []byte {
	for len(key) < 16 {
		key = append(key, key...)
	}
	return key[:16]
}

// PKCS7Padding PKCS7补码, 可以参考下http://blog.studygolang.com/167.html
func PKCS7Padding(data []byte) []byte {
	blockSize := 16
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)

}

// UnPKCS7Padding 去除PKCS7的补码
func UnPKCS7Padding(data []byte) []byte {
	length := len(data)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(data[length-1])
	if length <= unpadding {
		return nil
	}
	return data[:(length - unpadding)]
}

// Encrypt aes encrypt
func Encrypt(data []byte, key, iv []byte) (string, error) {
	key = expendKey(key)
	iv = expendKey(iv)

	var err error
	ckey, err = aes.NewCipher(key)
	if nil != err {
		return "", err
	}

	encrypter := cipher.NewCBCEncrypter(ckey, iv)

	// PKCS7补码
	str := PKCS7Padding([]byte(data))
	out := make([]byte, len(str))

	encrypter.CryptBlocks(out, str)

	return base64.StdEncoding.EncodeToString(out), nil
}

// Decrypt aes decrypt
func Decrypt(base64Str string, key, iv []byte) ([]byte, error) {
	key = expendKey(key)
	iv = expendKey(iv)

	var err error
	ckey, err = aes.NewCipher(key)
	if nil != err {
		return nil, err
	}

	decrypter := cipher.NewCBCDecrypter(ckey, iv)

	base64In, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}

	out := make([]byte, len(base64In))
	decrypter.CryptBlocks(out, base64In)

	// 去除PKCS7补码
	out = UnPKCS7Padding(out)
	return out, nil
}
