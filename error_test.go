package mixin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsErrorCodes(t *testing.T) {
	err := createError(200, 401, "invalid token")
	assert.True(t, IsErrorCodes(err, 400, 401))
	assert.False(t, IsErrorCodes(err, 400, 403))
}
