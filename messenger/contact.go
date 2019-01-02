package messenger

import (
	"context"
)

// Participant conversation participant
type Participant struct {
	Type      string `json:"type"`
	UserID    string `json:"userID"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

// Conversation conversation
type Conversation struct {
	ConversationID string `json:"conversation_id"`
	CreatorID      string `json:"creator_id"`
	Category       string `json:"category,omitempty"`
	Name           string `json:"name,omitempty"`
	IconURL        string `json:"icon_url,omitempty"`
	Announcement   string `json:"announcement,omitempty"`
	CreatedAt      string `json:"created_at"`
	CodeID         string `json:"code_id,omitempty"`
	CodeURL        string `json:"code_url,omitempty"`

	Participants []*Participant `json:"participants"`
}

// CreateConversation crate conversation
func (m Messenger) CreateConversation(ctx context.Context, category, participants, action, role, userID string) (*Conversation, error) {
	return nil, nil
}
