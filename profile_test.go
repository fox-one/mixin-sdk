package mixin

import (
	"context"
	"testing"

	"github.com/fox-one/pkg/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFetchProfile(t *testing.T) {
	ctx := context.Background()
	const token = "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE5MDc4MTEyMTIsImlhdCI6MTU5MjQ1MTIxMiwianRpIjoiNGQzZGU0ZWQtOWQ0NS00ZGUyLWJkNWItOWEyZDliYjc2ZmM2Iiwic2NwIjoiRlVMTCIsInNpZCI6ImRiMmYzMmJiLWYyYTUtNDJiMS1iOTQ2LTYzYTRlMTI5YjAyYyIsInNpZyI6IjVlNmI1OGZmYTEwYjNiYzUxNzI0ZmYwYmJkMmFmYjkxYzQ3NzFlZTM0MGY1ZDY4NTM0MGRmYTRjODU0YmFmYmEiLCJ1aWQiOiI1YzRmMzBhNi0xZjQ5LTQzYzMtYjM3Yi1jMDFhYWU1MTkxYWYifQ.ni3kFh2yI9HqGldTR3GMCNlpfE6OXrpUYduhZGmZ9-wbCyFEO9zZeiZtOjzqKrIMdLsL6DIBsZibvaJp_UVopYFiK9tLs9_8OVg31cETtsq9YOSvnAw346kRn6iw3OmDuhJm4afVAsmtvUyxAS9FoeCxrDiUDKMuv1h3TjeFzBY"
	profile, err := FetchProfile(WithRequestID(ctx, uuid.New()), token)
	if assert.Nilf(t, err, "FetchProfile should success") {
		assert.Lenf(t, profile.UserID, 36, "user id should be uuid")
	}
}

func TestSearchUser(t *testing.T) {
	token := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhaWQiOiJiNjlkMTE2NC0yZGMxLTRmN2UtYThlMy1mMjk2ZGMwMTZkODYiLCJleHAiOjE2Mjk0NDAxMjUsImlhdCI6MTU5NzkwNDEyNSwiaXNzIjoiMjgyN2Q4MWYtNmFlMC00ODQyLWI5MmYtNjU3NmFmZTM2ODYzIiwic2NwIjoiUEhPTkU6UkVBRCBQUk9GSUxFOlJFQUQgTUVTU0FHRVM6UkVQUkVTRU5UIEFTU0VUUzpSRUFEIn0.VVkev1XsbT4Np8RQ0Z2GSgpVj3d41ErODJl0qRi8V87GGW8Kc98WRVmcYzVUPljs2LyfSgS4fACKUV9K-USyNsMgQEsmLIxnfP089PcUGAaWcNFdO_8n2bfIzOhO9nWs8wN9Cx5t8W7AKX2EyUAH4zLtOG9vKe9Mc3kwDpDw4Zg"
	user, e := SearchUser(context.Background(), "37261734", token)
	if assert.Nilf(t, e, "search user success") {
		assert.NotNil(t, user)
	}
}
