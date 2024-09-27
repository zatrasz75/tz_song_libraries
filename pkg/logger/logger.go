package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type LoggersInterface interface {
	Error(message string, err error)
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Fatal(message string, err error)
	Debug(message string, args ...interface{})
}
type MyLogger struct {
	logger  *log.Logger
	console *os.File
}

func NewLogger() LoggersInterface {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	myLogger := &MyLogger{
		logger: logger,
	}

	return myLogger
}

func logWithCallerInfo(file string, line int, level string, message string, args ...interface{}) {
	caller := fmt.Sprintf("%s:%d", filepath.Base(file), line)

	var str strings.Builder
	str.WriteString("[")
	str.WriteString(level)
	str.WriteString("]")
	str.WriteString(" ")
	str.WriteString(caller)
	str.WriteString(" ")
	str.WriteString(message)

	formattedMessage := fmt.Sprintf(str.String(), args...)

	log.Println(formattedMessage)
}

// Error записывает сообщение об ошибке в лог вместе с контекстом вызова функции.
// Параметр err содержит ошибку, связанную с данным сообщением.
func (l *MyLogger) Error(message string, err error) {
	_, file, line, _ := runtime.Caller(1)
	if l.logger != nil {
		logWithCallerInfo(file, line, "ERROR", "%s: %v", message, err)
	} else {
		log.Println("No logger available.")
	}
}

// Info записывает информационное сообщение в лог вместе с контекстом вызова функции.
// Параметры args содержат дополнительные данные для сообщения.
func (l *MyLogger) Info(message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if l.logger != nil {
		logWithCallerInfo(file, line, "INFO", message, args...)
	} else {
		log.Println("No logger available.")
	}
}

// Warn записывает предупреждение в лог вместе с контекстом вызова функции.
// Параметры args содержат дополнительные данные для сообщения.
func (l *MyLogger) Warn(message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if l.logger != nil {
		logWithCallerInfo(file, line, "WARN", message, args...)
	} else {
		log.Println("No logger available.")
	}
}

// Fatal записывает фатальное сообщение в лог вместе с контекстом вызова функции
// и завершает приложение с кодом ошибки 1.
// Параметр err содержит ошибку, связанную с данным сообщением.
func (l *MyLogger) Fatal(message string, err error) {
	_, file, line, _ := runtime.Caller(1)
	if l.logger != nil {
		logWithCallerInfo(file, line, "FATAL", "%s: %v", message, err)
		os.Exit(1) // Завершаем приложение с кодом ошибки
	} else {
		log.Println("No logger available.")
	}
}

// Debug записывает информационное сообщение в лог вместе с контекстом вызова функции.
// Параметры args содержат дополнительные данные для сообщения.
func (l *MyLogger) Debug(message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if l.logger != nil {
		logWithCallerInfo(file, line, "DEBUG", message, args...)
	} else {
		log.Println("No logger available.")
	}
}
