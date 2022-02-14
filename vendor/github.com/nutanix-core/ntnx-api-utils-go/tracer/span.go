//
// Copyright (c) 2022 Nutanix Inc. All rights reserved.
//
// Author: akash.shinde@nutanix.com
//
// Implements Span interface. Span interface defined here is a wrapper over
// opentelemetry trace.Span
//

package tracer

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	oteltracing "go.opentelemetry.io/otel/trace"
)

// Span struct, stores opentracing span.
type Span struct {
	s *oteltracing.Span
}

// Span interface implements the functions that can be used from services for
// operations like close span, set tag in span etc.

type SpanIfc interface {
	Finish()
	SetTag(key string, value interface{})
	getSpan() oteltracing.Span
}

// Checks that Span struct implements SpanIfc
var _ SpanIfc = (*Span)(nil)

// Finishes span.
func (span *Span) Finish() {
	if span == nil || span.s == nil {
		return
	}
	span.getSpan().End()
}

// Adds tag to span.
// NOTE: Tag values can be numeric types, strings, or bools. At the of the
//       writing of this function, behavior of other tag value types is
//       undefined at the OpenTracing level. If caller is using value that is
//       not of any of the above type, it is advised to test the tag before
//       merging.
func (span *Span) SetTag(key string, value interface{}) {
	if span == nil || span.s == nil {
		return
	}

	if str_val, ok := value.(string); ok {
		span.getSpan().SetAttributes(attribute.String(key, str_val))
	} else if int_val, ok := value.(int); ok {
		span.getSpan().SetAttributes(attribute.Int(key, int_val))
	} else if int64_val, ok := value.(int64); ok {
		span.getSpan().SetAttributes(attribute.Int64(key, int64_val))
	} else if float64_val, ok := value.(float64); ok {
		span.getSpan().SetAttributes(attribute.Float64(key, float64_val))
	} else if bool_val, ok := value.(bool); ok {
		span.getSpan().SetAttributes(attribute.Bool(key, bool_val))
	} else {
		span.getSpan().SetAttributes(attribute.String(key, fmt.Sprintf("%v", value)))
	}
}

func (span *Span) getSpan() oteltracing.Span {
	if span == nil {
		return nil
	}
	return *span.s
}
