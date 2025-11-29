package v1

import "github.com/I-Van-Radkov/corporate-messenger/api-gateway/internal/config"

type Route struct {
	Pattern string
	Target  string
	Public  bool
}

func LoadRoutes(cfg config.RoutesConfig) []Route {
	var routes = []Route{
		{Pattern: "/api/v1/auth/*", Target: cfg.AuthServicePath, Public: true},
	}

	return routes
}
