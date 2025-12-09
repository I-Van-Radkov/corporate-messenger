package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/I-Van-Radkov/corporate-messenger/chat-service/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	Conn      *websocket.Conn
	Send      chan []byte

	IsClosed  atomic.Bool
	closeOnce sync.Once
	mu        sync.Mutex

	cancelFunc context.CancelFunc
	ctx        context.Context
}

func NewClient(userID uuid.UUID, conn *websocket.Conn) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		UserID:     userID,
		SessionID:  uuid.New(),
		Conn:       conn,
		Send:       make(chan []byte, 256),
		ctx:        ctx,
		cancelFunc: cancel,
	}
	client.IsClosed.Store(false)

	return client
}

func (c *Client) close() {
	c.closeOnce.Do(func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.IsClosed.Store(true)

		if c.cancelFunc != nil {
			c.cancelFunc()
		}

		if c.Conn != nil {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			c.Conn.Close()
		}

		close(c.Send)
	})
}

func (c *Client) getIsClosed() bool {
	return c.IsClosed.Load()
}

type ChatUsecase interface {
	SendMsgToStorage(ctx context.Context, userId uuid.UUID, msg []byte)
	GetMessagesFromStorage(ctx context.Context, userId uuid.UUID) [][]byte

	SendMessageToDb(ctx context.Context, msg *dto.MessageDTO) (string, error)
	//EditMessage(ctx context.Context, messageID, userID, content string) error
	//DeleteMessage(ctx context.Context, messageID, userID string) error
	//MarkAsRead(ctx context.Context, userID, chatID, messageID string) error
	GetChatMembers(ctx context.Context, chatID string) ([]*dto.ChatMemberDTO, error)
	IsUserInChat(ctx context.Context, userID, chatID string) (bool, error)
}

type WebsocketHandlers struct {
	chatUsecase ChatUsecase
	factory     *MessageFactory

	upgrader websocket.Upgrader
	conns    map[uuid.UUID][]*Client
	mu       sync.Mutex
}

func NewWebsockethandlers(chatUsecase ChatUsecase) *WebsocketHandlers {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	factoryMsg := NewMessageFactory()

	return &WebsocketHandlers{
		chatUsecase: chatUsecase,
		factory:     factoryMsg,
		upgrader:    upgrader,
		conns:       make(map[uuid.UUID][]*Client),
	}
}

func (h *WebsocketHandlers) HandleConnection(c *gin.Context) {
	userIdStr := c.GetString("user_id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("failed to get user_id: %w", err)})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("failed to upgrade: %w", err)})
		return
	}

	client := NewClient(userId, conn)

	h.addClient(client)

	conn.SetReadLimit(1048576)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	go h.handlePingPong(client)

	go h.readPump(client)
	go h.writePump(client)

	h.downloadMessages(client)
}

func (h *WebsocketHandlers) downloadMessages(client *Client) {
	messages := h.chatUsecase.GetMessagesFromStorage(context.Background(), client.UserID)
	if len(messages) == 0 {
		return
	}

	for _, msg := range messages {
		select {
		case client.Send <- msg:
			// отправлено
		default:
			go h.removeClient(client.UserID, client.SessionID)
		}
	}
}

func (h *WebsocketHandlers) addClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.conns[client.UserID] = append(h.conns[client.UserID], client)
}

func (h *WebsocketHandlers) removeClient(userId, sessionId uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients, exists := h.conns[userId]
	if !exists {
		return
	}

	for i, client := range clients {
		if client.SessionID == sessionId {
			client.close()
			h.conns[userId] = append(clients[:i], clients[i+1:]...)

			if len(h.conns[userId]) == 0 {
				delete(h.conns, userId)
			}
			return
		}
	}
}

func (h *WebsocketHandlers) handlePingPong(client *Client) {
	ticker := time.NewTicker(25 * time.Second)
	defer func() {
		ticker.Stop()
		h.removeClient(client.UserID, client.SessionID)
	}()

	for {
		select {
		case <-ticker.C:
			if client.getIsClosed() {
				return
			}

			err := client.Conn.WriteControl(websocket.PingMessage, []byte(time.Now().String()), time.Now().Add(10*time.Second))
			if err != nil {
				return
			}
		case <-client.ctx.Done():
			return
		}
	}
}

func (h *WebsocketHandlers) writePump(client *Client) {
	ticker := time.NewTicker(50 * time.Second)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			if !ok {
				if client.cancelFunc != nil {
					client.cancelFunc()
				}
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			n := len(client.Send)
			for i := 0; i < n; i++ {
				w.Write(<-client.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-client.ctx.Done():
			return
		}
	}
}

func (h *WebsocketHandlers) readPump(client *Client) {
	defer func() {
		if client.cancelFunc != nil {
			client.cancelFunc()
		}
	}()

	for {
		select {
		case <-client.ctx.Done():
			return
		default:
			client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

			messageType, message, err := client.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Websocket error: %v", err)
				}
				return
			}

			if messageType == websocket.TextMessage {
				h.handleIncomingMessage(client, message)
			}
		}
	}
}

func (h *WebsocketHandlers) handleIncomingMessage(client *Client, message []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var incoming dto.IncomingMessage
	if err := json.Unmarshal(message, &incoming); err != nil {
		h.sendError(client, dto.ErrInvalidMsgFormat, "Неверный формат сообщения")
		return
	}

	switch incoming.Type {
	case dto.TypeSendMessage:
		h.handleSendMessage(ctx, client, incoming.Payload, incoming.ChatID)
	default:
		h.sendError(client, dto.ErrInvalidMsgType, fmt.Sprintf("Неподдерживаемый тип сообщения: %s", incoming.Type))
	}
}

func (h *WebsocketHandlers) handleSendMessage(ctx context.Context, client *Client, payload json.RawMessage, chatId uuid.UUID) {
	var req dto.SendMessageIncPayload
	if err := json.Unmarshal(payload, &req); err != nil {
		h.sendError(client, dto.ErrInvalidPayload, "Неверный формат запроса")
		return
	}

	if chatId == uuid.Nil {
		h.sendError(client, dto.ErrDataIsEmpty, "ID сообщения не может быть пустым")
		return
	}
	if req.Content == "" {
		h.sendError(client, dto.ErrDataIsEmpty, "Содержание сообщения не может быть пустым")
		return
	}
	if req.Type == "" {
		req.Type = "text"
	}

	hasAccess, err := h.chatUsecase.IsUserInChat(ctx, client.UserID.String(), chatId.String())
	if err != nil || !hasAccess {
		h.sendError(client, dto.ErrAccessDenied, "Нет доступа к чату")
		return
	}

	msg := &dto.MessageDTO{
		ID:       uuid.New(),
		ChatID:   chatId,
		SenderID: client.UserID,
		Content:  req.Content,
		Type:     req.Type,
		ReplyTo:  req.ReplyTo,
		SentAt:   time.Now(),
	}

	savedId, err := h.chatUsecase.SendMessageToDb(ctx, msg)
	if err != nil {
		h.sendError(client, dto.ErrSaveFailed, "Не удалось сохранить сообщение")
		return
	}

	msg.ID = uuid.MustParse(savedId)

	outgoing := h.factory.NewOutgoingMessage(msg)

	h.broadcastToChat(ctx, msg.ChatID, outgoing)
}

func (h *WebsocketHandlers) sendError(client *Client, code dto.ErrorType, message string) {
	errorMsg := h.factory.NewError(code, message, "")

	data, err := json.Marshal(errorMsg)
	if err != nil {
		return
	}

	select {
	case client.Send <- data:
		// отправлено
	default:
		go h.removeClient(client.UserID, client.SessionID)
	}
}

func (h *WebsocketHandlers) broadcastToChat(ctx context.Context, chatID uuid.UUID, event *dto.OutgoingMessage) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	members, err := h.chatUsecase.GetChatMembers(ctx, chatID.String())
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	for _, member := range members {
		clients, exists := h.conns[member.UserID]
		if !exists {
			h.chatUsecase.SendMsgToStorage(ctx, member.UserID, data)
			continue
		}

		for _, client := range clients {
			select {
			case <-ctx.Done():
				// таймаут
			default:
				select {
				case client.Send <- data:
					// отправлено
				default:
					go h.removeClient(client.UserID, client.SessionID)
				}
			}
		}
	}
}
