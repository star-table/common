package tracing

import (
    "github.com/star-table/common/core/config"
    "github.com/opentracing/opentracing-go"
)

func EnableTracing() bool{
    return config.GetJaegerConfig() != nil
}

func StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
    return opentracing.GlobalTracer().StartSpan(operationName, opts...)
}

func Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
    return opentracing.GlobalTracer().Inject(sm, format, carrier)
}

func Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
    return opentracing.GlobalTracer().Extract(format, carrier)
}
