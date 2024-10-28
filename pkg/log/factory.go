package log

import (
	"ascale/pkg/conf/env"
	"fmt"

	"go.opentelemetry.io/otel/api/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string)
	Debugf(msg string, a ...interface{})
	DebugWithFields(msg string, fields ...zapcore.Field)

	Info(msg string)
	Infof(msg string, a ...interface{})
	InfoWithFields(msg string, fields ...zapcore.Field)

	Warn(msg string)
	Warnf(msg string, a ...interface{})
	WarnWithFields(msg string, fields ...zapcore.Field)

	Error(msg string)
	Errorf(msg string, a ...interface{})
	ErrorWithFields(msg string, fields ...zapcore.Field)

	DPanic(msg string)
	DPanicf(msg string, a ...interface{})
	DPanicWithFields(msg string, fields ...zapcore.Field)

	Panic(msg string)
	Panicf(msg string, a ...interface{})
	PanicWithFields(msg string, fields ...zapcore.Field)

	Fatal(msg string)
	Fatalf(msg string, a ...interface{})
	FatalWithFields(msg string, fields ...zapcore.Field)
}

// logger delegates all calls to the underlying zap.Logger
type logger struct {
	logger *zap.Logger
}

func (l logger) Debug(msg string) {
	fields := appendFields()
	l.logger.Debug(msg, fields...)
}

func (l logger) Debugf(msg string, a ...interface{}) {
	fields := appendFields()
	l.logger.Debug(fmt.Sprintf(msg, a...), fields...)
}

func (l logger) DebugWithFields(msg string, fields ...zapcore.Field) {
	fields = append(fields, appendFields()...)
	l.logger.Debug(msg, fields...)
}

// Info logs a message at InfoLevel
func (l logger) Info(msg string) {
	fields := appendFields()
	l.logger.Info(msg, fields...)
}

func (l logger) Infof(msg string, a ...interface{}) {
	fields := appendFields()
	l.logger.Info(fmt.Sprintf(msg, a...), fields...)
}

func (l logger) InfoWithFields(msg string, fields ...zapcore.Field) {
	fields = append(fields, appendFields()...)
	l.logger.Info(msg, fields...)
}

func (l logger) Warn(msg string) {
	fields := appendFields()
	l.logger.Warn(msg, fields...)
}

func (l logger) Warnf(msg string, a ...interface{}) {
	fields := appendFields()
	l.logger.Warn(fmt.Sprintf(msg, a...), fields...)
}

func (l logger) WarnWithFields(msg string, fields ...zapcore.Field) {
	fields = append(fields, appendFields()...)
	l.logger.Warn(msg, fields...)
}

func (l logger) Error(msg string) {
	fields := appendFields()
	l.logger.Error(msg, fields...)
}

func (l logger) Errorf(msg string, a ...interface{}) {
	fields := appendFields()
	l.logger.Error(fmt.Sprintf(msg, a...), fields...)
}

func (l logger) ErrorWithFields(msg string, fields ...zapcore.Field) {
	fields = append(fields, appendFields()...)
	l.logger.Error(msg, fields...)
}

func (l logger) DPanic(msg string) {
	fields := appendFields()
	l.logger.DPanic(msg, fields...)
}

func (l logger) DPanicf(msg string, a ...interface{}) {
	fields := appendFields()
	l.logger.DPanic(fmt.Sprintf(msg, a...), fields...)
}

func (l logger) DPanicWithFields(msg string, fields ...zapcore.Field) {
	fields = append(fields, appendFields()...)
	l.logger.DPanic(msg, fields...)
}

func (l logger) Panic(msg string) {
	fields := appendFields()
	l.logger.Panic(msg, fields...)
}

func (l logger) Panicf(msg string, a ...interface{}) {
	fields := appendFields()
	l.logger.Panic(fmt.Sprintf(msg, a...), fields...)
}

func (l logger) PanicWithFields(msg string, fields ...zapcore.Field) {
	fields = append(fields, appendFields()...)
	l.logger.Panic(msg, fields...)
}

func (l logger) Fatal(msg string) {
	fields := appendFields()
	l.logger.Fatal(msg, fields...)
}

func (l logger) Fatalf(msg string, a ...interface{}) {
	fields := appendFields()
	l.logger.Fatal(fmt.Sprintf(msg, a...), fields...)
}

func (l logger) FatalWithFields(msg string, fields ...zapcore.Field) {
	fields = append(fields, appendFields()...)
	l.logger.Fatal(msg, fields...)
}

// spanLogger
type spanLogger struct {
	logger *zap.Logger
	span   trace.Span
}

func (l spanLogger) Debug(msg string) {
	spanCtx := l.span.SpanContext()
	fields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	l.logger.Debug(msg, fields...)
}

func (l spanLogger) Debugf(msg string, a ...interface{}) {
	spanCtx := l.span.SpanContext()
	fields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	l.logger.Debug(fmt.Sprintf(msg, a...), fields...)
}

func (l spanLogger) DebugWithFields(msg string, fields ...zapcore.Field) {
	spanCtx := l.span.SpanContext()
	traceFields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	fields = append(fields, traceFields...)
	l.logger.Debug(msg, fields...)
}

// Info logs a message at InfoLevel
func (l spanLogger) Info(msg string) {
	spanCtx := l.span.SpanContext()
	fields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	l.logger.Info(msg, fields...)
}

func (l spanLogger) Infof(msg string, a ...interface{}) {
	spanCtx := l.span.SpanContext()
	fields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	l.logger.Info(fmt.Sprintf(msg, a...), fields...)
}

func (l spanLogger) InfoWithFields(msg string, fields ...zapcore.Field) {
	spanCtx := l.span.SpanContext()
	traceFields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	fields = append(fields, traceFields...)
	l.logger.Info(msg, fields...)
}

func (l spanLogger) Warn(msg string) {
	spanCtx := l.span.SpanContext()
	fields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	l.logger.Warn(msg, fields...)
}

func (l spanLogger) Warnf(msg string, a ...interface{}) {
	spanCtx := l.span.SpanContext()
	fields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	l.logger.Warn(fmt.Sprintf(msg, a...), fields...)
}

func (l spanLogger) WarnWithFields(msg string, fields ...zapcore.Field) {
	spanCtx := l.span.SpanContext()
	traceFields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	fields = append(fields, traceFields...)
	l.logger.Warn(msg, fields...)
}

func (l spanLogger) Error(msg string) {
	spanCtx := l.span.SpanContext()
	fields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	l.logger.Error(msg, fields...)
}

func (l spanLogger) Errorf(msg string, a ...interface{}) {
	spanCtx := l.span.SpanContext()
	fields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	l.logger.Error(fmt.Sprintf(msg, a...), fields...)
}

func (l spanLogger) ErrorWithFields(msg string, fields ...zapcore.Field) {
	spanCtx := l.span.SpanContext()
	traceFields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	fields = append(fields, traceFields...)
	l.logger.Error(msg, fields...)
}

func (l spanLogger) DPanic(msg string) {
	spanCtx := l.span.SpanContext()
	fields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	l.logger.DPanic(msg, fields...)
}

func (l spanLogger) DPanicf(msg string, a ...interface{}) {
	spanCtx := l.span.SpanContext()
	fields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	l.logger.DPanic(fmt.Sprintf(msg, a...), fields...)
}

func (l spanLogger) DPanicWithFields(msg string, fields ...zapcore.Field) {
	spanCtx := l.span.SpanContext()
	traceFields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	fields = append(fields, traceFields...)
	l.logger.DPanic(msg, fields...)
}

func (l spanLogger) Panic(msg string) {
	spanCtx := l.span.SpanContext()
	fields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	l.logger.Panic(msg, fields...)
}

func (l spanLogger) Panicf(msg string, a ...interface{}) {
	spanCtx := l.span.SpanContext()
	fields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	l.logger.Panic(fmt.Sprintf(msg, a...), fields...)
}

func (l spanLogger) PanicWithFields(msg string, fields ...zapcore.Field) {
	spanCtx := l.span.SpanContext()
	traceFields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	fields = append(fields, traceFields...)
	l.logger.Panic(msg, fields...)
}

func (l spanLogger) Fatal(msg string) {
	spanCtx := l.span.SpanContext()
	fields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	l.logger.Fatal(msg, fields...)
}

func (l spanLogger) Fatalf(msg string, a ...interface{}) {
	spanCtx := l.span.SpanContext()
	fields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	l.logger.Fatal(fmt.Sprintf(msg, a...), fields...)
}

func (l spanLogger) FatalWithFields(msg string, fields ...zapcore.Field) {
	spanCtx := l.span.SpanContext()
	traceFields := TraceContext(
		spanCtx.TraceID.String(),
		spanCtx.SpanID.String(),
		true,
		env.ProjectID,
	)
	fields = append(fields, appendFields()...)
	fields = append(fields, traceFields...)
	l.logger.Fatal(msg, fields...)
}
