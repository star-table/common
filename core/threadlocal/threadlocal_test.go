package threadlocal

import (
	"github.com/star-table/common/core/consts"
	"github.com/star-table/common/core/util/uuid"
	"github.com/jtolds/gls"
	"testing"
)

func TestSetTraceId(t *testing.T) {

	SetTraceId()
	t.Log(GetTraceId())

	Mgr.SetValues(gls.Values{consts.TraceIdKey: uuid.NewUuid()}, func() {

		t.Log("in ", GetTraceId())
	})
}