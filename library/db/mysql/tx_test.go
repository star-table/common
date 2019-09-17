package mysql

import (
	"fmt"
	"gitea.bjx.cloud/allstar/common/core/config"
	"testing"
	"upper.io/db.v3/lib/sqlbuilder"
)

func TestTransX(t *testing.T) {
	config.LoadConfig("F:\\polaris-backend\\polaris-server\\config", "application")
	a := []int64{1}
	TransX(func(tx sqlbuilder.Tx) error {
		a = append(a, 2)
		return nil
	})
	fmt.Println(a)

}
