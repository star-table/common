package threadlocal

import (
	"github.com/jtolds/gls"
	"github.com/star-table/common/core/consts"
	"github.com/star-table/common/core/util/uuid"
	"testing"
)

func TestSetTraceId(t *testing.T) {

	SetTraceId()
	t.Log(GetTraceId())

	Mgr.SetValues(gls.Values{consts.TraceIdKey: uuid.NewUuid()}, func() {

		t.Log("in ", GetTraceId())
	})
}
