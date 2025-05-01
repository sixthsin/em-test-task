package logging

import (
	"log"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	InfoLogger  *logrus.Logger
	ErrorLogger *logrus.Logger
}

func InitLogger(infoLogFileName, errorLogFileName string) *Logger {
	logsDir := "./logs"
	if err := os.MkdirAll(logsDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	infoLogFile := filepath.Join(logsDir, infoLogFileName)
	errorLogFile := filepath.Join(logsDir, errorLogFileName)

	infoLogWriter := &lumberjack.Logger{
		Filename:   infoLogFile,
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	errorLogWriter := &lumberjack.Logger{
		Filename:   errorLogFile,
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	infoLogger := logrus.New()
	infoLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	infoLogger.SetOutput(infoLogWriter)
	infoLogger.SetLevel(logrus.InfoLevel)

	errorLogger := logrus.New()
	errorLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	errorLogger.SetOutput(errorLogWriter)
	errorLogger.SetLevel(logrus.ErrorLevel)

	return &Logger{
		InfoLogger:  infoLogger,
		ErrorLogger: errorLogger,
	}
}

func (l *Logger) Close() {
	if lumberjackWriter, ok := l.InfoLogger.Out.(*lumberjack.Logger); ok {
		lumberjackWriter.Close()
	}
	if lumberjackWriter, ok := l.ErrorLogger.Out.(*lumberjack.Logger); ok {
		lumberjackWriter.Close()
	}
}
