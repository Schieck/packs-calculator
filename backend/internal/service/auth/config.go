package service

import (
	"time"

	entity "github.com/Schieck/packs-calculator/internal/domain/entity/auth"
)

type DefaultTimeProvider struct{}

func (p *DefaultTimeProvider) Now() time.Time {
	return time.Now()
}

type DefaultAuthConfig struct {
	authSecret      string
	tokenExpiration time.Duration
	issuer          string
}

func NewDefaultAuthConfig(authSecret string, issuer string) entity.AuthConfig {
	return &DefaultAuthConfig{
		authSecret:      authSecret,
		tokenExpiration: 24 * time.Hour, // Default set to 24 hours
		issuer:          issuer,
	}
}

func (c *DefaultAuthConfig) GetAuthSecret() string {
	return c.authSecret
}

func (c *DefaultAuthConfig) GetTokenExpiration() time.Duration {
	return c.tokenExpiration
}

func (c *DefaultAuthConfig) GetIssuer() string {
	return c.issuer
}
