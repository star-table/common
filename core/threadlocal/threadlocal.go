package threadlocal

import (
	"gitea.bjx.cloud/allstar/common/core/consts"
	"gitea.bjx.cloud/allstar/common/core/model"
	"github.com/jtolds/gls"
)

var (
	Mgr = gls.NewContextManager()
)

func GetHttpContext() *model.HttpContext {
	if traceId, ok := Mgr.GetValue(consts.HttpContextKey); ok {
		if traceId != nil {
			httpContext := traceId.(model.HttpContext)
			return &httpContext
		}
	}
	return nil
}

func GetTraceId() string {
	if traceId, ok := Mgr.GetValue(consts.TraceIdKey); ok {
		if traceId != nil {
			return traceId.(string)
		}
	}
	return ""
}
