package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBool(t *testing.T) {
	data := map[string]interface{}{
		"a": true,
		"b": false,
	}

	a, _ := data["a"]
	assert.Equal(t, 1, ParseInt(a))

	b, _ := data["b"]
	assert.Equal(t, 0, ParseInt(b))
}
