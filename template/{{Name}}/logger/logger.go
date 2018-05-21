package logger

import (
    "bitbucket.org/gismart/{{Name}}/config"
    "github.com/sirupsen/logrus"
    "github.com/evalphobia/logrus_sentry"
)

var Logger *logrus.Logger

func init() {
    Logger = logrus.New()
    Logger.Formatter = &logrus.JSONFormatter{}

    lvl, err := logrus.ParseLevel(config.Config.Logger.LogLevel)
    if err != nil {
        logrus.Fatalf("Failed to parse log level. %v", err)
    }
    Logger.Level = lvl

    hook, err := logrus_sentry.NewSentryHook(config.Config.Logger.SentryDSN, []logrus.Level{
        logrus.PanicLevel,
        logrus.FatalLevel,
        logrus.ErrorLevel,
    })

    if err == nil {
        Logger.Hooks.Add(hook)
    }
}
