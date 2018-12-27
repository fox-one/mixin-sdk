package utils

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetheaders(t *testing.T) {
	req, _ := http.NewRequest("GET", "www.fox.one", bytes.NewBufferString("haha"))
	SetHeaders(req, "Sk", "sk1")
	SetHeaders(req, "sk", "sk2")

	if sk := req.Header.Get("sk"); sk != "sk1" {
		t.Error("sk should be sk1")
	}
}

func TestBuildUrl(t *testing.T) {
	path := "https://www.fox.one?a=b"
	u, err := BuildURL(path, "symbol", "eos", "limit")
	assert.Nil(t, err)
	assert.Equal(t, path+"&symbol=eos", u)
}
