package handlers

import (
	"context"

	"github.com/I-Van-Radkov/corporate-messenger/chat-service/internal/dto"
	"github.com/gin-gonic/gin"
)

type ChatUsecase interface {
	// Создание чата (private / group / department)
	CreateChat(ctx context.Context, input *dto.CreateChatRequest) (string, error)

	// Получить список чатов пользователя (с последним сообщением и unread count — можно в отдельной DTO)
	GetUserChats(ctx context.Context, userID string, limit, offset int) ([]*dto.ChatPreview, error)

	// Загрузить историю чата (с пагинацией)
	//GetMessages(ctx context.Context, params models.GetMessagesParams) (*models.MessagesPage, error)

	// Отметить сообщения как прочитанные до определённого (обычно вызывается при открытии чата)
	//MarkChatAsRead(ctx context.Context, userID, chatID, lastReadMessageID string) error

	// Добавить/удалить участников, сменить роли и т.д.
	AddMembers(ctx context.Context, chatID string, userIDs []string, adderID string) error
	RemoveMember(ctx context.Context, chatID, userID, removerID string) error
	ChangeMemberRole(ctx context.Context, chatID, userID string, newRole string, changerID string) error

	// Получить участников чата
	GetChatMembers(ctx context.Context, chatID string) ([]*dto.ChatMemberDTO, error)
}

type Chathandlers struct {
	chatusecase ChatUsecase
}

func NewChatHandlers(chatusecase ChatUsecase) *Chathandlers {
	return &Chathandlers{
		chatusecase: chatusecase,
	}
}

func (h *Chathandlers) CreateChat(c *gin.Context)
func (h *Chathandlers) RemoveChat(c *gin.Context)
func (h *Chathandlers) GetUserChats(c *gin.Context)
func (h *Chathandlers) AddMembers(c *gin.Context)
func (h *Chathandlers) RemoveMember(c *gin.Context)
func (h *Chathandlers) ChangeMemberRole(c *gin.Context)
func (h *Chathandlers) GetChatMembers(c *gin.Context)
