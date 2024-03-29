package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func RequestLogger() func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&requestLogger{log.StandardLogger()})
}

type requestLogger struct {
	Logger *log.Logger
}

func (l *requestLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &requestLoggerEntry{Logger: log.NewEntry(l.Logger)}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	logFields := log.Fields{
		"ts":          time.Now().UTC().Format(time.RFC1123),
		"http_scheme": scheme,
		"http_proto":  r.Proto,
		"http_method": r.Method,
		"remote_addr": r.RemoteAddr,
		"user_agent":  r.UserAgent(),
		"uri":         fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI),
	}

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	entry.Logger = entry.Logger.WithFields(logFields)

	entry.Logger.Infoln("request started")

	return entry
}

type requestLoggerEntry struct {
	Logger log.FieldLogger
}

func (l *requestLoggerEntry) Write(status, bytes int, elapsed time.Duration) {
	l.Logger = l.Logger.WithFields(log.Fields{
		"resp_status":       status,
		"resp_bytes_length": bytes,
		"resp_elapsed_ms":   float64(elapsed.Nanoseconds()) / 1000000.0,
	})

	l.Logger.Infoln("request complete")
}

func (l *requestLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(log.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

func GetLogEntry(r *http.Request) (log.FieldLogger, error) {
	entry, ok := middleware.GetLogEntry(r).(*requestLoggerEntry)
	if !ok {
		return nil, errors.New("Unable to cast to 'requestLoggerEntry'")
	}
	return entry.Logger, nil
}
