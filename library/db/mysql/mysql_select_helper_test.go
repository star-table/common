package mysql

import (
	"fmt"
	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/util/json"
	"testing"
	"time"
)

func TestSelectByQuery(t *testing.T) {
	config.LoadUnitTestConfig()

	//s := &[]domains.PpmProProjectRelation{}

	//i := true
	//
	//for ;;{
	//	if i{
	//		s := &[]map[string]interface{}{}
	//		err := SelectByQuery("SELECT * FROM ppm_pri_issue limit 1", s)
	//		fmt.Println(err)
	//		fmt.Println(len(*s))
	//		fmt.Println(json.ToJsonIgnoreError(s))
	//
	//
	//		err = SelectByQuery("SELECT * FROM ppm_pri_issue limit 2", s)
	//		fmt.Println(err)
	//		fmt.Println(len(*s))
	//		fmt.Println(json.ToJsonIgnoreError(s))
	//
	//
	//		err = SelectByQuery("SELECT * FROM ppm_pri_issue limit 3", s)
	//		fmt.Println(err)
	//		fmt.Println(len(*s))
	//		fmt.Println(json.ToJsonIgnoreError(s))
	//
	//		time.Sleep(2 * time.Millisecond)
	//	}
	//}
	//
	//go func() {
	//	time.Sleep(60)
	//	i = false
	//}()

	s := &[]map[string]interface{}{}
	err := SelectByQuery("SELECT * FROM ppm_pri_issue limit 1", s)
	fmt.Println(err)
	fmt.Println(len(*s))
	fmt.Println(json.ToJsonIgnoreError(s))

	err = SelectByQuery("SELECT * FROM ppm_pri_issue limit 2", s)
	fmt.Println(err)
	fmt.Println(len(*s))
	fmt.Println(json.ToJsonIgnoreError(s))

	err = SelectByQuery("SELECT * FROM ppm_pri_issue limit 3", s)
	fmt.Println(err)
	fmt.Println(len(*s))
	fmt.Println(json.ToJsonIgnoreError(s))

	time.Sleep(2 * time.Millisecond)

}
