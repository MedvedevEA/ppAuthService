package logger

import (
	"log"
	"log/slog"
	"os"
)

func MustNew(env string) *slog.Logger {
	var lg *slog.Logger
	switch env {
	case "local":
		log.Printf("logger is configured for deployment environment 'local'\n")
		lg = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "dev":
		log.Printf("logger is configured for deployment environment 'dev'\n")
		lg = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log.Printf("logger is configured for deployment environment 'prod'\n")
		lg = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log.Fatalf("failed to configured logger: application deployment environment is not defined\n")
	}

	return lg
}
