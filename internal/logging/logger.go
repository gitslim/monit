// Модуль logging предназначен для управления логгированием.
package logging

import "go.uber.org/zap"

// Logger представляет собой логгер, который использует библиотеку zap.
type Logger struct {
	sugar *zap.SugaredLogger
}

// NewLogger создает новый логгер.
func NewLogger() (*Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	return &Logger{sugar: sugar}, nil
}

// Close закрывает логгер.
func (l *Logger) Close() {
	l.sugar.Sync()
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.sugar.Debugw(msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.sugar.Infow(msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.sugar.Warnw(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.sugar.Errorw(msg, args...)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.sugar.Fatalw(msg, args...)
}

func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.sugar.Debugf(msg, args...)
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	l.sugar.Infof(msg, args...)
}

func (l *Logger) Warnf(msg string, args ...interface{}) {
	l.sugar.Warnf(msg, args...)
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.sugar.Errorf(msg, args...)
}

func (l *Logger) Fatalf(msg string, args ...interface{}) {
	l.sugar.Fatalf(msg, args...)
}
