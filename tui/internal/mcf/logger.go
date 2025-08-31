package mcf

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger provides structured logging for the MCF TUI
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	logFile     *os.File
}

// LogLevel represents different log levels
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	ERROR
)

// NewLogger creates a new logger instance
func NewLogger(logDir string, enableDebug bool) (*Logger, error) {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create log file with timestamp
	logFileName := fmt.Sprintf("mcf-tui-%s.log", time.Now().Format("2006-01-02"))
	logPath := filepath.Join(logDir, logFileName)

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create multi-writers to log to both file and stdout (for debug)
	var infoWriter, errorWriter, debugWriter io.Writer

	if enableDebug {
		// Log to both file and stdout when debugging
		infoWriter = io.MultiWriter(logFile, os.Stdout)
		errorWriter = io.MultiWriter(logFile, os.Stderr)
		debugWriter = io.MultiWriter(logFile, os.Stdout)
	} else {
		// Log only to file in normal mode
		infoWriter = logFile
		errorWriter = logFile
		debugWriter = logFile
	}

	logger := &Logger{
		infoLogger:  log.New(infoWriter, "[INFO] ", log.LstdFlags|log.Lshortfile),
		errorLogger: log.New(errorWriter, "[ERROR] ", log.LstdFlags|log.Lshortfile),
		debugLogger: log.New(debugWriter, "[DEBUG] ", log.LstdFlags|log.Lshortfile),
		logFile:     logFile,
	}

	logger.Info("MCF TUI Logger initialized", "logFile", logPath)
	return logger, nil
}

// Close closes the log file
func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// Info logs an info message
func (l *Logger) Info(message string, keyvals ...interface{}) {
	l.infoLogger.Printf("%s %s", message, l.formatKeyvals(keyvals...))
}

// Error logs an error message
func (l *Logger) Error(message string, err error, keyvals ...interface{}) {
	errorMsg := message
	if err != nil {
		errorMsg = fmt.Sprintf("%s: %v", message, err)
	}
	l.errorLogger.Printf("%s %s", errorMsg, l.formatKeyvals(keyvals...))
}

// Debug logs a debug message
func (l *Logger) Debug(message string, keyvals ...interface{}) {
	l.debugLogger.Printf("%s %s", message, l.formatKeyvals(keyvals...))
}

// formatKeyvals formats key-value pairs for logging
func (l *Logger) formatKeyvals(keyvals ...interface{}) string {
	if len(keyvals) == 0 {
		return ""
	}

	result := "["
	for i := 0; i < len(keyvals); i += 2 {
		if i > 0 {
			result += " "
		}

		key := fmt.Sprintf("%v", keyvals[i])

		var value string
		if i+1 < len(keyvals) {
			value = fmt.Sprintf("%v", keyvals[i+1])
		} else {
			value = "nil"
		}

		result += fmt.Sprintf("%s=%s", key, value)
	}
	result += "]"

	return result
}

// LogCommandExecution logs command execution details
func (l *Logger) LogCommandExecution(commandName string, args []string, success bool, output string, duration time.Duration) {
	l.Info("Command executed",
		"command", commandName,
		"args", fmt.Sprintf("%v", args),
		"success", success,
		"duration", duration,
		"outputLength", len(output))

	if !success {
		l.Error("Command failed", nil, "command", commandName, "output", output)
	}
}

// LogMCFOperation logs MCF adapter operations
func (l *Logger) LogMCFOperation(operation string, details map[string]interface{}) {
	keyvals := []interface{}{"operation", operation}
	for k, v := range details {
		keyvals = append(keyvals, k, v)
	}
	l.Info("MCF operation", keyvals...)
}

// LogUIEvent logs UI events and interactions
func (l *Logger) LogUIEvent(event string, view string, details map[string]interface{}) {
	keyvals := []interface{}{"event", event, "view", view}
	for k, v := range details {
		keyvals = append(keyvals, k, v)
	}
	l.Debug("UI event", keyvals...)
}

// LogPerformance logs performance metrics
func (l *Logger) LogPerformance(operation string, duration time.Duration, details map[string]interface{}) {
	keyvals := []interface{}{"operation", operation, "duration", duration}
	for k, v := range details {
		keyvals = append(keyvals, k, v)
	}
	l.Debug("Performance", keyvals...)
}

