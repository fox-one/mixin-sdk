package mixin

import (
	"context"
)

type Profile struct {
	UserID         string `json:"user_id"`
	IdentityNumber string `json:"identity_number"`
	FullName       string `json:"full_name,omitempty"`
	AvatarURL      string `json:"avatar_url,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	Phone          string `json:"phone,omitempty"`
}

func fetchProfile(ctx context.Context) (*Profile, error) {
	resp, err := Request(ctx).Get("/me")
	if err != nil {
		return nil, err
	}

	var profile Profile
	if err := UnmarshalResponse(resp, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// FetchProfile fetch my profile
func (user *User) FetchProfile(ctx context.Context) (*Profile, error) {
	ctx = WithAuth(ctx, user)
	return fetchProfile(ctx)
}

// FetchProfile fetch user profile with edotoken
func (ed *EdOToken) FetchProfile(ctx context.Context) (*Profile, error) {
	ctx = WithAuth(ctx, ed)
	return fetchProfile(ctx)
}

// FetchProfile fetch user profile with accesstoken
func FetchProfile(ctx context.Context, accessToken string) (*Profile, error) {
	ctx = WithToken(ctx, accessToken)
	return fetchProfile(ctx)
}

func fetchFriends(ctx context.Context) ([]*User, error) {
	resp, err := Request(ctx).Get("/friends")
	if err != nil {
		return nil, err
	}

	var friends []*User
	if err := UnmarshalResponse(resp, &friends); err != nil {
		return nil, err
	}

	return friends, nil
}

// FetchFriends fetch friends with accesstoken
func FetchFriends(ctx context.Context, accessToken string) ([]*User, error) {
	ctx = WithToken(ctx, accessToken)

	return fetchFriends(ctx)
}

// FetchFriends fetch friends with edo token
func (ed *EdOToken) FetchFriends(ctx context.Context) ([]*User, error) {
	ctx = WithAuth(ctx, ed)

	return fetchFriends(ctx)
}

func UserMe(ctx context.Context, accessToken string) (*Profile, error) {
	return FetchProfile(ctx, accessToken)
}

// ModifyProfile update my profile
func (user *User) ModifyProfile(ctx context.Context, fullname, avatarBase64 string) (*Profile, error) {
	paras := map[string]interface{}{}
	if len(fullname) > 0 {
		paras["full_name"] = fullname
	}
	if len(avatarBase64) > 0 {
		paras["avatar_base64"] = avatarBase64
	}

	var profile Profile
	if err := user.Request(ctx, "POST", "/me", paras, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

// ModifyPreference update my preference
func (user *User) ModifyPreference(ctx context.Context, receiveMessageSource, acceptConversationSource string) (*Profile, error) {
	paras := map[string]interface{}{}
	if len(receiveMessageSource) > 0 {
		paras["receive_message_source"] = receiveMessageSource
	}
	if len(acceptConversationSource) > 0 {
		paras["accept_conversation_source"] = acceptConversationSource
	}

	var profile Profile
	if err := user.Request(ctx, "POST", "/me/preferences", paras, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

// FetchUsers fetch users
func (user *User) FetchUsers(ctx context.Context, userIDS ...string) ([]*Profile, error) {
	if len(userIDS) == 0 {
		return nil, nil
	}

	var profiles []*Profile
	if err := user.Request(ctx, "POST", "/users/fetch", userIDS, &profiles); err != nil {
		return nil, err
	}
	return profiles, nil
}

// FetchUser fetch user
func (user *User) FetchUser(ctx context.Context, userID string) (*Profile, error) {
	var profile Profile
	if err := user.Request(ctx, "GET", "/users/"+userID, nil, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

// SearchUser search user; q is String: Mixin Id or Phone number
func (user *User) SearchUser(ctx context.Context, q string) (*Profile, error) {
	var profile Profile
	if err := user.Request(ctx, "GET", "/search/"+q, nil, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

// FetchFriends fetch friends
func (user *User) FetchFriends(ctx context.Context) ([]*Profile, error) {
	var profiles []*Profile
	if err := user.Request(ctx, "GET", "/friends", nil, &profiles); err != nil {
		return nil, err
	}
	return profiles, nil
}
