package identity

import "time"

type IdentityServiceConfig struct {
	URL             string        `env:"AUTH_SERVICE_PATH" env-default:"http://localhost:8081/auth/introspect"`
	Timeout         time.Duration `env:"DIRECTORY_SERVICE_TIMEOUT" env-default:"5s"`
	MaxIdleConns    int           `env:"DIRECTORY_MAX_IDLE_CONNS" env-default:"100"`
	MaxConnsPerHost int           `env:"DIRECTORY_MAX_CONNS_PER_HOST" env-default:"50"`
	IdleConnTimeout time.Duration `env:"DIRECTORY_IDLE_CONN_TIMEOUT" env-default:"90s"`
	RetryCount      int           `env:"DIRECTORY_RETRY_COUNT" env-default:"3"`
	CircuitBreaker  struct {
		FailureThreshold int           `env:"DIRECTORY_CB_FAILURE_THRESHOLD" env-default:"5"`
		SuccessThreshold int           `env:"DIRECTORY_CB_SUCCESS_THRESHOLD" env-default:"1"`
		Timeout          time.Duration `env:"DIRECTORY_CB_TIMEOUT" env-default:"30s"`
	}
}
