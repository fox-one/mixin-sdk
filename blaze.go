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
	"golang.org/x/sync/errgroup"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 10 * time.Second
	pingPeriod = pongWait * 8 / 10

	ackBatch = 80

	CreateMessageAction      = "CREATE_MESSAGE"
	AcknowledgeReceiptAction = "ACKNOWLEDGE_MESSAGE_RECEIPT"
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

	// ack status
	ack bool
}

func (m *MessageView) reset() {
	m.ack = false
	m.RepresentativeID = ""
	m.QuoteMessageID = ""
}

// Ack mark messageView as acked
// otherwise sdk will ack this message
func (m *MessageView) Ack() {
	m.ack = true
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
	user         *User
	readDeadline time.Time
}

type BlazeListener interface {
	OnAckReceipt(ctx context.Context, msg *MessageView, userID string) error
	OnMessage(ctx context.Context, msg *MessageView, userID string) error
}

func NewBlazeClient(user *User) *BlazeClient {
	client := BlazeClient{
		user: user,
	}
	return &client
}

func (b *BlazeClient) SetReadDeadline(conn *websocket.Conn, t time.Time) error {
	if err := conn.SetReadDeadline(t); err != nil {
		return err
	}

	b.readDeadline = t
	return nil
}

func (b *BlazeClient) Loop(ctx context.Context, listener BlazeListener) error {
	conn, err := connectMixinBlaze(b.user)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go tick(ctx, conn)

	ackBuffer := make(chan string)
	defer close(ackBuffer)

	go b.ack(ctx, conn, ackBuffer)

	_ = b.SetReadDeadline(conn, time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		return b.SetReadDeadline(conn, time.Now().Add(pongWait))
	})

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

		message.reset()
		if err := jsoniter.Unmarshal(blazeMessage.Data, &message); err != nil {
			continue
		}

		switch blazeMessage.Action {
		case CreateMessageAction:
			messageID := message.MessageID
			if err := listener.OnMessage(ctx, &message, b.user.UserID); err != nil {
				return err
			}

			if !message.ack {
				ackBuffer <- messageID
			}
		case AcknowledgeReceiptAction:
			if err := listener.OnAckReceipt(ctx, &message, b.user.UserID); err != nil {
				return err
			}
		}

		if time.Until(b.readDeadline) < time.Second {
			// 可能因为收到的消息过多或者消息处理太慢或者 ack 太慢
			// 导致没有及时处理 pong frame 而 read deadline 没有刷新
			// 这种情况下不应该读超时，在这里重置一下 read deadline
			_ = b.SetReadDeadline(conn, time.Now().Add(pongWait))
		}
	}
}

func connectMixinBlaze(user *User) (*websocket.Conn, error) {
	token, err := user.SignToken(RequestSig("GET", "/", nil), uuid.New(), time.Minute)
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

func tick(ctx context.Context, conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		conn.Close()
		ticker.Stop()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait)); err != nil {
				return
			}
		}
	}
}

func (b *BlazeClient) ack(ctx context.Context, _ *websocket.Conn, ackBuffer <-chan string) {
	const dur = time.Second
	t := time.NewTimer(dur)

	const maxBatch = 8 * ackBatch // 640

	requests := make([]*AcknowledgementRequest, 0, ackBatch)

	for {
		select {
		case id, ok := <-ackBuffer:
			if !ok {
				return
			}

			requests = append(requests, &AcknowledgementRequest{
				MessageID: id,
				Status:    "READ",
			})

			if count := len(requests); count >= maxBatch {
				count = maxBatch
				if err := b.sendAcknowledgements(ctx, requests[:count]); err == nil {
					remain := requests[count:]
					copy(requests, remain)
					requests = requests[:len(remain)]

					if len(requests) == 0 {
						if !t.Stop() {
							<-t.C
						}

						t.Reset(dur)
					}
				}
			}
		case <-t.C:
			if count := len(requests); count > 0 {
				if count > maxBatch {
					count = maxBatch
				}

				if err := b.sendAcknowledgements(ctx, requests[:count]); err == nil {
					remain := requests[count:]
					copy(requests, remain)
					requests = requests[:len(remain)]
				}
			}

			t.Reset(dur)
		}
	}
}

func (b *BlazeClient) sendAcknowledgements(ctx context.Context, requests []*AcknowledgementRequest) error {
	if len(requests) <= ackBatch {
		return b.user.SendAcknowledgements(ctx, requests)
	}

	var g errgroup.Group
	for idx := 0; idx < len(requests); idx += ackBatch {
		right := idx + ackBatch
		if right > len(requests) {
			right = len(requests)
		}

		batch := requests[idx:right]
		g.Go(func() error {
			return b.user.SendAcknowledgements(ctx, batch)
		})
	}

	return g.Wait()
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
	if err := conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}

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
