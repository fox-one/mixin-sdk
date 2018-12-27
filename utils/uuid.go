package utils

import (
	"crypto/md5"
	"io"

	"github.com/satori/go.uuid"
)

// UUIDWithString generate uuid with text
func UUIDWithString(text string) string {
	h := md5.New()
	io.WriteString(h, text)
	sum := h.Sum(nil)
	sum[6] = (sum[6] & 0x0f) | 0x30
	sum[8] = (sum[8] & 0x3f) | 0x80
	return uuid.FromBytesOrNil(sum).String()
}
