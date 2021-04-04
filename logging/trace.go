package logging

import (
	"context"
	"fmt"
	"net/http"

	"go.opencensus.io/trace"
)

func WithTracing(ctx context.Context, h http.Handler) http.HandlerFunc {
	logger := FromContext(ctx)
	return func(w http.ResponseWriter, r *http.Request) {
		ctxReq := r.Context()
		h.ServeHTTP(w, r.WithContext(WithLogger(ctxReq,
			logger.With(traceFromContext(ctxReq)...),
		)))
	}
}

const (
	traceKey        = "logging.googleapis.com/trace"
	spanKey         = "logging.googleapis.com/spanId"
	traceSampledKey = "logging.googleapis.com/trace_sampled"
)

// traceFromContext adds the correct Stackdriver trace fields.
//
// see: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
func traceFromContext(ctx context.Context) []interface{} {
	span := trace.FromContext(ctx)

	if span == nil {
		return nil
	}

	sc := span.SpanContext()
	return []interface{}{
		traceKey, fmt.Sprintf("projects/%s/traces/%s", "viecco", sc.TraceID),
		spanKey, sc.SpanID.String(),
		traceSampledKey, sc.IsSampled(),
	}
}
