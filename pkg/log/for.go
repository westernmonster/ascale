package log

import (
	"context"

	"go.opentelemetry.io/otel/api/trace"
)

func For(c context.Context) Logger {
	span := trace.SpanFromContext(c)
	if span.IsRecording() {
		return spanLogger{logger: l, span: span}
	}

	return logger{logger: l}
}

func GetLogger() Logger {
	return logger{logger: l}
}
