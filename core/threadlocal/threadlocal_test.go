package threadlocal

import (
	"github.com/galaxy-book/common/core/consts"
	"github.com/galaxy-book/common/core/util/uuid"
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