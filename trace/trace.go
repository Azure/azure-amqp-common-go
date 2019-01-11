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
		return ctx, new(noOpSpanner)
	}
	return tracer.StartSpan(ctx, operationName, opts)
}

// StartSpanWithRemoteParent starts a new child span of the span from the given parent.
func StartSpanWithRemoteParent(ctx context.Context, operationName string, carrier Carrier, opts ...interface{}) (context.Context, Spanner) {
	if tracer == nil {
		return ctx, new(noOpSpanner)
	}
	return tracer.StartSpanWithRemoteParent(ctx, operationName, carrier, opts)
}

// FromContext returns the Span stored in a context, or nil if there isn't one.
func FromContext(ctx context.Context) Spanner {
	if tracer == nil {
		return new(noOpSpanner)
	}
	return tracer.FromContext(ctx)
}

type (
	// Attribute is a key value pair for decorating spans
	Attribute struct {
		Key   string
		Value interface{}
	}

	// Carrier is an abstraction over OpenTracing and OpenCensus propagation carrier
	Carrier interface {
		Set(key string, value interface{})
		GetKeyValues() map[string]interface{}
	}

	// Spanner is an abstraction over OpenTracing and OpenCensus Spans
	Spanner interface {
		AddAttributes(attributes ...Attribute)
		End()
		Logger() Logger
		Inject(carrier Carrier) error
		InternalSpan() interface{}
	}

	// Tracer is an abstraction over OpenTracing and OpenCensus trace implementations
	Tracer interface {
		StartSpan(ctx context.Context, operationName string, opts ...interface{}) (context.Context, Spanner)
		StartSpanWithRemoteParent(ctx context.Context, operationName string, carrier Carrier, opts ...interface{}) (context.Context, Spanner)
		FromContext(ctx context.Context) Spanner
		NewContext(parent context.Context, span Spanner) context.Context
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

	noOpLogger struct{}

	noOpSpanner struct{}
)

// AddAttributes is a nop
func (ns *noOpSpanner) AddAttributes(attributes ...Attribute) {}

// End is a nop
func (ns *noOpSpanner) End() {}

// Logger returns a nopLogger
func (ns *noOpSpanner) Logger() Logger {
	return noOpLogger{}
}

// Inject is a nop
func (ns *noOpSpanner) Inject(carrier Carrier) error {
	return nil
}

// InternalSpan returns nil
func (ns *noOpSpanner) InternalSpan() interface{} {
	return nil
}

// For will return a logger for a given context
func For(ctx context.Context) Logger {
	if span := tracer.FromContext(ctx); span != nil {
		return span.Logger()
	}
	return new(noOpLogger)
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
func (sl noOpLogger) Info(msg string, attributes ...Attribute) {}

// Error nops log entry
func (sl noOpLogger) Error(err error, attributes ...Attribute) {}

// Fatal nops log entry
func (sl noOpLogger) Fatal(msg string, attributes ...Attribute) {}

// Debug nops log entry
func (sl noOpLogger) Debug(msg string, attributes ...Attribute) {}
