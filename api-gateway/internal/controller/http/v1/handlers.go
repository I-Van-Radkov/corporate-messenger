package v1

import (
	"net/http/httputil"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProxyProvider interface {
	Get(string) *httputil.ReverseProxy
}

type ProxyHandlers struct {
	factory ProxyProvider
}

func NewProxyHandlers(factory ProxyProvider) *ProxyHandlers {
	return &ProxyHandlers{
		factory: factory,
	}
}

func (h *ProxyHandlers) ProxyTo(target string) gin.HandlerFunc {
	revProxy := h.factory.Get(target)
	return func(c *gin.Context) {
		path := c.Param("proxyPath")
		if path != "" && !strings.HasPrefix(path, "/") {
			path = "/" + path
		}

		c.Request.URL.Path = path
		revProxy.ServeHTTP(c.Writer, c.Request)
	}
}
