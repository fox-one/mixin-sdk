package messenger

import (
	"context"

	mixinsdk "github.com/fox-one/mixin-sdk"
)

// User messenger user entity
type User struct {
	UserID         string `json:"user_id"`
	IdentityNumber string `json:"identity_number"`
	FullName       string `json:"full_name,omitempty"`
	AvatarURL      string `json:"avatar_url,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	Phone          string `json:"phone,omitempty"`
}

func fetchProfile(ctx context.Context) (*User, error) {
	resp, err := mixinsdk.Request(ctx).Get("/me")
	if err != nil {
		return nil, err
	}

	var user User
	err = mixinsdk.UnmarshalResponse(resp, &user)
	return &user, err
}

// FetchProfile fetch my profile
func (m Messenger) FetchProfile(ctx context.Context) (*User, error) {
	ctx = mixinsdk.WithAuth(ctx, m.User)
	return fetchProfile(ctx)
}

func UserMe(ctx context.Context, accessToken string) (*User, error) {
	ctx = mixinsdk.WithToken(ctx, accessToken)
	return fetchProfile(ctx)
}

// ModifyProfile update my profile
func (m Messenger) ModifyProfile(ctx context.Context, fullname, avatarBase64 string) (*User, error) {
	paras := map[string]interface{}{}
	if len(fullname) > 0 {
		paras["full_name"] = fullname
	}
	if len(avatarBase64) > 0 {
		paras["avatar_base64"] = avatarBase64
	}

	var user User
	if err := m.SendRequest(ctx, "POST", "/me", paras, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// ModifyPreference update my preference
func (m Messenger) ModifyPreference(ctx context.Context, receiveMessageSource, acceptConversationSource string) (*User, error) {
	paras := map[string]interface{}{}
	if len(receiveMessageSource) > 0 {
		paras["receive_message_source"] = receiveMessageSource
	}
	if len(acceptConversationSource) > 0 {
		paras["accept_conversation_source"] = acceptConversationSource
	}

	var user User
	if err := m.SendRequest(ctx, "POST", "/me/preferences", paras, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// FetchUsers fetch users
func (m Messenger) FetchUsers(ctx context.Context, userIDS ...string) ([]*User, error) {
	if len(userIDS) == 0 {
		return []*User{}, nil
	}

	var users []*User
	if err := m.SendRequest(ctx, "POST", "/users/fetch", userIDS, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// FetchUser fetch user
func (m Messenger) FetchUser(ctx context.Context, userID string) (*User, error) {
	var user User
	if err := m.SendRequest(ctx, "GET", "/users/"+userID, nil, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// SearchUser search user; q is String: Mixin Id or Phone Numbe
func (m Messenger) SearchUser(ctx context.Context, q string) (*User, error) {
	var user User
	if err := m.SendRequest(ctx, "GET", "/search/"+q, nil, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// FetchFriends fetch friends
func (m Messenger) FetchFriends(ctx context.Context) ([]*User, error) {
	var users []*User
	if err := m.SendRequest(ctx, "GET", "/friends", nil, &users); err != nil {
		return nil, err
	}
	return users, nil
}
