package vin

import (
	"net/http"
	"net/url"

	otelglobal "go.opentelemetry.io/otel/api/global"
	otelpropagation "go.opentelemetry.io/otel/api/propagation"
	oteltrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/semconv"
)

const (
	tracerKey = "otel-go-contrib-tracer"
)

const defaultComponentName = "net/http"

type mwOptions struct {
	opNameFunc    func(r *http.Request) string
	spanObserver  func(span oteltrace.Span, r *http.Request)
	urlTagFunc    func(u *url.URL) string
	componentName string
}

// MWOption controls the behavior of the Middleware.
type MWOption func(*mwOptions)

// OperationNameFunc returns a MWOption that uses given function f
// to generate operation name for each server-side span.
func OperationNameFunc(f func(r *http.Request) string) MWOption {
	return func(options *mwOptions) {
		options.opNameFunc = f
	}
}

// MWComponentName returns a MWOption that sets the component name
// for the server-side span.
func MWComponentName(componentName string) MWOption {
	return func(options *mwOptions) {
		options.componentName = componentName
	}
}

// MWSpanObserver returns a MWOption that observe the span
// for the server-side span.
func MWSpanObserver(f func(span oteltrace.Span, r *http.Request)) MWOption {
	return func(options *mwOptions) {
		options.spanObserver = f
	}
}

// MWURLTagFunc returns a MWOption that uses given function f
// to set the span's http.url tag. Can be used to change the default
// http.url tag, eg to redact sensitive information.
func MWURLTagFunc(f func(u *url.URL) string) MWOption {
	return func(options *mwOptions) {
		options.urlTagFunc = f
	}
}

var (
	_ignorePaths = map[string]bool{
		"/monitor/ping": true,
		"/favicon.ico":  true,
	}
)

// Middleware is a vin native version of the equivalent middleware in:
func Trace(service string, options ...MWOption) HandlerFunc {
	opts := mwOptions{
		opNameFunc: func(r *http.Request) string {
			return r.URL.Path + " " + r.Method
		},
		spanObserver: func(span oteltrace.Span, r *http.Request) {},
		urlTagFunc: func(u *url.URL) string {
			return u.String()
		},
	}
	for _, opt := range options {
		opt(&opts)
	}

	return func(c *Context) {
		tracer := otelglobal.Tracer(service)
		c.Set(tracerKey, tracer)

		if _ignorePaths[c.Request.URL.Path] {
			c.Next()
			return
		}

		ctx := otelpropagation.ExtractHTTP(c, otelglobal.Propagators(), c.Request.Header)
		startOpts := []oteltrace.StartOption{
			oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", c.Request)...),
			oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(c.Request)...),
			oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(service, c.FullPath(), c.Request)...),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}
		op := opts.opNameFunc(c.Request)
		ctx, span := tracer.Start(ctx, op, startOpts...)
		defer span.End()

		c.Context = oteltrace.ContextWithSpan(c.Context, span)

		c.Next()

		status := c.Writer.Status()
		attrs := semconv.HTTPAttributesFromHTTPStatusCode(status)
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(status)
		span.SetAttributes(attrs...)
		span.SetStatus(spanStatus, spanMessage)
		if len(c.Errors) > 0 {
			span.SetAttributes(label.String("vin.errors", c.Errors.String()))
		}
	}
}
