package logger

import "context"

// Logger is the application-wide logging port used by the domain and application layers.
//
// It follows a structured logging style: pass a message and optional key-value pairs
// (even number of values) to annotate the log entry.
//
// Implementations live in the infrastructure layer.
type Logger interface {
	Debug(msg string, keyVals ...any)
	Info(msg string, keyVals ...any)
	Warn(msg string, keyVals ...any)
	Error(msg string, keyVals ...any)

	// With returns a derived logger that always includes the supplied
	// structured fields.
	With(keyVals ...any) Logger

	// WithContext attaches context to the logger if supported by the implementation.
	WithContext(ctx context.Context) Logger
}
