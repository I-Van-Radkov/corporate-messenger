package v1

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/I-Van-Radkov/corporate-messenger/chat-service/internal/adapter"
	"github.com/I-Van-Radkov/corporate-messenger/chat-service/internal/controller/http/v1/handlers"
	"github.com/I-Van-Radkov/corporate-messenger/chat-service/internal/controller/websocket"
	"github.com/I-Van-Radkov/corporate-messenger/chat-service/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	srv *http.Server
	db  *pgxpool.Pool
}

func NewServer(port int, readTimeout, writeTimeout time.Duration, db *pgxpool.Pool) *Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Handler:      nil,
	}

	return &Server{
		srv: srv,
		db:  db,
	}
}

func (s *Server) RegisterHandlers() error {
	chatRepo := adapter.NewChatRepo(s.db)
	chatUsecase := usecase.NewChatUsecase(chatRepo)

	wsHandlers := websocket.NewWebsockethandlers(chatUsecase)
	chatHandlers := handlers.NewChatHandlers(chatUsecase)

	router := gin.Default()
	router.Use(ExtractUserInfoMiddleware())

	router.GET("/ws", wsHandlers.HandleConnection)

	chats := router.Group("/chats")
	{
		// Все пользователи
		chats.GET("", chatHandlers.GetUserChats)
		chats.POST("", chatHandlers.CreateChat)
		chats.GET("/:chat_id/members", chatHandlers.GetChatMembers)
		router.GET("/:chat_id") // ДОПИСАТЬ ЗАГРУЗКУ СООБЩЕНИЙ ЧАТА

		// Только staff (admin/moderator/support)
		staffOnly := chats.Group("")
		staffOnly.Use(RequireStaff())
		{
			staffOnly.DELETE("/:chat_id/members/:user_id", chatHandlers.RemoveMember)
			staffOnly.POST("/:chat_id/members/:user_id/role", chatHandlers.ChangeMemberRole)
			staffOnly.POST("/:chat_id/members", chatHandlers.AddMembers)
			staffOnly.DELETE("/:chat_id", chatHandlers.RemoveChat)
		}
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
