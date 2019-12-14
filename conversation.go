package sdk

import (
	"context"
	"time"
)

// Participant conversation participant
type Participant struct {
	Type      string     `json:"type,omitempty"`
	UserID    string     `json:"user_id,omitempty"`
	Role      string     `json:"role,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

// Conversation conversation
type Conversation struct {
	ConversationID string `json:"conversation_id,omitempty"`
	CreatorID      string `json:"creator_id,omitempty"`
	Category       string `json:"category,omitempty"`
	Name           string `json:"name,omitempty"`
	IconURL        string `json:"icon_url,omitempty"`
	Announcement   string `json:"announcement,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	CodeID         string `json:"code_id,omitempty"`
	CodeURL        string `json:"code_url,omitempty"`

	Participants []*Participant `json:"participants,omitempty"`
}

// CreateConversation crate conversation
func (user *User) CreateConversation(ctx context.Context, category, conversationID, name, action, role, userID string, participants []*Participant) (*Conversation, error) {
	params := map[string]interface{}{
		"category": category,
	}
	if conversationID != "" {
		params["conversation_id"] = conversationID
	}
	if name != "" {
		params["name"] = name
	}
	if action != "" {
		params["action"] = action
	}
	if role != "" {
		params["role"] = role
	}
	if userID != "" {
		params["user_id"] = userID
	}
	if len(participants) > 0 {
		params["participants"] = participants
	}

	var conversation Conversation
	if err := user.Request(ctx, "POST", "/conversations", params, &conversation); err != nil {
		return nil, err
	}
	return &conversation, nil
}

// ReadConversation read conversation
func (user *User) ReadConversation(ctx context.Context, conversationID string) (*Conversation, error) {
	var conversation Conversation
	if err := user.Request(ctx, "GET", "/conversations/"+conversationID, nil, &conversation); err != nil {
		return nil, err
	}
	return &conversation, nil
}
