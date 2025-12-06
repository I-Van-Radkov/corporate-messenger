package identity

import "time"

type IdentityServiceConfig struct {
	URL             string        `env:"DIRECTORY_SERVICE_PATH,required"`
	Timeout         time.Duration `env:"DIRECTORY_SERVICE_TIMEOUT" envDefault:"5s"`
	MaxIdleConns    int           `env:"DIRECTORY_MAX_IDLE_CONNS" envDefault:"100"`
	MaxConnsPerHost int           `env:"DIRECTORY_MAX_CONNS_PER_HOST" envDefault:"50"`
	IdleConnTimeout time.Duration `env:"DIRECTORY_IDLE_CONN_TIMEOUT" envDefault:"90s"`
	RetryCount      int           `env:"DIRECTORY_RETRY_COUNT" envDefault:"3"`
	CircuitBreaker  struct {
		FailureThreshold int           `env:"DIRECTORY_CB_FAILURE_THRESHOLD" envDefault:"5"`
		SuccessThreshold int           `env:"DIRECTORY_CB_SUCCESS_THRESHOLD" envDefault:"1"`
		Timeout          time.Duration `env:"DIRECTORY_CB_TIMEOUT" envDefault:"30s"`
	}
}
