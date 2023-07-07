package mysql

import (
	"fmt"
	"github.com/star-table/common/core/config"
	"testing"
	"upper.io/db.v3/lib/sqlbuilder"
)

func TestTransX(t *testing.T) {
	config.LoadUnitTestConfig()
	a := []int64{1}
	TransX(func(tx sqlbuilder.Tx) error {
		a = append(a, 2)
		return nil
	})
	fmt.Println(a)

}
