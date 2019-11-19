package threadlocal

import (
	"gitea.bjx.cloud/allstar/common/core/consts"
	"gitea.bjx.cloud/allstar/common/core/model"
	"gitea.bjx.cloud/allstar/common/core/util/uuid"
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
