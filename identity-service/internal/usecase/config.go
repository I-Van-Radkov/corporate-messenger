package usecase

import "time"

type AuthConfig struct {
	JwtSecret    string        `env:"JWT_SECRET,required"`
	JwtExpiresIn time.Duration `env:"JWT_EXPIRES_IN" env-default:"15m"`
}
