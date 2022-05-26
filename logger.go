package main

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func NewLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}

type StructuredLogger struct {
	Logger *zap.SugaredLogger
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	l.Logger.With(
		"ts", time.Now().UTC().Format(time.RFC1123),
		"method", r.Method,
		"addr", r.RemoteAddr,
		"uri", fmt.Sprintf("%s%s", r.Host, r.RequestURI),
	).Info("request started")

	return l
}

func (l *StructuredLogger) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger = l.Logger.With(
		"status", http.StatusText(status),
		"bytes_length", bytes,
		"latency", float64(elapsed.Nanoseconds())/1000000.0,
	)

	l.Logger.Info("request complete")
}

func (l *StructuredLogger) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.With(
		"stack", string(stack),
		"panic", fmt.Sprintf("%+v", v),
	)
}
