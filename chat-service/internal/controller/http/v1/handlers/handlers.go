package handlers

import (
	"context"
	"net/http"

	"github.com/I-Van-Radkov/corporate-messenger/chat-service/internal/dto"
	"github.com/I-Van-Radkov/corporate-messenger/chat-service/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatUsecase interface {
	CreateChat(ctx context.Context, req *dto.CreateChatRequest, creatorID uuid.UUID) (*dto.CreateChatResponse, error)
	RemoveChat(ctx context.Context, chatID uuid.UUID) error
	GetUserChats(ctx context.Context, userID uuid.UUID) (*dto.GetUserChatsResponse, error)
	GetChatMessages(ctx context.Context, chatID, userID uuid.UUID, limit int, before *uuid.UUID) (*dto.GetMessagesResponse, error)
	GetChatMembers(ctx context.Context, chatID uuid.UUID) (*dto.ChatMembers, error)

	AddMembers(ctx context.Context, chatID uuid.UUID, userIDs []uuid.UUID, adderID uuid.UUID) error
	RemoveMember(ctx context.Context, chatID, userID, removerID uuid.UUID) error
	ChangeMemberRole(ctx context.Context, chatID, userID uuid.UUID, role dto.MemberRole, changerID uuid.UUID) error
}

type Chathandlers struct {
	chatusecase ChatUsecase
}

func NewChatHandlers(chatusecase ChatUsecase) *Chathandlers {
	return &Chathandlers{
		chatusecase: chatusecase,
	}
}

func (h *Chathandlers) CreateChat(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var req dto.CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.chatusecase.CreateChat(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Chathandlers) RemoveChat(c *gin.Context) {
	chatIdStr := c.Param("chat_id")
	chatId, err := uuid.Parse(chatIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.chatusecase.RemoveChat(c.Request.Context(), chatId)
	if err != nil {
		if err == errors.ErrChatNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
func (h *Chathandlers) GetUserChats(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))

	resp, err := h.chatusecase.GetUserChats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Chathandlers) GetChatMessages(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	var req dto.GetMessagesRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid params"})
		return
	}

	if req.Limit == 0 {
		req.Limit = 50
	}

	resp, err := h.chatusecase.GetChatMessages(c.Request.Context(), req.ChatID, userID, req.Limit, req.Before)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Chathandlers) AddMembers(c *gin.Context) {
	adderID, _ := uuid.Parse(c.GetString("user_id"))
	chatID, _ := uuid.Parse(c.Param("chat_id"))

	var req dto.AddMembersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_ids"})
		return
	}

	if err := h.chatusecase.AddMembers(c.Request.Context(), chatID, req.UserIDs, adderID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Chathandlers) RemoveMember(c *gin.Context) {
	removerID, _ := uuid.Parse(c.GetString("user_id"))
	chatID, _ := uuid.Parse(c.Param("chat_id"))
	userID, _ := uuid.Parse(c.Param("user_id"))

	if err := h.chatusecase.RemoveMember(c.Request.Context(), chatID, userID, removerID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Chathandlers) ChangeMemberRole(c *gin.Context) {
	changerID, _ := uuid.Parse(c.GetString("user_id"))
	chatID, _ := uuid.Parse(c.Param("chat_id"))
	userID, _ := uuid.Parse(c.Param("user_id"))

	var req dto.ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	if err := h.chatusecase.ChangeMemberRole(c.Request.Context(), chatID, userID, req.Role, changerID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Chathandlers) GetChatMembers(c *gin.Context) {
	chatID, _ := uuid.Parse(c.Param("chat_id"))

	resp, err := h.chatusecase.GetChatMembers(c.Request.Context(), chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
