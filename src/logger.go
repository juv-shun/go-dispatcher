package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger(level string) (*Logger, error) {
	var err error

	logger := &Logger{logrus.New()}

	logger.Out = os.Stdout

	logger.Level, err = logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	logger.Formatter = &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	}
	return logger, nil
}

func (logger *Logger) WithWorkerInfo(w *worker) *logrus.Entry {
	return logger.WithFields(logrus.Fields{
		"dispatcher": w.dispatcher.name,
		"workerID":   w.id,
	})
}
