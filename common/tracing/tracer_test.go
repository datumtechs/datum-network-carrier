package tracing

import (
	"context"
	"errors"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/traceutil"
	"go.opencensus.io/trace"
	"testing"
	"time"
)

func TestJaeger(t *testing.T) {

	Setup(
		"carrier-network-4", // service name
		"p2p",
		"http://10.1.1.1:14268/api/traces",
		0.20,
		true,
	)
	ctx := context.TODO()
	for i := 0; i < 20; i++ {
		ctx, span := trace.StartSpan(ctx, fmt.Sprintf("sync.rpc"))
		//defer span.End()
		span.AddAttributes(trace.StringAttribute("topic",fmt.Sprintf("mytopic %d", i)))
		span.AddAttributes(trace.StringAttribute("peer", fmt.Sprintf("peerId %d", i)))
		traceutil.AnnotateError(span, errors.New(fmt.Sprintf("err %d", i)))
		time.Sleep(3*time.Second)
		span.End()

		ctx, span_1 := trace.StartSpan(ctx, fmt.Sprintf("sync.rpc_1 %d", i))
		span_1.AddAttributes(trace.StringAttribute("peer_1 tags", fmt.Sprintf("peerId_1 %d", i)))
		traceutil.AnnotateError(span_1, errors.New(fmt.Sprintf("err_1 %d", i)))
		time.Sleep(4*time.Second)
		span_1.End()

		_, span_2 := trace.StartSpan(ctx, fmt.Sprintf("sync.rpc_2 %d", i))
		span_2.AddAttributes(trace.StringAttribute("peer_2 tags", fmt.Sprintf("peerId_2 %d", i)))
		traceutil.AnnotateError(span_2, errors.New(fmt.Sprintf("err_2 %d", i)))
		time.Sleep(4*time.Second)
		span_2.End()
	}

	ch := make(chan int64, 1)
	go func() {
		time.Sleep(60 * time.Second)
		ch <- 10
	}()
	<- ch
}
