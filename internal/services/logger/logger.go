package logger

import (
	"context"
	"log/slog"
	"os"
)

var Log *slog.Logger

func Init(level string) {
	// Устанавливаем уровень логирования
	var slogLevel slog.Level
	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	// Создаем логгер с JSON-форматом вывода
	Log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slogLevel,
	}))

	// Устанавливаем как логгер по умолчанию
	slog.SetDefault(Log)
}

// Helper методы для удобного использования

func Debug(msg string, args ...any) {
	Log.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	Log.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Log.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Log.Error(msg, args...)
}

func LogSlice[T any](key string, slice []*T) {
	if len(slice) == 0 {
		Log.Info("Empty slice", "key", key)
		return
	}

	attrs := make([]slog.Attr, len(slice))
	for i, v := range slice {
		attrs[i] = slog.Any(key, v)
	}

	Log.LogAttrs(context.TODO(), slog.LevelInfo, "Slice contents", attrs...)
}
