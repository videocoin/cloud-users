package logger

import (
	"os"
	"time"

	logrussentry "github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
)

func NewLogrusLogger(serviceName string, serviceVersion string) *logrus.Entry {
	l := logrus.New()
	logger := logrus.NewEntry(l)

	loglevel = os.Getenv("LOGLEVEL")
	if loglevel == "" {
		loglevel = logrus.InfoLevel.String()
	}
	level, err := logrus.ParseLevel(loglevel)
	if err != nil {
		level = logrus.InfoLevel
	}

	l.SetLevel(level)

	if level == logrus.DebugLevel {
		l.SetFormatter(&logrus.TextFormatter{TimestampFormat: time.RFC3339Nano})
	} else {
		l.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	}

	sentryDSN = os.Getenv("SENTRY_DSN")
	if sentryDSN != "" {
		sentryLevels := []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		}
		sentryTags := map[string]string{
			"service": serviceName,
			"version": serviceVersion,
		}
		sentryHook, err := logrussentry.NewAsyncWithTagsSentryHook(
			sentryDSN,
			sentryTags,
			sentryLevels,
		)
		sentryHook.StacktraceConfiguration.Enable = true
		sentryHook.Timeout = 5 * time.Second
		sentryHook.SetRelease(serviceVersion)

		if err != nil {
			logger.Warning(err)
		} else {
			l.AddHook(sentryHook)
		}
	}

	logger.WithFields(logrus.Fields{
		"service": serviceName,
		"version": serviceVersion,
	})

	return logger
}
