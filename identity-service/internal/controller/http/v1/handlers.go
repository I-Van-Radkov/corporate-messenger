package v1

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/I-Van-Radkov/corporate-messenger/identity-service/internal/dto"
	"github.com/I-Van-Radkov/corporate-messenger/identity-service/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthUsecase interface {
	CreateAccount(ctx context.Context, req *dto.CreateAccountRequest) (*dto.AccountResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (string, error)

	IntrospectToken(ctx context.Context, req *dto.IntrospectRequest) (*dto.IntrospectResponse, error)

	ListAccounts(ctx context.Context) ([]*dto.AccountResponse, error)
	UpdateAccount(ctx context.Context, accountID uuid.UUID, req *dto.UpdateAccountRequest) (*dto.AccountResponse, error)
	DeactivateAccount(ctx context.Context, accountID uuid.UUID) error
}

type AuthHandlers struct {
	authUsecase AuthUsecase
}

func NewAuthHandlers(authUsecase AuthUsecase) *AuthHandlers {
	return &AuthHandlers{
		authUsecase: authUsecase,
	}
}

func (h *AuthHandlers) CreateAccountHandler(c *gin.Context) {
	var req dto.CreateAccountRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authUsecase.CreateAccount(c.Request.Context(), &req)
	if err != nil {
		if err == errors.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandlers) LoginHandler(c *gin.Context) {
	log.Println("я в логин хендлере")

	var req dto.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(req.Email, req.Password)

	tokenString, err := h.authUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		if err == errors.ErrInvalidPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if err == errors.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "Bearer",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}

func (h *AuthHandlers) IntrospectToken(c *gin.Context) {
	log.Println("я в интроспект хендлере")
	var req dto.IntrospectRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authUsecase.IntrospectToken(c.Request.Context(), &req)
	if err != nil {
		if err == errors.ErrInvalidToken {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// В internal/controller/http/v1/handlers/auth_handlers.go добавь:

func (h *AuthHandlers) ListAccountsHandler(c *gin.Context) {
	accounts, err := h.authUsecase.ListAccounts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accounts": accounts, "total": len(accounts)})
}

func (h *AuthHandlers) UpdateAccountHandler(c *gin.Context) {
	accountIDStr := c.Param("id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	var req dto.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authUsecase.UpdateAccount(c.Request.Context(), accountID, &req)
	if err != nil {
		if err == errors.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandlers) DeactivateAccountHandler(c *gin.Context) {
	accountIDStr := c.Param("id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	if err := h.authUsecase.DeactivateAccount(c.Request.Context(), accountID); err != nil {
		if err == errors.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
