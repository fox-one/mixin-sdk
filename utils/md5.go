package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 md5
func MD5(str string) []byte {
	h := md5.New()
	h.Write([]byte(str))
	return h.Sum(nil)
}

// MD5Hex md5 hexdigest
func MD5Hex(str string) string {
	return hex.EncodeToString(MD5(str))
}
