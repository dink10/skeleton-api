package server

import (
    "github.com/sirupsen/logrus"
    "bitbucket.org/gismart/{{Name}}/logger"
    "net/http"
    "github.com/go-chi/chi/middleware"
    "time"
    "fmt"
)

func RequestLogger() func(next http.Handler) http.Handler {
    return middleware.RequestLogger(&requestLogger{logger.Logger})
}

type requestLogger struct {
    Logger *logrus.Logger
}

func (l *requestLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
    entry := &requestLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
    logFields := logrus.Fields{}

    logFields["ts"] = time.Now().UTC().Format(time.RFC1123)

    if reqID := middleware.GetReqID(r.Context()); reqID != "" {
        logFields["req_id"] = reqID
    }

    scheme := "http"
    if r.TLS != nil {
        scheme = "https"
    }
    logFields["http_scheme"] = scheme
    logFields["http_proto"] = r.Proto
    logFields["http_method"] = r.Method

    logFields["remote_addr"] = r.RemoteAddr
    logFields["user_agent"] = r.UserAgent()

    logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

    entry.Logger = entry.Logger.WithFields(logFields)

    entry.Logger.Infoln("request started")

    return entry
}

type requestLoggerEntry struct {
    Logger logrus.FieldLogger
}

func (l *requestLoggerEntry) Write(status, bytes int, elapsed time.Duration) {
    l.Logger = l.Logger.WithFields(logrus.Fields{
        "resp_status": status, "resp_bytes_length": bytes,
        "resp_elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
    })

    l.Logger.Infoln("request complete")
}

func (l *requestLoggerEntry) Panic(v interface{}, stack []byte) {
    l.Logger = l.Logger.WithFields(logrus.Fields{
        "stack": string(stack),
        "panic": fmt.Sprintf("%+v", v),
    })
}

func GetLogEntry(r *http.Request) logrus.FieldLogger {
    entry := middleware.GetLogEntry(r).(*requestLoggerEntry)
    return entry.Logger
}
