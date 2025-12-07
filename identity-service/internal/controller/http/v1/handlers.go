package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/I-Van-Radkov/corporate-messenger/identity-service/internal/dto"
	"github.com/I-Van-Radkov/corporate-messenger/identity-service/internal/errors"
	"github.com/gin-gonic/gin"
)

type AuthUsecase interface {
	CreateAccount(ctx context.Context, req *dto.CreateAccountRequest) (*dto.AccountResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (string, error)

	IntrospectToken(ctx context.Context, req *dto.IntrospectRequest) (*dto.IntrospectResponse, error)
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
	var req dto.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
		Secure:   gin.Mode() == gin.ReleaseMode,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}

func (h *AuthHandlers) IntrospectToken(c *gin.Context) {
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
