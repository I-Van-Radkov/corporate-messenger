package v1

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/I-Van-Radkov/corporate-messenger/api-gateway/internal/config"
	"github.com/I-Van-Radkov/corporate-messenger/api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
)

type Server struct {
	srv    *http.Server
	routes []Route
}

func NewServer(cfg *config.Config) *Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", cfg.Port),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Handler:      nil,
	}

	routes := LoadRoutes(cfg.RoutesConfig)

	return &Server{
		srv:    srv,
		routes: routes,
	}
}

func (s *Server) RegisterHandlers() error {
	router := gin.New()

	factory := proxy.NewFactory()
	proxyHandlers := NewProxyHandlers(factory)

	protected := router.Group("/")
	public := router.Group("/")

	for _, route := range s.routes {
		pattern := strings.Replace(route.Pattern, "*", "{proxyPath:*}", 1)

		group := public
		if !route.Public {
			group = protected
		}

		group.Any(pattern, proxyHandlers.ProxyTo(route.Target))
	}

	return nil
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
