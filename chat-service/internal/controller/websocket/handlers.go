package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/I-Van-Radkov/corporate-messenger/chat-service/internal/dto"
	"github.com/I-Van-Radkov/corporate-messenger/chat-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	Conn      *websocket.Conn
	Send      chan []byte

	IsClosed bool
	mu       sync.Mutex
}

type ChatUsecase interface {
	SendMessage(ctx context.Context, msg *dto.MessageDTO) (string, error)
	EditMessage(ctx context.Context, messageID, userID, content string) error
	DeleteMessage(ctx context.Context, messageID, userID string) error
	MarkAsRead(ctx context.Context, userID, chatID, messageID string) error
	GetChatMembers(ctx context.Context, chatID string) ([]*dto.ChatMemberDTO, error)
	IsUserInChat(ctx context.Context, userID, chatID string) (bool, error)
}

type WebsocketHandlers struct {
	chatUsecase ChatUsecase
	upgrader    websocket.Upgrader
	conns       map[uuid.UUID][]*Client
	mu          sync.Mutex
}

func NewWebsockethandlers(chatUsecase ChatUsecase) *WebsocketHandlers {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return &WebsocketHandlers{
		chatUsecase: chatUsecase,
		upgrader:    upgrader,
		conns:       make(map[uuid.UUID][]*Client),
	}
}

func (h *WebsocketHandlers) HandleConnection(c *gin.Context) {
	ctx := context.Background()

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

	sessionId := uuid.New()

	client := &Client{
		UserID:    userId,
		SessionID: sessionId,
		Conn:      conn,
		Send:      make(chan []byte),
	}

	h.addClient(client)

	conn.SetReadLimit(1048576)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	go h.handlePingPong(ctx, client)

	go h.readPump(ctx, client)
	go h.writePump(ctx, client)
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

func (c *Client) close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.IsClosed {
		c.Conn.Close()
		close(c.Send)
		c.IsClosed = true
	}
}

func (h *WebsocketHandlers) handlePingPong(ctx context.Context, client *Client) {
	ticker := time.NewTicker(25 * time.Second)
	defer func() {
		ticker.Stop()
		h.removeClient(client.UserID, client.SessionID)
	}()

	for {
		select {
		case <-ticker.C:
			err := client.Conn.WriteControl(websocket.PingMessage, []byte(time.Now().String()), time.Now().Add(10*time.Second))
			if err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (h *WebsocketHandlers) writePump(ctx context.Context, client *Client) {
	ticker := time.NewTicker(50 * time.Second)
	defer func() {
		ticker.Stop()
		client.close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
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
		case <-ctx.Done():
			return
		}
	}
}

func (h *WebsocketHandlers) readPump(ctx context.Context, client *Client) {
	defer func() {
		h.removeClient(client.UserID, client.SessionID)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			messageType, message, err := client.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Websocket error: %v", err)
				}
				return
			}

			if messageType == websocket.TextMessage {
				h.handleMessage(client, message)
			}
		}
	}
}

func (h *WebsocketHandlers) handleMessage(client *Client, message []byte) {
	var wsMessage WSMessage

	if err := json.Unmarshal(message, &wsMessage); err != nil {
		h.sendError(client, "invalid message format")
		return
	}

	log.Printf("Received message type: %s from user: %s", wsMessage.Type, client.UserID)

	switch wsMessage.Type {
	case "send_message":
		h.handleSendMessage(client, wsMessage.Payload, wsMessage.ChatID)
	// case "edit_message":
	// 	h.handleEditMessage(client, wsMessage.Payload)
	// case "delete_message":
	// 	h.handleDeleteMessage(client, wsMessage.Payload)
	// case "mark_as_read":
	// 	h.handleMarkAsRead(client, wsMessage.Payload)
	// case "typing":
	// 	h.handleTyping(client, wsMessage.Payload)
	default:
		h.sendError(client, "unknown message type")
	}
}

func (h *WebsocketHandlers) handleSendMessage(client *Client, payload json.RawMessage, chatID string) {
	if ok, err := h.chatUsecase.IsUserInChat(context.Background(), client.UserID.String(), chatID); err != nil || !ok {
		h.sendError(client, "access denied to chat")
		return
	}

	var sendPayload WSSendMessagePayload
	if err := json.Unmarshal(payload, &sendPayload); err != nil {
		h.sendError(client, "invalid send message payload")
		return
	}

	message := &dto.MessageDTO{
		ID:       uuid.New().String(),
		ChatID:   chatID,
		SenderID: client.UserID.String(),
		Content:  sendPayload.Content,
		Type:     models.MessageType(sendPayload.Type),
		ReplyTo:  &sendPayload.ReplyTo,
		SentAt:   time.Now(),
	}

	savedID, err := h.chatUsecase.SendMessage(context.Background(), message)
	if err != nil {
		h.sendError(client, "failed to save message")
		return
	}

	// Добавляем ID сохраненного сообщения
	message.ID = savedID

	// Рассылаем участникам чата
	h.broadcastToChat(chatID, WSEventMessage{
		Event:     "message_sent",
		Data:      message,
		Timestamp: time.Now(),
	})
}

func (h *WebsocketHandlers) sendError(client *Client, message string) {
	errorMsg := map[string]interface{}{
		"type": "error",
		"payload": map[string]string{
			"message": message,
		},
		"timestamp": time.Now().Unix(),
	}

	if data, err := json.Marshal(errorMsg); err == nil {
		client.Send <- data
	}
}

func (h *WebsocketHandlers) broadcastToChat(chatID string, event WSEventMessage) {
	// TODO: Реализовать логику получения участников чата из БД
	members, err := h.chatUsecase.GetChatMembers(context.Background(), chatID)
	if err != nil {
		log.Printf("Error getting members: %v", err)
		return
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling broadcast message: %v", err)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for _, member := range members {
		userID, err := uuid.Parse(member.UserID)
		if err != nil {
			continue
		}

		if clients, exists := h.conns[userID]; exists {
			for _, client := range clients {
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
