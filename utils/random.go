package utils

import (
	"bytes"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// RandomStr random str
func RandomStr(length int) string {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	alpLen := len(alphanum)

	var result bytes.Buffer
	for i := 0; i < length; i++ {
		result.WriteString(string(alphanum[RandInt(0, alpLen)]))
	}

	return result.String()
}

// RandInt random integer value
func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
