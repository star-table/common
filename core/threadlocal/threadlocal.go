package threadlocal

import (
	"gitea.bjx.cloud/allstar/common/core/model"
	"github.com/jtolds/gls"
)

// context key
const (
	TraceIdKey     = "_traceId"
	HttpContextKey = "_httpContext"
)

var (
	Mgr = gls.NewContextManager()
)

func GetHttpContext() *model.HttpContext {
	if traceId, ok := Mgr.GetValue(HttpContextKey); ok {
		if traceId != nil {
			httpContext := traceId.(model.HttpContext)
			return &httpContext
		}
	}
	return nil
}

func GetTraceId() string {
	if traceId, ok := Mgr.GetValue(TraceIdKey); ok {
		if traceId != nil {
			return traceId.(string)
		}
	}
	return ""
}
