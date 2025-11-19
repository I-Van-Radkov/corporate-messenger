package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	srv *http.Server
}

func NewServer(port int, readTimeout, writeTimeout time.Duration) *Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Handler:      nil,
	}

	return &Server{
		srv: srv,
	}
}

func (s *Server) RegisterHandlers() {
	router := gin.New()

	api := router.Group("/api/v1")
	{
		//роуты для ендпоинтов
	}
}
