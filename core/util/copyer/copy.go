package copyer

import (
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/json"
	"github.com/galaxy-book/common/core/util/strs"
	"github.com/pkg/errors"
)

var log = logger.GetDefaultLogger()

func Copy(src interface{}, source interface{}) error {
	jsonStr, err := json.ToJson(src)
	if err != nil {
		log.Error(strs.ObjectToString(err))
		return errors.New("json转换异常")
	}
	err = json.FromJson(jsonStr, source)
	if err != nil {
		log.Error(strs.ObjectToString(err))
	}
	return nil
}
