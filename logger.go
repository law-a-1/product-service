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

type SugaredRequestLogger struct {
	Logger  *zap.SugaredLogger
	request *http.Request
}

func (l *SugaredRequestLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	l.request = r

	return l
}

func (l *SugaredRequestLogger) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger.Infof(
		"%s %s%s - %s %dB in %f",
		l.request.Method, l.request.Host, l.request.RequestURI, http.StatusText(status), bytes,
		float64(elapsed.Nanoseconds())/1000000.0)
}

func (l *SugaredRequestLogger) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.With(
		"stack", string(stack),
		"panic", fmt.Sprintf("%+v", v),
	)
}
