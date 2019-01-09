package opencensus

import (
	"context"
	"github.com/Azure/azure-amqp-common-go/trace"
	oct "go.opencensus.io/trace"
)

func init() {
	trace.Register(new(Trace))
}

type (
	// Trace is the implementation of the OpenCensus trace abstraction
	Trace struct{}

	// Span is the implementation of the OpenCensus Span abstraction
	Span struct {
		span *oct.Span
	}
)

// StartSpan starts a new child span of the current span in the context. If
// there is no span in the context, creates a new trace and span.
//
// Returned context contains the newly created span. You can use it to
// propagate the returned span in process.
func (t *Trace) StartSpan(ctx context.Context, operationName string, opts ...interface{}) (context.Context, trace.Spanner) {
	ctx, span := oct.StartSpan(ctx, operationName, toOCOption(opts...)...)
	return ctx, &Span{span: span}
}

// StartSpanWithRemoteParent starts a new child span of the span from the given parent.
//
// If the incoming context contains a parent, it ignores. StartSpanWithRemoteParent is
// preferred for cases where the parent is propagated via an incoming request.
//
// Returned context contains the newly created span. You can use it to
// propagate the returned span in process.
func (t *Trace) StartSpanWithRemoteParent(ctx context.Context, operationName string, reference interface{}, opts ...interface{}) (context.Context, trace.Spanner) {
	if sp, ok := reference.(oct.SpanContext); ok {
		ctx, span := oct.StartSpanWithRemoteParent(ctx, operationName, sp, toOCOption(opts...)...)
		return ctx, &Span{span: span}
	}
	return t.StartSpan(ctx, operationName)
}

// FromContext returns the Span stored in a context, or nil if there isn't one.
func (t *Trace) FromContext(ctx context.Context) trace.Spanner {
	sp := oct.FromContext(ctx)
	return &Span{span: sp}
}

// AddAttributes sets attributes in the span.
//
// Existing attributes whose keys appear in the attributes parameter are overwritten.
func (s *Span) AddAttributes(attributes ...trace.Attribute) {
	s.span.AddAttributes(attributesToOCAttributes(attributes...)...)
}

// End ends the span.
func (s *Span) End() {
	s.span.End()
}

// Logger returns a trace.Logger for the span
func (s *Span) Logger() trace.Logger {
	return &trace.SpanLogger{Span: s}
}

func toOCOption(opts ...interface{}) []oct.StartOption {
	var ocStartOptions []oct.StartOption
	for _, opt := range opts {
		if o, ok := opt.(oct.StartOption); ok {
			ocStartOptions = append(ocStartOptions, o)
		}
	}
	return ocStartOptions
}

func attributesToOCAttributes(attributes ...trace.Attribute) []oct.Attribute {
	var ocAttributes []oct.Attribute
	for _, attr := range attributes {
		switch attr.Value.(type) {
		case int64:
			ocAttributes = append(ocAttributes, oct.Int64Attribute(attr.Key, attr.Value.(int64)))
		case string:
			ocAttributes = append(ocAttributes, oct.StringAttribute(attr.Key, attr.Value.(string)))
		case bool:
			ocAttributes = append(ocAttributes, oct.BoolAttribute(attr.Key, attr.Value.(bool)))
		}
	}
	return ocAttributes
}
