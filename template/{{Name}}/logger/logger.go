package logger

import (
    sentry "github.com/evalphobia/logrus_sentry"
    "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"fmt"
)

// InitLogger initialises logger
func InitLogger(logLevel, sentryDSN string) {
	log.SetFormatter(&log.JSONFormatter{})

	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse log level: %s, %v", logLevel, err))
	}
	log.SetLevel(lvl)

	if sentryDSN != "" {
		hook, err := sentry.NewSentryHook(sentryDSN, []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
		})
		if err != nil {
			log.Error(errors.Wrap(err, "can not init sentry hook"))
		}

		if hook != nil {
			log.AddHook(hook)
		} else {
			log.Warnf("can not add sentry hook, %s", sentryDSN)
		}
	}
}
