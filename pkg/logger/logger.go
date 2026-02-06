package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
)

type Level uint32

const (
	DebugLevel Level = Level(logrus.DebugLevel)
	InfoLevel  Level = Level(logrus.InfoLevel)
	WarnLevel  Level = Level(logrus.WarnLevel)
	ErrorLevel Level = Level(logrus.ErrorLevel)
	FatalLevel Level = Level(logrus.FatalLevel)
	PanicLevel Level = Level(logrus.PanicLevel)
)

type Logger struct {
	logger   *logrus.Logger
	mu       sync.Mutex
	handlers []Handler
}

type Handler struct {
	Writer    io.Writer
	Formatter logrus.Formatter
}

var (
	defaultLogger *Logger
	once          sync.Once
)

func init() {
	once.Do(func() {
		defaultLogger = New()
	})
}

func New() *Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: false,
		TimestampFormat:  "2006-01-02 15:04:05",
		FullTimestamp:    true,
		DisableColors:    false,
	})

	return &Logger{
		logger:   logger,
		handlers: make([]Handler, 0),
	}
}

func Get() *Logger {
	return defaultLogger
}

func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}

func (l *Logger) SetLevel(level Level) {
	l.logger.SetLevel(logrus.Level(level))
}

func (l *Logger) SetOutput(writer io.Writer) {
	l.logger.SetOutput(writer)
}

func (l *Logger) SetFormatter(formatter logrus.Formatter) {
	l.logger.SetFormatter(formatter)
}

func (l *Logger) AddHandler(writer io.Writer, formatter logrus.Formatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.handlers = append(l.handlers, Handler{Writer: writer, Formatter: formatter})
}

func (l *Logger) log(level logrus.Level, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)

	l.mu.Lock()
	defer l.mu.Unlock()

	l.logger.Log(level, msg)

	for _, handler := range l.handlers {
		handler.Writer.Write([]byte(msg + "\n"))
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(logrus.DebugLevel, format, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(logrus.DebugLevel, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(logrus.InfoLevel, format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(logrus.InfoLevel, format, args...)
}

func (l *Logger) Success(format string, args ...interface{}) {
	l.log(logrus.InfoLevel, "[SUCCESS] "+format, args...)
}

func (l *Logger) Successf(format string, args ...interface{}) {
	l.log(logrus.InfoLevel, "[SUCCESS] "+format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(logrus.WarnLevel, format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(logrus.WarnLevel, format, args...)
}

func (l *Logger) Warning(format string, args ...interface{}) {
	l.log(logrus.WarnLevel, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(logrus.ErrorLevel, format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(logrus.ErrorLevel, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(logrus.FatalLevel, format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(logrus.FatalLevel, format, args...)
}

func (l *Logger) Panic(format string, args ...interface{}) {
	l.log(logrus.PanicLevel, format, args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.log(logrus.PanicLevel, format, args...)
}

func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.logger.WithField(key, value)
}

func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.logger.WithFields(fields)
}

func (l *Logger) WithError(err error) *logrus.Entry {
	return l.logger.WithError(err)
}

func (l *Logger) SetLogFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %w", err)
	}

	l.logger.SetOutput(io.MultiWriter(os.Stdout, file))
	return nil
}

func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func Success(format string, args ...interface{}) {
	defaultLogger.Success(format, args...)
}

func Successf(format string, args ...interface{}) {
	defaultLogger.Successf(format, args...)
}

func Warn(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	defaultLogger.Fatal(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

func SetFormatter(formatter logrus.Formatter) {
	defaultLogger.SetFormatter(formatter)
}

func SetLogFile(path string) error {
	return defaultLogger.SetLogFile(path)
}
