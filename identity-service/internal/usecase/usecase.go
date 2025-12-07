package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/I-Van-Radkov/corporate-messenger/identity-service/internal/clients/directory"
	"github.com/I-Van-Radkov/corporate-messenger/identity-service/internal/dto"
	errModels "github.com/I-Van-Radkov/corporate-messenger/identity-service/internal/errors"
	"github.com/I-Van-Radkov/corporate-messenger/identity-service/internal/models"
	"github.com/I-Van-Radkov/corporate-messenger/identity-service/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthRepo interface {
	Create(ctx context.Context, account *models.Account) error
	FindByID(ctx context.Context, accountID string) (*models.Account, error)
	FindByEmail(ctx context.Context, email string) (*models.Account, error)
	FindByUserID(ctx context.Context, userID string) (*models.Account, error)
}

type AuthUsecase struct {
	authrepo  AuthRepo
	dirClient directory.Client
	authCfg   AuthConfig
}

func NewAuthUsecase(authrepo AuthRepo, dirClient directory.Client, authCfg AuthConfig) *AuthUsecase {
	return &AuthUsecase{
		authrepo:  authrepo,
		dirClient: dirClient,
		authCfg:   authCfg,
	}
}

func (u *AuthUsecase) CreateAccount(ctx context.Context, input *dto.CreateAccountRequest) (*dto.AccountResponse, error) {
	// Проверка на существование работника в directory-db
	exists, err := u.dirClient.UserExists(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from directory-service: %w", err)
	}
	if !exists {
		return nil, errModels.ErrUserNotFound
	}

	// Хеширование пароля
	passwordHash, err := utils.HashPasswordBase64(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Сохранение в БД
	account := &models.Account{
		AccountID:    uuid.New(),
		UserID:       uuid.MustParse(input.UserID),
		Email:        input.Email,
		PasswordHash: passwordHash,
		Role:         models.AccountRole(input.Role),
		IsActive:     true,
		LastLogin:    nil,
		CreatedAt:    time.Now(),
	}

	err = u.authrepo.Create(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("failed to save account to db: %w", err)
	}

	accountResponse := &dto.AccountResponse{
		AccountID: account.AccountID.String(),
		UserID:    account.UserID.String(),
		Email:     account.Email,
		Role:      dto.AccountRole(account.Role),
		IsActive:  account.IsActive,
		LastLogin: account.LastLogin,
		CreatedAt: account.CreatedAt,
	}

	return accountResponse, nil
}

func (u *AuthUsecase) Login(ctx context.Context, input *dto.LoginRequest) (string, error) {
	// Проверка на наличие аккаунта
	account, err := u.authrepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errModels.ErrUserNotFound
		}

		return "", fmt.Errorf("failed to find account by id from db: %w", err)
	}

	// Подтверждение пароля
	ok, err := utils.VerifyPassword(input.Password, account.PasswordHash)
	if err != nil {
		return "", fmt.Errorf("failed to verify password: %w", err)
	}
	if !ok {
		return "", errModels.ErrInvalidPassword
	}

	// Генерация токена
	tokenString, err := utils.SignToken(account.UserID, account.Role, u.authCfg.JwtSecret, u.authCfg.JwtExpiresIn)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	// Возврат токена
	return tokenString, nil
}

func (u *AuthUsecase) IntrospectToken(ctx context.Context, input *dto.IntrospectRequest) (*dto.IntrospectResponse, error) {
	parsedToken, err := jwt.Parse(input.Token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signed method: %v", t.Header["alg"])
		}

		return []byte(u.authCfg.JwtSecret), nil
	})

	output := &dto.IntrospectResponse{
		Active: false,
	}

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return output, nil
		}
		return nil, errModels.ErrInvalidToken
	}

	if !parsedToken.Valid {
		return nil, errModels.ErrInvalidToken
	}

	output.Active = true
	return output, nil
}
