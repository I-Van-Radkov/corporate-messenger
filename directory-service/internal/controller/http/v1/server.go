package v1

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/I-Van-Radkov/corporate-messenger/directory-service/internal/adapter"
	"github.com/I-Van-Radkov/corporate-messenger/directory-service/internal/usecase"
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
	dirRepo := adapter.NewDirectoryRepo(s.db)
	dirUsecase := usecase.NewDirectoryUsecase(dirRepo)

	dirhandlers := NewDirectoryHandlers(dirUsecase)

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
	router.Use(ExtractUserInfoMiddleware())

	adminAvail := router.Group("/directory")
	adminAvail.Use(RequireAdminOnly())
	{
		adminAvail.POST("/departments", dirhandlers.CreateDepartment)
		adminAvail.GET("/departments", dirhandlers.GetDepartments)
		adminAvail.DELETE("/departments/:department_id", dirhandlers.RemoveDepartment)
		adminAvail.GET("/departments/:department_id/users", dirhandlers.GetDepartmentMembers)

		adminAvail.POST("/users", dirhandlers.CreateUser)
		adminAvail.GET("/users", dirhandlers.GetUsers)
		adminAvail.GET("/users/:user_id", dirhandlers.GetUser)
		adminAvail.DELETE("/users/:user_id", dirhandlers.RemoveUser)
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
