package config

import (
	dom "payment-sse/internal/domain"
	"payment-sse/internal/protocol/http"
	postgres "payment-sse/internal/repo/postgres/provider"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Http http.Config `envPrefix:"HTTP_"`

	PG postgres.Config `envPrefix:"PG_"`

	Dom dom.Config `envPrefix:"DOMAIN_"`
}

func NewConfig() (Config, error) {
	return env.ParseAsWithOptions[Config](env.Options{
		Prefix: "PAYMENT_SSE_",
	})
}
