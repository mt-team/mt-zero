package middleware

import (
	"net/http"

	"mtzero/src/gateway/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type TraceMiddleware struct {
	Conf config.Config
}

func NewTraceMiddleware(conf config.Config) *TraceMiddleware {
	return &TraceMiddleware{
		Conf: conf,
	}
}

func (m *TraceMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		propagator := otel.GetTextMapPropagator()
		serviceName := m.Conf.Name
		tracer := otel.GetTracerProvider().Tracer(serviceName)
		path := r.RequestURI

		ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		spanCtx, span := tracer.Start(
			ctx,
			path,
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(
				serviceName, path, r)...),
		)
		defer span.End()

		sc := span.SpanContext()
		w.Header().Set("X-Trace-ID", sc.TraceID().String())

		r = r.WithContext(spanCtx)
		next(w, r)
	}
}
