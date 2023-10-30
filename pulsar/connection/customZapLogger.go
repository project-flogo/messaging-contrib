package connection

import (
	"strings"

	pulsarLogger "github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/log"
)

// Do not print pulsar cleint's Info logs untill engine log level is Debug or lower
type zapLoggerWrapper struct {
	logger log.Logger
}

func (z *zapLoggerWrapper) SubLogger(fields pulsarLogger.Fields) pulsarLogger.Logger {
	var fieldValues []interface{}
	for k, v := range fields {
		fieldValues = append(fieldValues, k)
		fieldValues = append(fieldValues, v)
	}
	return &zapLoggerWrapper{logger: log.ChildLoggerWithFields(z.logger, fieldValues...)}
}
func (z *zapLoggerWrapper) WithFields(fields pulsarLogger.Fields) pulsarLogger.Entry {
	var fieldValues []interface{}
	for k, v := range fields {
		fieldValues = append(fieldValues, k)
		fieldValues = append(fieldValues, v)
	}
	return zapEntry{logger: log.ChildLoggerWithFields(z.logger, fieldValues...)}
}
func (z *zapLoggerWrapper) WithField(name string, value interface{}) pulsarLogger.Entry {
	return zapEntry{logger: log.ChildLoggerWithFields(z.logger, name, value)}
}
func (z *zapLoggerWrapper) WithError(err error) pulsarLogger.Entry {
	return zapEntry{logger: log.ChildLoggerWithFields(z.logger, "error", err.Error())}
}
func (z *zapLoggerWrapper) Debug(args ...interface{}) {
	z.logger.Debug(args)
}
func (z *zapLoggerWrapper) Info(args ...interface{}) {
	if log.ToLogLevel(engineLogLevel) <= log.DebugLevel || isPrintable(args...) {
		z.logger.Info(args)
	}
}

func isPrintable(args ...interface{}) bool {
	for _, v := range args {
		vs, _ := coerce.ToString(v)
		if strings.Contains(vs, "Reconnected consumer to broker") || strings.Contains(vs, "Reconnected producer to broker") {
			return true
		}
	}
	return false
}

func (z *zapLoggerWrapper) Warn(args ...interface{}) {
	// Intermittent connection logs getting logged in with logging level 'WARN". Chaning them to log level "ERROR".
	if isConnectionLog(args...) {
		z.logger.Error(args)
		return
	}
	z.logger.Warn(args)

}
func (z *zapLoggerWrapper) Error(args ...interface{}) {
	z.logger.Error(args)
}
func (z *zapLoggerWrapper) Debugf(format string, args ...interface{}) {
	z.logger.Debugf(format, args)
}
func (z *zapLoggerWrapper) Infof(format string, args ...interface{}) {
	if log.ToLogLevel(engineLogLevel) <= log.DebugLevel {
		z.logger.Infof(format, args)
	}
}
func (z *zapLoggerWrapper) Warnf(format string, args ...interface{}) {
	// Intermittent connection logs getting logged in with logging level 'WARN". Chaning them to log level "ERROR".
	if isConnectionLog(args...) {
		z.logger.Errorf(format, args)
		return
	}
	z.logger.Warnf(format, args)

}
func (z *zapLoggerWrapper) Errorf(format string, args ...interface{}) {
	z.logger.Errorf(format, args)
}

type zapEntry struct {
	logger log.Logger
}

func (z zapEntry) WithFields(fields pulsarLogger.Fields) pulsarLogger.Entry {
	var fieldValues []interface{}
	for k, v := range fields {
		fieldValues = append(fieldValues, k)
		fieldValues = append(fieldValues, v)
	}
	return zapEntry{logger: log.ChildLoggerWithFields(z.logger, fieldValues...)}
}
func (z zapEntry) WithField(name string, value interface{}) pulsarLogger.Entry {
	return zapEntry{log.ChildLoggerWithFields(z.logger, name, value)}
}
func (z zapEntry) Debug(args ...interface{}) {
	z.logger.Debug(args)
}
func (z zapEntry) Info(args ...interface{}) {
	if log.ToLogLevel(engineLogLevel) <= log.DebugLevel || isPrintable(args...) {
		z.logger.Info(args)
	}
}
func (z zapEntry) Warn(args ...interface{}) {
	// Intermittent connection logs getting logged in with logging level 'WARN". Chaning them to log level "ERROR".
	if isConnectionLog(args...) {
		z.logger.Error(args)
		return
	}
	z.logger.Warn(args)
}
func (z zapEntry) Error(args ...interface{}) {
	z.logger.Error(args)
}
func (z zapEntry) Debugf(format string, args ...interface{}) {
	z.logger.Debugf(format, args)
}
func (z zapEntry) Infof(format string, args ...interface{}) {
	if log.ToLogLevel(engineLogLevel) <= log.DebugLevel {
		z.logger.Infof(format, args)
	}
}
func (z zapEntry) Warnf(format string, args ...interface{}) {
	// Intermittent connection logs getting logged in with logging level 'WARN". Chaning them to log level "ERROR".
	if isConnectionLog(args...) {
		z.logger.Errorf(format, args)
		return
	}
	z.logger.Warnf(format, args)
}
func (z zapEntry) Errorf(format string, args ...interface{}) {
	z.logger.Errorf(format, args)
}

// check if log is a connection log.
func isConnectionLog(args ...interface{}) bool {
	for _, v := range args {
		vs, _ := coerce.ToString(v)
		if strings.Contains(strings.ToLower(vs), "connection") || strings.Contains(strings.ToLower(vs), "failed") {
			return true
		}
	}
	return false
}
