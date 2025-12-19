package v1

import (
	"log"
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
		// Получаем оригинальный путь
		originalPath := c.Request.URL.Path

		// Определяем какой сервис и какой путь удалять
		var prefixToRemove string

		switch {
		case strings.HasPrefix(originalPath, "/api/v1/admin"):
			prefixToRemove = "/api/v1/admin"
		case strings.HasPrefix(originalPath, "/api/v1/auth"):
			prefixToRemove = "/api/v1/auth"
		case strings.HasPrefix(originalPath, "/api/v1/directory"):
			prefixToRemove = "/api/v1/directory"
		case strings.HasPrefix(originalPath, "/api/v1/chats"):
			prefixToRemove = "/api/v1/chats"
		default:
			prefixToRemove = "/api/v1"
		}

		// Удаляем префикс
		newPath := strings.TrimPrefix(originalPath, prefixToRemove)

		// Если путь пустой, ставим "/"
		if newPath == "" {
			newPath = "/"
		}

		// Устанавливаем новый путь
		c.Request.URL.Path = newPath

		// Логируем для отладки
		log.Printf("[Proxy] %s -> %s%s", originalPath, target, newPath)

		// Проксируем запрос
		revProxy.ServeHTTP(c.Writer, c.Request)
	}
}
