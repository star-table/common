package copyer

import (
	"gitea.bjx.cloud/allstar/common/core/logger"
	"gitea.bjx.cloud/allstar/common/core/util/json"
	"gitea.bjx.cloud/allstar/common/core/util/strs"
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
