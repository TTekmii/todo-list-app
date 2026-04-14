// internal/lib/logger/logger.go
package logger

import (
	"log/slog"
	"os"

	"github.com/TTekmii/todo-list-app/internal/lib/logger/handlers/slogdiscard"
	"github.com/TTekmii/todo-list-app/internal/lib/logger/handlers/slogpretty"
)

type Config struct {
	Level  string `mapstructure:"level" env:"LOG_LEVEL" default:"info"`    // debug, info, warn, error
	Format string `mapstructure:"format" env:"LOG_FORMAT" default:"json"`  // json, text, pretty
	Env    string `mapstructure:"env" env:"APP_ENV" default:"development"` // development, staging, production
}

func New(cfg Config) *slog.Logger {
	var handler slog.Handler

	switch {
	case cfg.Env == "production" || cfg.Format == "json":
		level := parseLevel(cfg.Level)
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})

	case cfg.Format == "pretty":
		level := parseLevel(cfg.Level)
		opts := &slog.HandlerOptions{Level: level}
		handler = slogpretty.PrettyHandlerOptions{SlogOpts: opts}.NewPrettyHandler(os.Stdout)

	case cfg.Format == "text":
		level := parseLevel(cfg.Level)
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})

	case cfg.Level == "discard" || cfg.Env == "test":
		handler = slogdiscard.NewDiscardHandler()

	default:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	return slog.New(handler)
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
