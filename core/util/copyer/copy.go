package copyer

import (
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/strs"
	"github.com/pkg/errors"
)

var log = logger.GetDefaultLogger()

func Copy(src interface{}, source interface{}) error {
	jsonStr, err := json.Marshal(src)
	if err != nil {
		log.Error(strs.ObjectToString(err))
		return errors.New("json转换异常")
	}
	err = json.Unmarshal(jsonStr, source)
	if err != nil {
		log.Error(strs.ObjectToString(err))
	}
	return nil
}
