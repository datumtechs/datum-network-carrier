package traceutil

import (
	"github.com/RosettaFlow/Carrier-Go/common/hashutil"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"go.opencensus.io/trace"
)

// AnnotateError on span.
//This should be used any time a particular span experiences an error.
func AnnotateError(span *trace.Span, err error) {
	if err == nil {
		return
	}
	span.AddAttributes(trace.BoolAttribute("error", true))
	span.SetStatus(trace.Status{
		Code:    trace.StatusCodeUnknown,
		Message: err.Error(),
	})
}

func GenerateParentSpan(pid peer.ID, msg *pubsub.Message) trace.SpanContext {
	h := hashutil.Hash([]byte(msg.String()))
	traceId := trace.TraceID{}
	copy(traceId[0:], h[:16])
	return trace.SpanContext{
		TraceID:      traceId,
		SpanID:       trace.SpanID{0, 1, 2, 3, 4, 5, 6, 7},
		TraceOptions: 1,
	}
}