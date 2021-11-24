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

func GenerateParentSpan(pid peer.ID, msg *pubsub.Message, spanId trace.SpanID) trace.SpanContext {
	h := hashutil.Hash([]byte(msg.String()))
	traceId := trace.TraceID{}
	copy(traceId[0:], h[:16])
	return trace.SpanContext{
		TraceID:      traceId,
		SpanID:       spanId,
		TraceOptions: 1,
	}
}

func GenerateParentSpanWithGossipTestData(pid peer.ID, msg *pubsub.Message) trace.SpanContext {
	return GenerateParentSpan(pid, msg, trace.SpanID{1, 1, 1, 1, 1, 1, 1, 1})
}

func GenerateParentSpanWithCommitMsg(pid peer.ID, msg *pubsub.Message) trace.SpanContext {
	return GenerateParentSpan(pid, msg, trace.SpanID{2, 2, 2, 2, 2, 2, 2, 2})
}

func GenerateParentSpanWithConfirmMsg(pid peer.ID, msg *pubsub.Message) trace.SpanContext {
	return GenerateParentSpan(pid, msg, trace.SpanID{3, 3, 3, 3, 3, 3, 3, 3})
}

func GenerateParentSpanWithConfirmVote(pid peer.ID, msg *pubsub.Message) trace.SpanContext {
	return GenerateParentSpan(pid, msg, trace.SpanID{4, 4, 4, 4, 4, 4, 4, 4})
}

func GenerateParentSpanWithPrepareMsg(pid peer.ID, msg *pubsub.Message) trace.SpanContext {
	return GenerateParentSpan(pid, msg, trace.SpanID{5, 5, 5, 5, 5, 5, 5, 5})
}

func GenerateParentSpanWithPrepareVote(pid peer.ID, msg *pubsub.Message) trace.SpanContext {
	return GenerateParentSpan(pid, msg, trace.SpanID{6, 6, 6, 6, 6, 6, 6, 6})
}