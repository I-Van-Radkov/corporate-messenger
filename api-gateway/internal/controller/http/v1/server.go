package v1

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/I-Van-Radkov/corporate-messenger/api-gateway/internal/clients/identity"
	"github.com/I-Van-Radkov/corporate-messenger/api-gateway/internal/config"
	"github.com/I-Van-Radkov/corporate-messenger/api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
)

type Server struct {
	srv        *http.Server
	authClient identity.Client
	routes     []Route
}

func NewServer(cfg *config.Config, authClient identity.Client) *Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", cfg.Port),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Handler:      nil,
	}

	routes := LoadRoutes(cfg.RoutesConfig)

	return &Server{
		srv:        srv,
		authClient: authClient,
		routes:     routes,
	}
}

func (s *Server) RegisterHandlers() error {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS,PATCH")
		c.Header("Access-Control-Allow-Headers",
			"Content-Type,Authorization,Accept,Origin,X-Requested-With,X-User-ID,X-User-Role")
		c.Header("Access-Control-Expose-Headers",
			"Content-Length,Content-Range,Authorization,X-User-ID,X-User-Role")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	factory := proxy.NewFactory()
	proxyHandlers := NewProxyHandlers(factory)

	protected := router.Group("/")
	protected.Use(AuthMiddleware(s.authClient))

	public := router.Group("/")

	for _, route := range s.routes {
		// Исправление: используем /*catchall для wildcard
		pattern := strings.TrimSuffix(route.Pattern, "*") + "/*catchall"

		group := public
		if !route.Public {
			group = protected
		}

		group.Any(pattern, proxyHandlers.ProxyTo(route.Target))
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
