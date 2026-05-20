package slog

import (
	"context"
	stdslog "log/slog"
	"os"
	"sync"

	domainlogger "github.com/pascalallen/go-clean-arch/internal/app/domain/logger"
)

type slogLogger struct {
	l   *stdslog.Logger
	ctx context.Context
	mu  sync.RWMutex
}

// New creates a slog-backed Logger. In production (APP_ENV=production), it
// uses JSON output at Info level. Otherwise, it uses a human-friendly text
// handler at Debug level.
func New() domainlogger.Logger {
	handler := defaultHandler()
	return &slogLogger{l: stdslog.New(handler)}
}

func defaultHandler() stdslog.Handler {
	env := os.Getenv("APP_ENV")
	switch env {
	case "production", "prod":
		return stdslog.NewJSONHandler(os.Stdout, &stdslog.HandlerOptions{Level: stdslog.LevelInfo})
	default:
		return stdslog.NewTextHandler(os.Stdout, &stdslog.HandlerOptions{Level: stdslog.LevelDebug})
	}
}

func (s *slogLogger) Debug(msg string, keyVals ...any) {
	s.mu.RLock()
	if s.ctx != nil {
		s.l.Log(s.ctx, stdslog.LevelDebug, msg, keyVals...)
	} else {
		s.l.Debug(msg, keyVals...)
	}
	s.mu.RUnlock()
}

func (s *slogLogger) Info(msg string, keyVals ...any) {
	s.mu.RLock()
	if s.ctx != nil {
		s.l.Log(s.ctx, stdslog.LevelInfo, msg, keyVals...)
	} else {
		s.l.Info(msg, keyVals...)
	}
	s.mu.RUnlock()
}

func (s *slogLogger) Warn(msg string, keyVals ...any) {
	s.mu.RLock()
	if s.ctx != nil {
		s.l.Log(s.ctx, stdslog.LevelWarn, msg, keyVals...)
	} else {
		s.l.Warn(msg, keyVals...)
	}
	s.mu.RUnlock()
}

func (s *slogLogger) Error(msg string, keyVals ...any) {
	s.mu.RLock()
	if s.ctx != nil {
		s.l.Log(s.ctx, stdslog.LevelError, msg, keyVals...)
	} else {
		s.l.Error(msg, keyVals...)
	}
	s.mu.RUnlock()
}

func (s *slogLogger) With(keyVals ...any) domainlogger.Logger {
	s.mu.RLock()
	nl := s.l.With(keyVals...)
	s.mu.RUnlock()
	return &slogLogger{l: nl, ctx: s.ctx}
}

func (s *slogLogger) WithContext(ctx context.Context) domainlogger.Logger {
	s.mu.RLock()
	base := s.l
	s.mu.RUnlock()
	return &slogLogger{l: base, ctx: ctx}
}
