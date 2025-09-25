package log

import (
	"log/slog"
	"os"
)

var (
	InfoLogger    *slog.Logger
	WarningLogger *slog.Logger
	ErrorLogger   *slog.Logger
)

func init() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}

	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{AddSource: true})

	baseLogger := slog.New(handler)

	InfoLogger = baseLogger.With("level", slog.LevelInfo)
	WarningLogger = baseLogger.With("level", slog.LevelWarn)
	ErrorLogger = baseLogger.With("level", slog.LevelError)

	slog.SetDefault(baseLogger)
}
