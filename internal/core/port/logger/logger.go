package logger

// Logger defines the interface for logging operations in the domain layer.
// This interface is framework-agnostic and can be implemented by any logging library.
type Logger interface {
	// Info logs an informational message with optional key-value fields
	Info(msg string, fields ...any)

	// Error logs an error message with optional key-value fields
	Error(msg string, fields ...any)
}
