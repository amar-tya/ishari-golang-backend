package logger

import (
    "log"
    "os"
    "path/filepath"
)

type Logger interface {
    Info(msg string, fields ...any)
    Error(msg string, fields ...any)
}

type stdLogger struct{}

func New() Logger { return &stdLogger{} }

func (l *stdLogger) Info(msg string, fields ...any) {
    args := append([]any{"INFO:", msg}, fields...)
    log.Println(args...)
}

func (l *stdLogger) Error(msg string, fields ...any) {
    args := append([]any{"ERROR:", msg}, fields...)
    log.Println(args...)
}

type fileLogger struct {
    lg *log.Logger
}

// NewFile creates a logger that writes to the given file path.
// The directory is created if it does not exist; logs are appended.
func NewFile(path string) (Logger, error) {
    if dir := filepath.Dir(path); dir != "." && dir != "" {
        if err := os.MkdirAll(dir, 0o755); err != nil {
            return nil, err
        }
    }
    f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
    if err != nil {
        return nil, err
    }
    // Note: we intentionally do not expose close here; the OS will close on process exit.
    // If needed later, we can extend the interface.
    lg := log.New(f, "", log.LstdFlags|log.LUTC)
    return &fileLogger{lg: lg}, nil
}

func (l *fileLogger) Info(msg string, fields ...any) {
    args := append([]any{"INFO:", msg}, fields...)
    l.lg.Println(args...)
}

func (l *fileLogger) Error(msg string, fields ...any) {
    args := append([]any{"ERROR:", msg}, fields...)
    l.lg.Println(args...)
}
