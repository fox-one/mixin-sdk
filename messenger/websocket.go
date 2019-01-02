package messenger

import (
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/fox-one/mixin-sdk/mixin"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

const (
	keepAlivePeriod = 3 * time.Second
	writeWait       = 10 * time.Second
	pongWait        = 10 * time.Second
	pingPeriod      = (pongWait * 9) / 10

	createMessageAction = "CREATE_MESSAGE"
)

const (
	MessageCategoryPlainText             = "PLAIN_TEXT"
	MessageCategoryPlainImage            = "PLAIN_IMAGE"
	MessageCategoryPlainData             = "PLAIN_DATA"
	MessageCategoryPlainSticker          = "PLAIN_STICKER"
	MessageCategorySystemConversation    = "SYSTEM_CONVERSATION"
	MessageCategorySystemAccountSnapshot = "SYSTEM_ACCOUNT_SNAPSHOT"
)

type BlazeMessage struct {
	Id     string                 `json:"id"`
	Action string                 `json:"action"`
	Params map[string]interface{} `json:"params,omitempty"`
	Data   interface{}            `json:"data,omitempty"`
	Error  *mixin.Error           `json:"error,omitempty"`
}

type MessageView struct {
	ConversationId string    `json:"conversation_id"`
	UserId         string    `json:"user_id"`
	MessageId      string    `json:"message_id"`
	Category       string    `json:"category"`
	Data           string    `json:"data"`
	Status         string    `json:"status"`
	Source         string    `json:"source"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type TransferView struct {
	Type          string    `json:"type"`
	SnapshotId    string    `json:"snapshot_id"`
	CounterUserId string    `json:"counter_user_id"`
	AssetId       string    `json:"asset_id"`
	Amount        string    `json:"amount"`
	TraceId       string    `json:"trace_id"`
	Memo          string    `json:"memo"`
	CreatedAt     time.Time `json:"created_at"`
}

type messageContext struct {
	transactions *tmap
	readDone     chan bool
	writeDone    chan bool
	readBuffer   chan MessageView
	writeBuffer  chan []byte
}

type systemConversationPayload struct {
	Action        string `json:"action"`
	ParticipantId string `json:"participant_id"`
	UserId        string `json:"user_id,omitempty"`
	Role          string `json:"role,omitempty"`
}

type BlazeClient struct {
	mc  *messageContext
	uid string
	sid string
	key string
}

type BlazeListener interface {
	OnMessage(ctx context.Context, msg MessageView, userId string) error
}

func NewBlazeClient(uid, sid, key string) *BlazeClient {
	client := BlazeClient{
		mc: &messageContext{
			transactions: newTmap(),
			readDone:     make(chan bool, 1),
			writeDone:    make(chan bool, 1),
			readBuffer:   make(chan MessageView, 102400),
			writeBuffer:  make(chan []byte, 102400),
		},
		uid: uid,
		sid: sid,
		key: key,
	}
	return &client
}

// Add
func BlazeServerError(ctx context.Context, err error) error {
	return ServerError(err)
}
func UuidNewV4() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

func (b *Messenger) Loop(ctx context.Context, listener BlazeListener) error {
	//func (b *BlazeClient) Loop(ctx context.Context, listener BlazeListener) error {
	//conn, err := connectMixinBlaze(b.uid, b.sid, b.key)
	conn, err := b.connectMixinBlaze()
	if err != nil {
		return err
	}
	defer conn.Close()
	go writePump(ctx, conn, b.mc)
	go readPump(ctx, conn, b.mc)
	if err = writeMessageAndWait(ctx, b.mc, "LIST_PENDING_MESSAGES", nil); err != nil {
		return BlazeServerError(ctx, err)
	}
	for {
		select {
		case <-b.mc.readDone:
			return nil
		case msg := <-b.mc.readBuffer:
			err = listener.OnMessage(ctx, msg, b.uid)
			if err != nil {
				return err
			}
			params := map[string]interface{}{"message_id": msg.MessageId, "status": "READ"}
			if err = writeMessageAndWait(ctx, b.mc, "ACKNOWLEDGE_MESSAGE_RECEIPT", params); err != nil {
				return BlazeServerError(ctx, err)
			}
		}
	}
}

func (b *BlazeClient) SendMessage(ctx context.Context, conversationId, recipientId, category, content, representativeId string) error {
	params := map[string]interface{}{
		"conversation_id":   conversationId,
		"recipient_id":      recipientId,
		"message_id":        UuidNewV4().String(),
		"representative_id": representativeId,
		"category":          category,
		"data":              base64.StdEncoding.EncodeToString([]byte(content)),
	}
	if err := writeMessageAndWait(ctx, b.mc, createMessageAction, params); err != nil {
		return BlazeServerError(ctx, err)
	}
	return nil
}

func (b *BlazeClient) SendPlainText(ctx context.Context, msg MessageView, content string) error {
	params := map[string]interface{}{
		"conversation_id": msg.ConversationId,
		"recipient_id":    msg.UserId,
		"message_id":      UuidNewV4().String(),
		"category":        "PLAIN_TEXT",
		"data":            base64.StdEncoding.EncodeToString([]byte(content)),
	}
	if err := writeMessageAndWait(ctx, b.mc, createMessageAction, params); err != nil {
		return BlazeServerError(ctx, err)
	}
	return nil
}

func (b *BlazeClient) SendContact(ctx context.Context, conversationId, recipientId, contactId string) error {
	contactMap := map[string]string{"user_id": contactId}
	contactData, _ := json.Marshal(contactMap)
	params := map[string]interface{}{
		"conversation_id": conversationId,
		"recipient_id":    recipientId,
		"message_id":      UuidNewV4().String(),
		"category":        "PLAIN_CONTACT",
		"data":            base64.StdEncoding.EncodeToString(contactData),
	}
	if err := writeMessageAndWait(ctx, b.mc, createMessageAction, params); err != nil {
		return BlazeServerError(ctx, err)
	}
	return nil
}

func (b *BlazeClient) SendAppButton(ctx context.Context, conversationId, recipientId, label, action, color string) error {
	btns, err := json.Marshal([]interface{}{map[string]string{
		"label":  label,
		"action": action,
		"color":  color,
	}})
	if err != nil {
		return BlazeServerError(ctx, err)
	}
	params := map[string]interface{}{
		"conversation_id": conversationId,
		"recipient_id":    recipientId,
		"message_id":      UuidNewV4().String(),
		"category":        "APP_BUTTON_GROUP",
		"data":            base64.StdEncoding.EncodeToString(btns),
	}
	err = writeMessageAndWait(ctx, b.mc, createMessageAction, params)
	if err != nil {
		return BlazeServerError(ctx, err)
	}
	return nil
}

type Button struct {
	Label  string `json:"label"`
	Action string `json:"action"`
	Color  string `json:"color"`
}

func (b *BlazeClient) SendAppButtons(ctx context.Context, conversationId, recipientId string, buttons ...Button) error {
	btns, err := json.Marshal(buttons)
	if err != nil {
		return BlazeServerError(ctx, err)
	}
	params := map[string]interface{}{
		"conversation_id": conversationId,
		"recipient_id":    recipientId,
		"message_id":      UuidNewV4().String(),
		"category":        "APP_BUTTON_GROUP",
		"data":            base64.StdEncoding.EncodeToString(btns),
	}
	err = writeMessageAndWait(ctx, b.mc, createMessageAction, params)
	if err != nil {
		return BlazeServerError(ctx, err)
	}
	return nil
}
func (m Messenger) connectMixinBlaze() (*websocket.Conn, error) {
	//func connectMixinBlaze(uid, sid, key string) (*websocket.Conn, error) {
	//token, err := SignAuthenticationToken(uid, sid, key, "GET", "/", "")
	token, err := m.SignToken("GET", "/", nil)
	if err != nil {
		return nil, err
	}
	header := make(http.Header)
	header.Add("Authorization", "Bearer "+token)
	u := url.URL{Scheme: "wss", Host: "blaze.mixin.one", Path: "/"}
	dialer := &websocket.Dialer{
		Subprotocols: []string{"Mixin-Blaze-1"},
	}
	conn, _, err := dialer.Dial(u.String(), header)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func readPump(ctx context.Context, conn *websocket.Conn, mc *messageContext) error {
	defer func() {
		conn.Close()
		mc.writeDone <- true
		mc.readDone <- true
	}()

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		err := conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return BlazeServerError(ctx, err)
		}
		return nil
	})

	for {
		err := conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return BlazeServerError(ctx, err)
		}
		messageType, wsReader, err := conn.NextReader()
		if err != nil {
			return BlazeServerError(ctx, err)
		}
		if messageType != websocket.BinaryMessage {
			return BlazeServerError(ctx, fmt.Errorf("invalid message type %d", messageType))
		}
		err = parseMessage(ctx, mc, wsReader)
		if err != nil {
			return BlazeServerError(ctx, err)
		}
	}
}

func writePump(ctx context.Context, conn *websocket.Conn, mc *messageContext) error {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()
	for {
		select {
		case data := <-mc.writeBuffer:
			err := writeGzipToConn(conn, data)
			if err != nil {
				return BlazeServerError(ctx, err)
			}
		case <-mc.writeDone:
			return nil
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return BlazeServerError(ctx, err)
			}
		}
	}
}

func writeMessageAndWait(ctx context.Context, mc *messageContext, action string, params map[string]interface{}) error {
	var resp = make(chan BlazeMessage, 1)
	var id = UuidNewV4().String()
	mc.transactions.set(id, func(t BlazeMessage) error {
		select {
		case resp <- t:
		case <-time.After(1 * time.Second):
			return fmt.Errorf("timeout to hook %s %s", action, id)
		}
		return nil
	})
	blazeMessage, err := json.Marshal(BlazeMessage{Id: id, Action: action, Params: params})
	if err != nil {
		return err
	}
	select {
	case <-time.After(keepAlivePeriod):
		return fmt.Errorf("timeout to write %s %v", action, params)
	case mc.writeBuffer <- blazeMessage:
	}
	select {
	case <-time.After(keepAlivePeriod):
		return fmt.Errorf("timeout to wait %s %v", action, params)
	case t := <-resp:
		if t.Error != nil && t.Error.Code != 403 {
			//return writeMessageAndWait(ctx, mc, action, params)
			return BlazeServerError(ctx, err)
		}
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
	if err := wsWriter.Close(); err != nil {
		return err
	}
	return nil
}

func parseMessage(ctx context.Context, mc *messageContext, wsReader io.Reader) error {
	var message BlazeMessage
	gzReader, err := gzip.NewReader(wsReader)
	if err != nil {
		return err
	}
	defer gzReader.Close()
	if err = json.NewDecoder(gzReader).Decode(&message); err != nil {
		return err
	}
	transaction := mc.transactions.retrive(message.Id)
	if transaction != nil {
		return transaction(message)
	}
	if message.Action != "CREATE_MESSAGE" {
		return nil
	}
	data, err := json.Marshal(message.Data)
	if err != nil {
		return err
	}
	var msg MessageView
	if err = json.Unmarshal(data, &msg); err != nil {
		return err
	}

	select {
	case <-time.After(keepAlivePeriod):
		return fmt.Errorf("timeout to handle %s %s", msg.Category, msg.MessageId)
	case mc.readBuffer <- msg:
	}
	return nil
}

type tmap struct {
	mutex sync.Mutex
	m     map[string]mixinTransaction
}

type mixinTransaction func(BlazeMessage) error

func newTmap() *tmap {
	return &tmap{
		m: make(map[string]mixinTransaction),
	}
}

func (m *tmap) retrive(key string) mixinTransaction {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	defer delete(m.m, key)
	return m.m[key]
}

func (m *tmap) set(key string, t mixinTransaction) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.m[key] = t
}
