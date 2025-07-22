package logger

import (
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func Init(level string) {
	log = logrus.New()
	log.SetOutput(os.Stdout)

	// Set formatter based on environment
	if strings.ToLower(os.Getenv("LOG_FORMAT")) == "text" {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := strings.Split(f.File, "/")
				return f.Function, filename[len(filename)-1] + ":" + string(rune(f.Line))
			},
		})
	} else {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z",
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := strings.Split(f.File, "/")
				return f.Function, filename[len(filename)-1]
			},
		})
	}

	// Enable caller information in development
	if strings.ToLower(os.Getenv("LOG_CALLER")) == "true" {
		log.SetReportCaller(true)
	}

	// Set log level
	switch strings.ToLower(level) {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn", "warning":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "panic":
		log.SetLevel(logrus.PanicLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
}

// GetLogger returns the logger instance
func GetLogger() *logrus.Logger {
	if log == nil {
		Init("info") // Default initialization
	}
	return log
}

// GetOutput returns the logger output
func GetOutput() io.Writer {
	return GetLogger().Out
}

// WithFields creates an entry from the standard logger and adds multiple fields to it
func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}

// WithField creates an entry from the standard logger and adds a field to it
func WithField(key string, value interface{}) *logrus.Entry {
	return GetLogger().WithField(key, value)
}

// WithError creates an entry from the standard logger and adds an error to it
func WithError(err error) *logrus.Entry {
	return GetLogger().WithError(err)
}

// Debug logs a message at level Debug
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Debugf logs a message at level Debug
func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

// Debugln logs a message at level Debug
func Debugln(args ...interface{}) {
	GetLogger().Debugln(args...)
}

// Info logs a message at level Info
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Infof logs a message at level Info
func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

// Infoln logs a message at level Info
func Infoln(args ...interface{}) {
	GetLogger().Infoln(args...)
}

// Warn logs a message at level Warn
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Warnf logs a message at level Warn
func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

// Warnln logs a message at level Warn
func Warnln(args ...interface{}) {
	GetLogger().Warnln(args...)
}

// Error logs a message at level Error
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// Errorf logs a message at level Error
func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

// Errorln logs a message at level Error
func Errorln(args ...interface{}) {
	GetLogger().Errorln(args...)
}

// Fatal logs a message at level Fatal then the process will exit with status set to 1
func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

// Fatalf logs a message at level Fatal then the process will exit with status set to 1
func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

// Fatalln logs a message at level Fatal then the process will exit with status set to 1
func Fatalln(args ...interface{}) {
	GetLogger().Fatalln(args...)
}

// Panic logs a message at level Panic
func Panic(args ...interface{}) {
	GetLogger().Panic(args...)
}

// Panicf logs a message at level Panic
func Panicf(format string, args ...interface{}) {
	GetLogger().Panicf(format, args...)
}

// Panicln logs a message at level Panic
func Panicln(args ...interface{}) {
	GetLogger().Panicln(args...)
}

// SetLevel sets the log level
func SetLevel(level logrus.Level) {
	GetLogger().SetLevel(level)
}

// GetLevel returns the current log level
func GetLevel() logrus.Level {
	return GetLogger().GetLevel()
}

// SetOutput sets the logger output
func SetOutput(output io.Writer) {
	GetLogger().SetOutput(output)
}

// SetFormatter sets the logger formatter
func SetFormatter(formatter logrus.Formatter) {
	GetLogger().SetFormatter(formatter)
}
