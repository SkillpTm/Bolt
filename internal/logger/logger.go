// Package logger handles and carries all relevant functions for the logging of error messages
package logger

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Logger carries the functions to write to the log file
type Logger struct {
	ErrorLogPath   string
	HistoryLogPath string
}

// History writes an event to the history.log file
func (l *Logger) History(event string, format string, args ...any) {
	file, err := os.OpenFile(l.HistoryLogPath, os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		defer file.Close()
		writeToFile(file, event, fmt.Sprintf(format, args...))
	}
}

// Error writes an error to the error.log file
func (l *Logger) Error(format string, args ...any) {
	file, err := os.OpenFile(l.ErrorLogPath, os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		defer file.Close()
		writeToFile(file, "ERROR", fmt.Sprintf(format, args...))
	}
}

// Fatal writes an error to the error.log file and then stops the program with os.Exit(0)
func (l *Logger) Fatal(format string, args ...any) {
	file, err := os.OpenFile(l.ErrorLogPath, os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		defer file.Close()
		writeToFile(file, "FATAL", fmt.Sprintf(format, args...))
	}

	os.Exit(0)
}

// Panic attempts to load the error.log itself and will definitly os.Exit(100)
func (l *Logger) Panic(format string, args ...any) {
	homeDir, err := os.UserHomeDir()
	if err == nil {
		file, err := os.OpenFile(fmt.Sprintf("%s/.local/share/bolt/error.log", homeDir), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			defer file.Close()
			writeToFile(file, "PANIC", fmt.Sprintf(format, args...))
		}
	}

	os.Exit(100)
}

// writeToFile writes the time, level and error message to the provided file
func writeToFile(file *os.File, level string, message string) {
	splitMessage := strings.Split(message, "\n")
	newMessage := ""
	rootErr := ""
	for index, part := range splitMessage {
		if index+1 == len(splitMessage) {
			rootErr = part
			break
		}

		newMessage += fmt.Sprintf("%s\"%s\"\n", strings.Repeat(" ", index), part)
	}

	file.WriteString(fmt.Sprintf("%s %s --> \"%s\"\n%s---\n", time.Now().Format("[2006-01-02 15:04:05]"), level, rootErr, newMessage))
}
