package logger

import (
    "bitbucket.org/gismart/{{Name}}/config"
    "github.com/sirupsen/logrus"
    sentry "github.com/evalphobia/logrus_sentry"
    "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.ParseLevel(config.Config.Logger.LogLevel)

	lvl, err := log.ParseLevel(config.Config.Logger.LogLevel)
	if err != nil {
		log.Fatalf("Failed to parse log level: %v", err)
	}
	log.SetLevel(lvl)

	hook, err := sentry.NewSentryHook(config.Config.Logger.SentryDSN, []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
	})

	if err != nil {
		log.Print(errors.Wrap(err, "logger init sentry hook"))
	}

	if hook != nil {
		log.AddHook(hook)
	}
}
