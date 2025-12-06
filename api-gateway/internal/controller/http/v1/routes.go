package v1

import "github.com/I-Van-Radkov/corporate-messenger/api-gateway/internal/config"

type Route struct {
	Pattern string
	Target  string
	Public  bool
}

func LoadRoutes(cfg config.RoutesConfig) []Route {
	var routes = []Route{
		{Pattern: "/api/v1/admin/*", Target: cfg.AuthServicePath + "/admin", Public: false},
		{Pattern: "/api/v1/auth/*", Target: cfg.AuthServicePath + "/auth", Public: true},
		{Pattern: "/api/v1/directory/*", Target: cfg.DirectoryServicePath + "/directory", Public: false},
		{Pattern: "/api/v1/chats/*", Target: cfg.ChatServicePath + "/chats", Public: false},
	}

	return routes
}
