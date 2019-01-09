package opentracing

import (
	"context"
	"github.com/Azure/azure-amqp-common-go/trace"
	"github.com/opentracing/opentracing-go"
)

func init() {
	trace.Register(new(Trace))
}

type (
	// Trace is the implementation of the OpenTracing trace abstraction
	Trace struct{}

	// Span is the implementation of the OpenTracing Span abstraction
	Span struct {
		span opentracing.Span
	}
)

// StartSpan starts and returns a Span with `operationName`, using
// any Span found within `ctx` as a ChildOfRef. If no such parent could be
// found, StartSpanFromContext creates a root (parentless) Span.
func (t *Trace) StartSpan(ctx context.Context, operationName string, opts ...interface{}) (context.Context, trace.Spanner) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName, toOTOption(opts...)...)
	return ctx, &Span{span: span}
}

// StartSpanWithRemoteParent starts and returns a Span with `operationName`, using
// reference span as FollowsFrom
func (t *Trace) StartSpanWithRemoteParent(ctx context.Context, operationName string, reference interface{}, opts ...interface{}) (context.Context, trace.Spanner) {
	if sp, ok := reference.(opentracing.SpanContext); ok {
		span := opentracing.StartSpan(operationName, append(toOTOption(opts...), opentracing.FollowsFrom(sp))...)
		ctx = opentracing.ContextWithSpan(ctx, span)
		return ctx, &Span{span: span}
	}
	return t.StartSpan(ctx, operationName)
}

// FromContext returns the `Span` previously associated with `ctx`, or
// `nil` if no such `Span` could be found.
func (t *Trace) FromContext(ctx context.Context) trace.Spanner {
	sp := opentracing.SpanFromContext(ctx)
	return &Span{span: sp}
}

// AddAttributes a tags to the span.
//
// If there is a pre-existing tag set for `key`, it is overwritten.
func (s *Span) AddAttributes(attributes ...trace.Attribute) {
	for _, attr := range attributes {
		s.span.SetTag(attr.Key, attr.Value)
	}
}

// End sets the end timestamp and finalizes Span state.
//
// With the exception of calls to Context() (which are always allowed),
// Finish() must be the last call made to any span instance, and to do
// otherwise leads to undefined behavior.
func (s *Span) End() {
	s.span.Finish()
}

// Logger returns a trace.Logger for the span
func (s *Span) Logger() trace.Logger {
	return &trace.SpanLogger{Span: s}
}

func toOTOption(opts ...interface{}) []opentracing.StartSpanOption {
	var ocStartOptions []opentracing.StartSpanOption
	for _, opt := range opts {
		if o, ok := opt.(opentracing.StartSpanOption); ok {
			ocStartOptions = append(ocStartOptions, o)
		}
	}
	return ocStartOptions
}
