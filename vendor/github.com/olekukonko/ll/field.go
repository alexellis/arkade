package ll

import (
	"fmt"
	"github.com/olekukonko/ll/lx"
	"os"
	"strings"
)

// FieldBuilder enables fluent addition of fields before logging.
// It acts as a builder pattern to attach key-value pairs (fields) to log entries,
// allowing structured logging with metadata. The builder supports chaining to add fields
// and log messages at various levels (Info, Debug, etc.) in a single expression.
type FieldBuilder struct {
	logger *Logger                // Associated logger instance for logging operations
	fields map[string]interface{} // Fields to include in the log entry as key-value pairs
}

// Logger creates a new logger with the builder’s fields embedded in its context.
// It clones the parent logger and copies the builder’s fields into the new logger’s context,
// enabling persistent field inclusion in subsequent logs. This method supports fluent chaining
// after Fields or Field calls.
// Example:
//
//	logger := New("app").Enable()
//	newLogger := logger.Fields("user", "alice").Logger()
//	newLogger.Info("Action") // Output: [app] INFO: Action [user=alice]
func (fb *FieldBuilder) Logger() *Logger {
	// Clone the parent logger to preserve its configuration
	newLogger := fb.logger.Clone()
	// Initialize a new context map to avoid modifying the parent’s context
	newLogger.context = make(map[string]interface{})
	// Copy builder’s fields into the new logger’s context
	for k, v := range fb.fields {
		newLogger.context[k] = v
	}
	return newLogger
}

// Info logs a message at Info level with the builder’s fields.
// It formats the message using the provided format string and arguments, then delegates
// to the logger’s internal log method. If fields are nil, it returns early to avoid logging.
// This method is part of the fluent API, typically called after adding fields.
func (fb *FieldBuilder) Info(format string, args ...any) {
	// Skip logging if fields are nil to prevent invalid log entries
	if fb.fields == nil {
		return
	}
	// Format the message using the provided arguments
	msg := fmt.Sprintf(format, args...)
	// Log at Info level with the builder’s fields, no stack trace
	fb.logger.log(lx.LevelInfo, lx.ClassText, msg, fb.fields, false)
}

// Debug logs a message at Debug level with the builder’s fields.
// It formats the message and delegates to the logger’s log method, returning early if
// fields are nil. This method is used for debugging information that may be disabled in
// production environments.
func (fb *FieldBuilder) Debug(format string, args ...any) {
	// Skip logging if fields are nil
	if fb.fields == nil {
		return
	}
	// Format the message
	msg := fmt.Sprintf(format, args...)
	// Log at Debug level with the builder’s fields, no stack trace
	fb.logger.log(lx.LevelDebug, lx.ClassText, msg, fb.fields, false)
}

// Warn logs a message at Warn level with the builder’s fields.
// It formats the message and delegates to the logger’s log method, returning early if
// fields are nil. This method is used for warning conditions that do not halt execution.
func (fb *FieldBuilder) Warn(format string, args ...any) {
	// Skip logging if fields are nil
	if fb.fields == nil {
		return
	}
	// Format the message
	msg := fmt.Sprintf(format, args...)
	// Log at Warn level with the builder’s fields, no stack trace
	fb.logger.log(lx.LevelWarn, lx.ClassText, msg, fb.fields, false)
}

// Error logs a message at Error level with the builder’s fields.
// It formats the message and delegates to the logger’s log method, returning early if
// fields are nil. This method is used for error conditions that may require attention.
func (fb *FieldBuilder) Error(format string, args ...any) {
	// Skip logging if fields are nil
	if fb.fields == nil {
		return
	}
	// Format the message
	msg := fmt.Sprintf(format, args...)
	// Log at Error level with the builder’s fields, no stack trace
	fb.logger.log(lx.LevelError, lx.ClassText, msg, fb.fields, false)
}

// Stack logs a message at Error level with a stack trace and the builder’s fields.
// It formats the message and delegates to the logger’s log method, including a stack trace.
// Returns early if fields are nil. This method is useful for debugging critical errors.
func (fb *FieldBuilder) Stack(format string, args ...any) {
	// Skip logging if fields are nil
	if fb.fields == nil {
		return
	}
	// Format the message
	msg := fmt.Sprintf(format, args...)
	// Log at Error level with the builder’s fields and a stack trace
	fb.logger.log(lx.LevelError, lx.ClassText, msg, fb.fields, true)
}

// Fatal logs a message at Error level with a stack trace and the builder’s fields, then exits.
// It constructs the message from variadic arguments, logs it with a stack trace, and terminates
// the program with exit code 1. Returns early if fields are nil. This method is used for
// unrecoverable errors.
func (fb *FieldBuilder) Fatal(args ...any) {
	// Skip logging if fields are nil
	if fb.fields == nil {
		return
	}
	// Build the message by concatenating arguments with spaces
	var builder strings.Builder
	for i, arg := range args {
		if i > 0 {
			builder.WriteString(lx.Space)
		}
		builder.WriteString(fmt.Sprint(arg))
	}
	// Log at Error level with the builder’s fields and a stack trace
	fb.logger.log(lx.LevelError, lx.ClassText, builder.String(), fb.fields, true)
	// Exit the program with status code 1
	os.Exit(1)
}

// Panic logs a message at Error level with a stack trace and the builder’s fields, then panics.
// It constructs the message from variadic arguments, logs it with a stack trace, and triggers
// a panic with the message. Returns early if fields are nil. This method is used for critical
// errors that require immediate program termination with a panic.
func (fb *FieldBuilder) Panic(args ...any) {
	// Skip logging if fields are nil
	if fb.fields == nil {
		return
	}
	// Build the message by concatenating arguments with spaces
	var builder strings.Builder
	for i, arg := range args {
		if i > 0 {
			builder.WriteString(lx.Space)
		}
		builder.WriteString(fmt.Sprint(arg))
	}
	msg := builder.String()
	// Log at Error level with the builder’s fields and a stack trace
	fb.logger.log(lx.LevelError, lx.ClassText, msg, fb.fields, true)
	// Trigger a panic with the formatted message
	panic(msg)
}

// Err adds one or more errors to the FieldBuilder as a field and logs them.
// It stores non-nil errors in the "error" field: a single error if only one is non-nil,
// or a slice of errors if multiple are non-nil. It logs the concatenated string representations
// of non-nil errors (e.g., "failed 1; failed 2") at the Error level. Returns the FieldBuilder
// for chaining, allowing further field additions or logging. Thread-safe via the logger’s mutex.
// Example:
//
//	logger := New("app").Enable()
//	err1 := errors.New("failed 1")
//	err2 := errors.New("failed 2")
//	logger.Fields("k", "v").Err(err1, err2).Info("Error occurred")
//	// Output: [app] ERROR: failed 1; failed 2
//	//         [app] INFO: Error occurred [error=[failed 1 failed 2] k=v]
func (fb *FieldBuilder) Err(errs ...error) *FieldBuilder {
	// Initialize fields map if nil
	if fb.fields == nil {
		fb.fields = make(map[string]interface{})
	}

	// Collect non-nil errors and build log message
	var nonNilErrors []error
	var builder strings.Builder
	count := 0
	for i, err := range errs {
		if err != nil {
			if i > 0 && count > 0 {
				builder.WriteString("; ")
			}
			builder.WriteString(err.Error())
			nonNilErrors = append(nonNilErrors, err)
			count++
		}
	}

	// Set error field and log if there are non-nil errors
	if count > 0 {
		if count == 1 {
			// Store single error directly
			fb.fields["error"] = nonNilErrors[0]
		} else {
			// Store slice of errors
			fb.fields["error"] = nonNilErrors
		}
		// Log concatenated error messages at Error level
		fb.logger.log(lx.LevelError, lx.ClassText, builder.String(), nil, false)
	}

	// Return FieldBuilder for chaining
	return fb
}

// Merge adds additional key-value pairs to the FieldBuilder.
// It processes variadic arguments as key-value pairs, expecting string keys. Non-string keys
// or uneven pairs generate an "error" field with a descriptive message. Returns the FieldBuilder
// for chaining to allow further field additions or logging.
// Example:
//
//	logger := New("app").Enable()
//	logger.Fields("k1", "v1").Merge("k2", "v2").Info("Action") // Output: [app] INFO: Action [k1=v1 k2=v2]
func (fb *FieldBuilder) Merge(pairs ...any) *FieldBuilder {
	// Process pairs as key-value, advancing by 2
	for i := 0; i < len(pairs)-1; i += 2 {
		// Ensure the key is a string
		if key, ok := pairs[i].(string); ok {
			fb.fields[key] = pairs[i+1]
		} else {
			// Log an error field for non-string keys
			fb.fields["error"] = fmt.Errorf("non-string key in Merge: %v", pairs[i])
		}
	}
	// Check for uneven pairs (missing value)
	if len(pairs)%2 != 0 {
		fb.fields["error"] = fmt.Errorf("uneven key-value pairs in Merge: [%v]", pairs[len(pairs)-1])
	}
	return fb
}
