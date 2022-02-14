//
// Copyright (c) 2021 Nutanix Inc. All rights reserved.
//
// Author: deepanshu.singhal@nutanix.com
//

package tracer

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/golang/glog"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

var tracer trace.Tracer

// Returns tracing options to be used with GRPC server options. This will also
// generate span when server recieves a request.
//
// Usage:
//  unaryTrace, streamTrace := GrpcServerTraceOptions()
//  var opts []grpc.ServerOption
//  if unaryTrace != nil {
//    unaryOpt := grpc.UnaryInterceptor(unaryTrace)
//    opts = append(opts, unaryOpt)
//  }
//  if streamTrace != nil {
//    streamOpt := grpc.StreamInterceptor(streamTrace)
//    opts = append(opts, streamOpt)
//  }
//  server := grpc.NewServer(opts...)

func GrpcServerTraceOptions() (grpc.UnaryServerInterceptor,
	grpc.StreamServerInterceptor) {
	if !tracingEnabled {
		return nil, nil
	}

	unaryTracer := otelgrpc.UnaryServerInterceptor()
	streamTracer := otelgrpc.StreamServerInterceptor()
	return unaryTracer, streamTracer
}

// Appends tracing options to GRPC server options. This will also generate span
// when server recieves a request.
//
// Usage:
//  opts := GrpcServerOptionsWithTrace()
//  server := grpc.NewServer(opts...)

func GrpcServerOptionsWithTrace(opts ...grpc.ServerOption) []grpc.ServerOption {
	if !tracingEnabled {
		return opts
	}
	unaryTracer := grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor())
	streamTracer := grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor())
	opts = append(opts, unaryTracer, streamTracer)
	return opts
}

// Returns tracing options to be used with GRPC request options. This will also
// generate span from caller/client side when a request is made to another
// service.
//
// Usage:
//  unaryTrace, streamTrace := GrpcRequestTraceOptions()
//  var opts []grpc.DialOption
//  if unaryTrace != nil {
//    unaryOpt := grpc.WithUnaryInterceptor(unaryTrace)
//    opts = append(opts, unaryOpt)
//  }
//  if streamTrace != nil {
//    streamOpt := grpc.WithStreamInterceptor(streamTrace)
//    opts = append(opts, streamTrace)
//  }
//  conn, err := grpc.Dial("1.2.3.4", opts...)

func GrpcRequestTraceOptions() (grpc.UnaryClientInterceptor,
	grpc.StreamClientInterceptor) {
	if !tracingEnabled {
		return nil, nil
	}

	unaryTracer := otelgrpc.UnaryClientInterceptor()
	streamTracer := otelgrpc.StreamClientInterceptor()
	return unaryTracer, streamTracer
}

// Appends tracing options to GRPC request options. This will also generate span
// from caller/client side when a request is made to another service.
// Usage example:
//  opts := GrpcRequestOptionsWithTrace(grpc.WithInsecure())
//  conn, err := grpc.Dial("1.2.3.4", opts...)
func GrpcRequestOptionsWithTrace(opts ...grpc.DialOption) []grpc.DialOption {
	if !tracingEnabled {
		return opts
	}
	unaryTracer := grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor())
	streamTracer := grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor())

	opts = append(opts, unaryTracer, streamTracer)
	return opts
}

// Starts a child span from context.
//
// Usage expample:
//  func A() {
//    ...
//    ctx := context.Background()
//    span, ctx := StartSpan(ctx, "test-span")
//    defer span.Finish()
//    B(ctx)
//    ...
//  }
//
// Args:
//  c: Context.
//  spanName: Name of the span.
//
// Returns: Span interface and context. If nil is passed in context, empty span
//          and nil is returned.
func StartSpan(c context.Context, spanName string) (SpanIfc, context.Context) {
	otelTracer := otel.Tracer(spanName)
	spanContext, otelSpan := otelTracer.Start(c, spanName)

	span := &Span{
		s: &otelSpan,
	}
	// otelSpan.End()
	return span, spanContext
}

// Starts span of follows from relation. Can be used to create a span that is
// not a child but related to parent span. E.g., VM create task is follwed from
// VM create RPC.
//
// Usage example:
//  func A(ctx context.Context) {
//    ...
//    parentSpan := GetSpanFromContext(ctx)
//    span, ctx := StartFollowsFromSpan(ctx, parentSpan, "follows-from-span")
//    defer span.Finish()
//    B(ctx)
//    ...
//  }
//
// Args:
//  c: Context.
//  parentSpan: Span from which the newly created span will be followed.
//  spanName: Name of span to be created
//
// Returns: Span interface and context. If nil is passed in context or
//          parentSpan, empty span and nil is returned.
func StartFollowsFromSpan(c context.Context, parentSpan SpanIfc, spanName string) (SpanIfc, context.Context) {
	if c == nil || parentSpan == nil || parentSpan.getSpan() == nil {
		return &Span{}, nil
	}

	otelTracer := otel.Tracer(spanName)
	linkFromSpan := trace.Link{SpanContext: parentSpan.getSpan().SpanContext()}
	links := []trace.Link{linkFromSpan}
	spanContext, otelSpan := otelTracer.Start(c, spanName, trace.WithLinks(links...))

	otelSpan.End()
	span := &Span{
		s: &otelSpan,
	}
	return span, spanContext
}

// Returns span from context.
func GetSpanFromContext(c context.Context) SpanIfc {
	if c == nil {
		return &Span{}
	}

	otSpan := trace.SpanFromContext(c)
	span := &Span{
		s: &otSpan,
	}
	return span
}

// Injects context into HTTP request headers.
//
// Usage example:
//  func makeHTTPCall(ctx context.Context) {
//    span, ctx := StartSpan(ctx, "call-server")
//    defer span.Finish()
//    request, _ := http.NewRequest(...)
//    InjectCtxIntoHTTPHeader(span, request)
//    ...
//  }
//
// Args:
//  span: Context from span will be insterted to HTTP request headers. This
//        context needs to be extracted from server to which HTTP call is made.
//  request: HTTP request.
func InjectCtxIntoHTTPHeader(span SpanIfc, request *http.Request) {
	if span == nil || span.getSpan() == nil {
		return
	}

	ctx := context.Background()
	sctx := trace.ContextWithSpanContext(ctx, span.getSpan().SpanContext())
	// This sets value of span.kind tag as client.
	//TODO: make below commented line workable
	// otel.GetTextMapPropagator().Inject(span.getSpan().SpanContext(), propagation.HeaderCarrier(request.Header))
	otel.GetTextMapPropagator().Inject(sctx, propagation.HeaderCarrier(request.Header))
}

// Converts context into bytes. Can be used to store or transfer context.
//
// Usage example:
//  func A(ctx context.Context) {
//    span, ctx := StartSpan(ctx, "test-span")
//    defer span.Finish()
//    ctxToBytes := InjectCtxIntoBytes(span)
//    ...
//  }
//
// Args:
//  span: Context from span will be converted to bytes.
// Returns: bytes containing span details.
func InjectCtxIntoBytes(span SpanIfc) []byte {
	if span == nil || span.getSpan() == nil {
		return nil
	}

	contextMap := make(map[string]string)
	ctx := context.Background()
	newContext := trace.ContextWithSpanContext(ctx, span.getSpan().SpanContext())
	// This sets value of span.kind tag as client.
	//TODO: make below commented line workable
	// otel.GetTextMapPropagator().Inject(span.getSpan().SpanContext(), propagation.MapCarrier(ctxMap))
	otel.GetTextMapPropagator().Inject(newContext, propagation.MapCarrier(contextMap))

	ctxBytes, err := json.Marshal(contextMap)
	if err != nil {
		glog.Error("Error while marshalling: ", err)
		return nil
	}

	return ctxBytes
}

// Starts a child span for context stored in byte array.
//
// Usage example:
//  func A(ctx context.Context) {
//    span, ctx := StartSpan(ctx, "test-span")
//    defer span.Finish()
//    ctxToBytes := InjectCtxIntoBytes(span)
//    ...
//    B(ctxToBytes)
//    ...
//  }
//
//  func B(ctxToBytes []byte) {
//    ctx := context.Background()
//    span, ctx := StartSpanFromBytesCtx(ctxToBytes, "new-span", ctx)
//    defer span.Finish()
//    ...
//  }

// Args:
//  ctxBytes: Byte array containg span details.
//  spanName: Name of span to be created.
//  ctx: Context.
// Returns: Span interface and context. If ctxBytes or ctx, returns empty span
//          and nil.

func StartSpanFromBytesCtx(ctxBytes []byte, spanName string,
	ctx context.Context) (SpanIfc, context.Context) {
	if ctxBytes == nil || ctx == nil {
		return &Span{}, nil
	}

	var ctxMap map[string]string
	err := json.Unmarshal(ctxBytes, &ctxMap)
	if err != nil {
		glog.Error("Error while unmarshalling: ", err)
		return StartSpan(ctx, spanName)
	}

	spanCtx := otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(ctxMap))

	if err != nil {
		glog.Error("Unable to extract context from context map: ", err)
		return StartSpan(ctx, spanName)
	}

	otSpan := trace.SpanFromContext(spanCtx)
	span := &Span{
		s: &otSpan,
	}
	return span, ctx
}

// Starts a span with follows from relationship using context stored in byte
// array.
//
// Usage example:
//  func A(ctx context.Context) {
//    span, ctx := StartSpan(ctx, "test-span")
//    defer span.Finish()
//    ctxToBytes := InjectCtxIntoBytes(span)
//    ...
//    B(ctxToBytes)
//    ...
//  }
//
//  func B(ctxToBytes []byte) {
//    ctx := context.Background()
//    span, ctx := StartFollowsSpanFromBytesCtx(ctxToBytes, "new-span", ctx)
//    defer span.Finish()
//    ...
//  }
//

// Args:
//  ctxBytes: Byte array containg span details.
//  spanName: Name of span to be created.
//  ctx: Context.
// Returns: Span interface and context. If ctxBytes or ctx is nil, empty span
//          and nil is returned.
func StartFollowsSpanFromBytesCtx(ctxBytes []byte, spanName string,
	ctx context.Context) (SpanIfc, context.Context) {
	if ctxBytes == nil || ctx == nil {
		return &Span{}, nil
	}

	var ctxMap map[string]string
	err := json.Unmarshal(ctxBytes, &ctxMap)
	if err != nil {
		glog.Error("Error while unmarshalling: ", err)
		return StartSpan(ctx, spanName)
	}

	extractedSpanContext := otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(ctxMap))

	linkFromSpan := trace.Link{SpanContext: trace.SpanContextFromContext(extractedSpanContext)}
	links := []trace.Link{linkFromSpan}
	spanContext, otelSpan := tracer.Start(extractedSpanContext, spanName, trace.WithLinks(links...))

	otelSpan.End()
	span := &Span{
		s: &otelSpan,
	}

	return span, spanContext

}

// TODO: Modify to opentel
// Starts span using context in HTTP request headers.

// Usage example:
//  func RequestHandler(req *http.Request, res http.ResponseWriter) {
//    span, ctx := StartSpanFromHTTPHeader(request, "call-received")
//    defer span.Finish()
//    ...
//  }

// Args:
//  request: HTTP request.
//  spanName: Name of span to be created.
func StartSpanFromHTTPHeader(request *http.Request, spanName string) (SpanIfc,
	context.Context) {

	spanCtx := otel.GetTextMapPropagator().Extract(request.Context(), propagation.HeaderCarrier(request.Header))

	if spanCtx == nil {
		glog.V(2).Info("Unable to extract context from HTTP headers")
		return StartSpan(request.Context(), spanName)
	}

	//TODO: find opentelemetry method parameter equivalent to
	//      ext.RPCServerOption(spanCtx) which is for opentracing
	spanContext, otelSpan := tracer.Start(request.Context(), spanName, trace.WithSpanKind(trace.SpanKindServer))

	// In the span value of span.kind tag is set as server.
	otelSpan.End()
	span := &Span{
		s: &otelSpan,
	}

	return span, spanContext
}
