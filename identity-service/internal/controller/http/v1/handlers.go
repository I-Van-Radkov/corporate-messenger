package v1

import (
	"context"

	"github.com/I-Van-Radkov/corporate-messenger/identity-service/internal/dto"
	"github.com/gin-gonic/gin"
)

type AuthUsecase interface {
	CreateAccount(ctx context.Context, req dto.CreateAccountRequest) (*dto.AccountResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (string, error)
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

}

func (h *AuthHandlers) LoginHandler(c *gin.Context) {

}

func (h *AuthHandlers) IntrospectToken(c *gin.Context) {

}
