package trace

import (
	"context"
)

var tracer Tracer

// Register a Tracer instance
func Register(t Tracer) {
	tracer = t
}

// BoolAttribute returns a bool-valued attribute.
func BoolAttribute(key string, value bool) Attribute {
	return Attribute{Key: key, Value: value}
}

// StringAttribute returns a string-valued attribute.
func StringAttribute(key, value string) Attribute {
	return Attribute{Key: key, Value: value}
}

// Int64Attribute returns an int64-valued attribute.
func Int64Attribute(key string, value int64) Attribute {
	return Attribute{Key: key, Value: value}
}

// StartSpan starts a new child span
func StartSpan(ctx context.Context, operationName string, opts ...interface{}) (context.Context, Spanner) {
	if tracer == nil {
		return ctx, new(nopSpanner)
	}
	return tracer.StartSpan(ctx, operationName, opts)
}

// StartSpanWithRemoteParent starts a new child span of the span from the given parent.
func StartSpanWithRemoteParent(ctx context.Context, operationName string, reference interface{}, opts ...interface{}) (context.Context, Spanner) {
	if tracer == nil {
		return ctx, new(nopSpanner)
	}
	return tracer.StartSpanWithRemoteParent(ctx, operationName, reference, opts)
}

// FromContext returns the Span stored in a context, or nil if there isn't one.
func FromContext(ctx context.Context) Spanner {
	if tracer == nil {
		return new(nopSpanner)
	}
	return tracer.FromContext(ctx)
}

type (
	// Attribute is a key value pair for decorating spans
	Attribute struct {
		Key   string
		Value interface{}
	}

	// Spanner is an abstraction over OpenTracing and OpenCensus Spans
	Spanner interface {
		AddAttributes(attributes ...Attribute)
		End()
		Logger() Logger
	}

	// Tracer is an abstraction over OpenTracing and OpenCensus trace implementations
	Tracer interface {
		StartSpan(ctx context.Context, operationName string, opts ...interface{}) (context.Context, Spanner)
		StartSpanWithRemoteParent(ctx context.Context, operationName string, reference interface{}, opts ...interface{}) (context.Context, Spanner)
		FromContext(ctx context.Context) Spanner
	}

	// Logger is a generic interface for logging
	Logger interface {
		Info(msg string, attributes ...Attribute)
		Error(err error, attributes ...Attribute)
		Fatal(msg string, attributes ...Attribute)
		Debug(msg string, attributes ...Attribute)
	}

	// SpanLogger is a Logger implementation which logs to a tracing span
	SpanLogger struct {
		Span Spanner
	}

	nopLogger struct{}

	nopSpanner struct{}
)

// AddAttributes is a nop
func (ns *nopSpanner) AddAttributes(attributes ...Attribute) {}

// End is a nop
func (ns *nopSpanner) End() {}

// Logger returns a nopLogger
func (ns *nopSpanner) Logger() Logger {
	return nopLogger{}
}

// For will return a logger for a given context
func For(ctx context.Context) Logger {
	if span := tracer.FromContext(ctx); span != nil {
		return span.Logger()
	}
	return new(nopLogger)
}

// Info logs an info tag with message to a span
func (sl SpanLogger) Info(msg string, attributes ...Attribute) {
	sl.logToSpan("info", msg, attributes...)
}

// Error logs an error tag with message to a span
func (sl SpanLogger) Error(err error, attributes ...Attribute) {
	attributes = append(attributes, BoolAttribute("error", true))
	sl.logToSpan("error", err.Error(), attributes...)
}

// Fatal logs an error tag with message to a span
func (sl SpanLogger) Fatal(msg string, attributes ...Attribute) {
	attributes = append(attributes, BoolAttribute("error", true))
	sl.logToSpan("fatal", msg, attributes...)
}

// Debug logs a debug tag with message to a span
func (sl SpanLogger) Debug(msg string, attributes ...Attribute) {
	sl.logToSpan("debug", msg, attributes...)
}

func (sl SpanLogger) logToSpan(level string, msg string, attributes ...Attribute) {
	attrs := append(attributes, StringAttribute("event", msg), StringAttribute("level", level))
	sl.Span.AddAttributes(attrs...)
}

// Info nops log entry
func (sl nopLogger) Info(msg string, attributes ...Attribute) {}

// Error nops log entry
func (sl nopLogger) Error(err error, attributes ...Attribute) {}

// Fatal nops log entry
func (sl nopLogger) Fatal(msg string, attributes ...Attribute) {}

// Debug nops log entry
func (sl nopLogger) Debug(msg string, attributes ...Attribute) {}
