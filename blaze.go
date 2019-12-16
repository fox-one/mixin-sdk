package mixin

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/fox-one/pkg/uuid"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 10 * time.Second
	pingPeriod = pongWait * 8 / 10

	createMessageAction = "CREATE_MESSAGE"
)

const (
	MessageCategoryPlainText             = "PLAIN_TEXT"
	MessageCategoryPlainImage            = "PLAIN_IMAGE"
	MessageCategoryPlainData             = "PLAIN_DATA"
	MessageCategoryPlainSticker          = "PLAIN_STICKER"
	MessageCategoryPlainLive             = "PLAIN_LIVE"
	MessageCategoryPlainContact          = "PLAIN_CONTACT"
	MessageCategorySystemConversation    = "SYSTEM_CONVERSATION"
	MessageCategorySystemAccountSnapshot = "SYSTEM_ACCOUNT_SNAPSHOT"
	MessageCategoryMessageRecall         = "MESSAGE_RECALL"
	MessageCategoryAppButtonGroup        = "APP_BUTTON_GROUP"
)

type BlazeMessage struct {
	Id     string                 `json:"id"`
	Action string                 `json:"action"`
	Params map[string]interface{} `json:"params,omitempty"`
	Data   json.RawMessage        `json:"data,omitempty"`
	Error  *Error                 `json:"error,omitempty"`
}

type MessageView struct {
	ConversationID   string    `json:"conversation_id"`
	UserID           string    `json:"user_id"`
	MessageID        string    `json:"message_id"`
	Category         string    `json:"category"`
	Data             string    `json:"data"`
	RepresentativeID string    `json:"representative_id"`
	QuoteMessageID   string    `json:"quote_message_id"`
	Status           string    `json:"status"`
	Source           string    `json:"source"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type TransferView struct {
	Type          string    `json:"type"`
	SnapshotID    string    `json:"snapshot_id"`
	CounterUserID string    `json:"counter_user_id"`
	AssetID       string    `json:"asset_id"`
	Amount        string    `json:"amount"`
	TraceID       string    `json:"trace_id"`
	Memo          string    `json:"memo"`
	CreatedAt     time.Time `json:"created_at"`
}

type SystemConversationPayload struct {
	Action        string `json:"action"`
	ParticipantID string `json:"participant_id"`
	UserID        string `json:"user_id,omitempty"`
	Role          string `json:"role,omitempty"`
}

type BlazeClient struct {
	user *User
}

type BlazeListener interface {
	OnMessage(ctx context.Context, msg *MessageView, userId string) error
}

func NewBlazeClient(user *User) *BlazeClient {
	client := BlazeClient{
		user: user,
	}
	return &client
}

func (b *BlazeClient) Loop(ctx context.Context, listener BlazeListener) error {
	conn, err := connectMixinBlaze(b.user)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	go tick(ctx, conn)

	messageIds := make(chan string, 1)
	go ack(ctx, conn, messageIds)

	if err = writeMessage(conn, "LIST_PENDING_MESSAGES", nil); err != nil {
		return fmt.Errorf("write LIST_PENDING_MESSAGES failed: %w", err)
	}

	var (
		blazeMessage BlazeMessage
		message      MessageView
	)

	for {
		typ, r, err := conn.NextReader()
		if err != nil {
			return err
		}

		if typ != websocket.BinaryMessage {
			return fmt.Errorf("invalid message type %d", typ)
		}

		if err := parseBlazeMessage(r, &blazeMessage); err != nil {
			return err
		}

		if blazeMessage.Error != nil {
			return err
		}

		if blazeMessage.Action != createMessageAction {
			continue
		}

		if err := jsoniter.Unmarshal(blazeMessage.Data, &message); err != nil {
			return err
		}

		if err := listener.OnMessage(ctx, &message, b.user.UserID); err != nil {
			return err
		}

		messageIds <- message.MessageID
	}
}

func connectMixinBlaze(user *User) (*websocket.Conn, error) {
	token, err := user.SignToken("GET", "/", nil)
	if err != nil {
		return nil, err
	}

	header := make(http.Header)
	header.Add("Authorization", "Bearer "+token)
	u := url.URL{Scheme: "wss", Host: "mixin-blaze.zeromesh.net", Path: "/"}
	dialer := &websocket.Dialer{
		Subprotocols:   []string{"Mixin-Blaze-1"},
		ReadBufferSize: 1024,
	}
	conn, _, err := dialer.Dial(u.String(), header)
	if err != nil {
		return nil, err
	}

	// no limit
	conn.SetReadLimit(0)
	return conn, nil
}

func tick(ctx context.Context, conn *websocket.Conn) error {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = conn.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return fmt.Errorf("write PING failed: %w", err)
			}
		}
	}
}

func ack(ctx context.Context, conn *websocket.Conn, ids <-chan string) error {
	defer conn.Close()

	var requests []*AcknowledgementRequest

	const dur = time.Second
	t := time.NewTimer(dur)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case id := <-ids:
			requests = append(requests, &AcknowledgementRequest{
				MessageID: id,
				Status:    "READ",
			})

			// ack limit 是 80，这里设置得稍微低一点
			if len(requests) >= 70 {
				if err := writeMessage(conn, "ACKNOWLEDGE_MESSAGE_RECEIPTS", map[string]interface{}{
					"messages": requests,
				}); err != nil {
					return err
				}

				requests = []*AcknowledgementRequest{}

				if !t.Stop() {
					<-t.C
				}

				t.Reset(dur)
			}
		case <-t.C:
			if len(requests) > 0 {
				if err := writeMessage(conn, "ACKNOWLEDGE_MESSAGE_RECEIPTS", map[string]interface{}{
					"messages": requests,
				}); err != nil {
					return err
				}

				requests = []*AcknowledgementRequest{}
			}

			t.Reset(dur)
		}
	}
}

func writeMessage(coon *websocket.Conn, action string, params map[string]interface{}) error {
	id := uuid.New()
	blazeMessage, err := jsoniter.Marshal(BlazeMessage{Id: id, Action: action, Params: params})
	if err != nil {
		return err
	}

	if err := writeGzipToConn(coon, blazeMessage); err != nil {
		return err
	}

	return nil
}

func writeGzipToConn(conn *websocket.Conn, msg []byte) error {
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	wsWriter, err := conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}
	gzWriter, err := gzip.NewWriterLevel(wsWriter, 3)
	if err != nil {
		return err
	}
	if _, err := gzWriter.Write(msg); err != nil {
		return err
	}

	if err := gzWriter.Close(); err != nil {
		return err
	}
	return wsWriter.Close()
}

func parseBlazeMessage(r io.Reader, msg *BlazeMessage) error {
	gzReader, err := gzip.NewReader(r)
	if err != nil {
		return err
	}

	err = jsoniter.NewDecoder(gzReader).Decode(msg)
	_ = gzReader.Close()
	return err
}
