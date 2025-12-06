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
	router := gin.New()

	factory := proxy.NewFactory()
	proxyHandlers := NewProxyHandlers(factory)

	protected := router.Group("/")
	protected.Use(AuthMiddleware(s.authClient))

	public := router.Group("/")
	for _, route := range s.routes {
		pattern := strings.Replace(route.Pattern, "*", "{proxyPath:*}", 1)

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
