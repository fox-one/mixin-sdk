package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipEncoding(t *testing.T) {
	text := "hello world"
	ret, err := GZipEncode([]byte(text))
	assert.Nil(t, err)
	decoded, err := GZipDecode(ret)
	assert.Nil(t, err)
	assert.Equal(t, text, string(decoded))
}
