package log

import (
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	_ log.Logger = LogrusLogger{}
	_ log.Logger = entry{}
)

type LogrusLogger struct {
	logger *logrus.Logger
}

func NewDefaultLogger() LogrusLogger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	return LogrusLogger{
		logger: logger,
	}
}

func NewLogger(cfg Config) LogrusLogger {
	logger := logrus.New()

	switch cfg.Format {
	case FormatText:
		logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	case FormatJSON:
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	switch strings.ToLower(cfg.Level) {
	case DebugLevel:
		logger.SetLevel(logrus.DebugLevel)
	case InfoLevel:
		logger.SetLevel(logrus.InfoLevel)
	case WarnLevel:
		logger.SetLevel(logrus.WarnLevel)
	case ErrorLevel:
		logger.SetLevel(logrus.ErrorLevel)
	}

	logger.SetOutput(os.Stdout)
	return LogrusLogger{
		logger: logger,
	}
}

func (l *LogrusLogger) SetOutput(output io.Writer) {
	l.logger.SetOutput(os.Stdout)
}

func (l LogrusLogger) Debug(msg string, keyvals ...interface{}) {
	l.logger.WithFields(argsToFields(keyvals...)).Debug(msg)
}

func (l LogrusLogger) Info(msg string, keyvals ...interface{}) {
	l.logger.WithFields(argsToFields(keyvals...)).Info(msg)
}

func (l LogrusLogger) Error(msg string, keyvals ...interface{}) {
	l.logger.WithFields(argsToFields(keyvals...)).Error(msg)
}

func (l LogrusLogger) With(keyvals ...interface{}) log.Logger {
	return entry{
		l.logger.WithFields(argsToFields(keyvals...)),
	}
}

type entry struct {
	*logrus.Entry
}

func (e entry) Debug(msg string, keyvals ...interface{}) {
	e.Entry.WithFields(argsToFields(keyvals...)).Debug(msg)
}

func (e entry) Info(msg string, keyvals ...interface{}) {
	e.Entry.WithFields(argsToFields(keyvals...)).Info(msg)
}

func (e entry) Error(msg string, keyvals ...interface{}) {
	e.Entry.WithFields(argsToFields(keyvals...)).Error(msg)
}

func (e entry) With(keyvals ...interface{}) log.Logger {
	return entry{
		e.WithFields(argsToFields(keyvals...)),
	}
}

func argsToFields(keyvals ...interface{}) logrus.Fields {
	var fields = make(logrus.Fields)
	if len(keyvals)%2 != 0 {
		return fields
	}

	for i := 0; i < len(keyvals); i += 2 {
		fields[keyvals[i].(string)] = keyvals[i+1]
	}
	return fields
}
