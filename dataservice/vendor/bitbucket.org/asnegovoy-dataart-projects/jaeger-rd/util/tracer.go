package util

import (
	"fmt"
	"io"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// InitJaeger returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func InitJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg, err := config.FromEnv()
	if err != nil {
		panic(fmt.Sprintf("Could not parse Jaeger env vars: %s", err.Error()))
	}

	cfg.ServiceName = service

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func GetSpanFromRPCReq(tracer opentracing.Tracer, r *http.Request, spanName string) opentracing.Span {
	carrier := opentracing.HTTPHeadersCarrier(r.Header)
	clientContext, err := tracer.Extract(opentracing.HTTPHeaders, carrier)

	var span opentracing.Span
	if err == nil {
		span = tracer.StartSpan(spanName, ext.RPCServerOption(clientContext))
	} else {
		span = tracer.StartSpan(spanName)
	}

	return span
}

func InjectSpanToReq(span opentracing.Span, req *http.Request) {
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, req.URL.String())
	ext.HTTPMethod.Set(span, req.Method)
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
}
// vim: tabstop=8 shiftwidth=8 expandtab!
