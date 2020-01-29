package mixin

import (
	"context"
	"fmt"
	"time"

	"github.com/fox-one/mixin-sdk/utils"
)

const (
	ConversationCategoryContact = "CONTACT"
	ConversationCategoryGroup   = "GROUP"

	ParticipantActionAdd    = "ADD"
	ParticipantActionRemove = "REMOVE"
	ParticipantActionJoin   = "JOIN"
	ParticipantActionExit   = "EXIT"
	ParticipantActionRole   = "ROLE"

	ParticipantRoleAdmin  = "ADMIN"
	ParticipantRoleMember = ""
)

// Participant conversation participant
type Participant struct {
	Action    string     `json:"action,omitempty"`
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

type CreateConversationInput struct {
	Category       string         `json:"category,omitempty"`
	ConversationID string         `json:"conversation_id,omitempty"`
	Name           string         `json:"name,omitempty"`
	Participants   []*Participant `json:"participants,omitempty"`
}

// CreateConversation crate conversation
func (user *User) CreateConversation(ctx context.Context, input *CreateConversationInput) (*Conversation, error) {
	var conversation Conversation
	if err := user.Request(ctx, "POST", "/conversations", input, &conversation); err != nil {
		return nil, err
	}
	return &conversation, nil
}

// CreateContactConversation create a conversation with a mixin messenger user
func (user *User) CreateContactConversation(ctx context.Context, userID string) (*Conversation, error) {
	return user.CreateConversation(ctx, &CreateConversationInput{
		Category:       ConversationCategoryContact,
		ConversationID: utils.UniqueConversationID(user.UserID, userID),
		Participants:   []*Participant{{UserID: userID}},
	})
}

// CreateGroupConversation create a group in mixin messenger with given participants
func (user *User) CreateGroupConversation(ctx context.Context, conversationID, name string, participants []*Participant) (*Conversation, error) {
	return user.CreateConversation(ctx, &CreateConversationInput{
		Category:       ConversationCategoryGroup,
		ConversationID: conversationID,
		Name:           name,
		Participants:   participants,
	})
}

// ReadConversation read conversation
func (user *User) ReadConversation(ctx context.Context, conversationID string) (*Conversation, error) {
	var conversation Conversation
	if err := user.Request(ctx, "GET", "/conversations/"+conversationID, nil, &conversation); err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (user *User) ManageConversation(ctx context.Context, conversationID, action string, participants []*Participant) (*Conversation, error) {
	var conversation Conversation
	path := fmt.Sprintf("/conversations/%s/participants/%s", conversationID, action)
	if err := user.Request(ctx, "POST", path, participants, &conversation); err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (user *User) AddParticipants(ctx context.Context, conversationID string, users ...string) (*Conversation, error) {
	var participants []*Participant
	for _, user := range users {
		participants = append(participants, &Participant{UserID: user})
	}

	return user.ManageConversation(ctx, conversationID, ParticipantActionAdd, participants)
}

func (user *User) RemoveParticipants(ctx context.Context, conversationID string, users ...string) (*Conversation, error) {
	var participants []*Participant
	for _, user := range users {
		participants = append(participants, &Participant{UserID: user})
	}

	return user.ManageConversation(ctx, conversationID, ParticipantActionRemove, participants)
}

func (user *User) AdminParticipants(ctx context.Context, conversationID string, users ...string) (*Conversation, error) {
	var participants []*Participant
	for _, user := range users {
		participants = append(participants, &Participant{
			UserID: user,
			Role:   ParticipantRoleAdmin,
		})
	}

	return user.ManageConversation(ctx, conversationID, ParticipantActionRole, participants)
}
