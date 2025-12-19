package v1

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/I-Van-Radkov/corporate-messenger/identity-service/internal/adapter"
	"github.com/I-Van-Radkov/corporate-messenger/identity-service/internal/clients/directory"
	"github.com/I-Van-Radkov/corporate-messenger/identity-service/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	srv       *http.Server
	dirClient directory.Client
	db        *pgxpool.Pool
}

func NewServer(port int, readTimeout, writeTimeout time.Duration, dirCfg directory.DirectoryServiceConfig, db *pgxpool.Pool) *Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Handler:      nil,
	}

	dirClient := directory.NewHTTPClient(dirCfg)

	return &Server{
		srv:       srv,
		dirClient: dirClient,
		db:        db,
	}
}

func (s *Server) RegisterHandlers(authCfg usecase.AuthConfig) error {
	authrepo := adapter.NewAuthRepo(s.db)
	authUsecase := usecase.NewAuthUsecase(authrepo, s.dirClient, authCfg)

	authHandlers := NewAuthHandlers(authUsecase)

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	adminAvail := router.Group("/admin")
	adminAvail.Use(ExtractUserInfoMiddleware())
	adminAvail.Use(RequireAdminOnly())
	{
		adminAvail.POST("/accounts", authHandlers.CreateAccountHandler)           // Create
		adminAvail.GET("/accounts", authHandlers.ListAccountsHandler)             // List
		adminAvail.PUT("/accounts/:id", authHandlers.UpdateAccountHandler)        // Update
		adminAvail.DELETE("/accounts/:id", authHandlers.DeactivateAccountHandler) // Deactivate (soft delete)
	}
	userAvail := router.Group("/auth")
	{
		userAvail.POST("/login", authHandlers.LoginHandler)
		userAvail.POST("/introspect", authHandlers.IntrospectToken)
		// userAvail.POST("/logout")
	}

	s.srv.Handler = router

	return nil
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
