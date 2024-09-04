package dom

import (
	"fmt"
	"time"
)

type Config struct {
	PaymentConfirmationIn       time.Duration `env:"PAYMENT_CONFIRMATION_IN,notEmpty"`
	InactivityOrderEventTimeout time.Duration `env:"INACTIVITY_ORDER_EVENT_TIMEOUT,notEmpty"`
	LogLevel                    LogLevel      `env:"LOG_LEVEL,notEmpty"`
}

type LogLevel string

const (
	Debug LogLevel = "DEBUG"
	Info  LogLevel = "INFO"
)

func LogLevelFromString(s string) (LogLevel, error) {
	switch s {
	case "DEBUG":
		return Debug, nil
	case "INFO":
		return Info, nil
	}

	return "", fmt.Errorf("unkwnown log level provided: %s", s)
}
