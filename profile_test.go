package mixin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchProfile(t *testing.T) {
	ctx := context.Background()
	const token = "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE5MDc4MTEyMTIsImlhdCI6MTU5MjQ1MTIxMiwianRpIjoiNGQzZGU0ZWQtOWQ0NS00ZGUyLWJkNWItOWEyZDliYjc2ZmM2Iiwic2NwIjoiRlVMTCIsInNpZCI6ImRiMmYzMmJiLWYyYTUtNDJiMS1iOTQ2LTYzYTRlMTI5YjAyYyIsInNpZyI6IjVlNmI1OGZmYTEwYjNiYzUxNzI0ZmYwYmJkMmFmYjkxYzQ3NzFlZTM0MGY1ZDY4NTM0MGRmYTRjODU0YmFmYmEiLCJ1aWQiOiI1YzRmMzBhNi0xZjQ5LTQzYzMtYjM3Yi1jMDFhYWU1MTkxYWYifQ.ni3kFh2yI9HqGldTR3GMCNlpfE6OXrpUYduhZGmZ9-wbCyFEO9zZeiZtOjzqKrIMdLsL6DIBsZibvaJp_UVopYFiK9tLs9_8OVg31cETtsq9YOSvnAw346kRn6iw3OmDuhJm4afVAsmtvUyxAS9FoeCxrDiUDKMuv1h3TjeFzBY"
	profile, err := FetchProfile(ctx, token)
	if assert.Nilf(t, err, "FetchProfile should success") {
		assert.Lenf(t, profile.UserID, 36, "user id should be uuid")
	}
}
