package threadlocal

import (
	"github.com/galaxy-book/common/core/consts"
	"github.com/galaxy-book/common/core/model"
	"github.com/galaxy-book/common/core/util/uuid"
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

//往threadlocal中设置UUID类型的 traceId
func SetTraceId() {

	Mgr.SetValues(gls.Values{consts.TraceIdKey: uuid.NewUuid()}, func() {})
}
