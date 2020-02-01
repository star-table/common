package mysql

import (
	"github.com/galaxy-book/common/core/config"
	"testing"
)

func TestSelectByQuery(t *testing.T) {
	config.LoadUnitTestConfig()

	//s := &[]domains.PpmProProjectRelation{}
	s := &[]int64{}
	err := SelectByQuery("SELECT id FROM ppm_pro_project_relation WHERE relation_id = 1007 AND is_delete = 2", s)
	t.Log(err)
	t.Log(len(*s))

}
