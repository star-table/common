package mysql

import (
	"gitea.bjx.cloud/allstar/common/core/logger"
	"gitea.bjx.cloud/allstar/common/core/util/strs"
	"sync"

	"gitea.bjx.cloud/allstar/common/core/config"
	"upper.io/db.v3/lib/sqlbuilder"
	upper "upper.io/db.v3/mysql"

	"errors"
	"strconv"
)

var mysqlMutex sync.Mutex
var settings *upper.ConnectionURL

func GetConnect() (sqlbuilder.Database, error) {
	if config.GetMysqlConfig() == nil {
		panic(errors.New("Mysql Datasource Configuration is missing!"))
	}

	if settings == nil {
		mysqlMutex.Lock()
		defer mysqlMutex.Unlock()
		if settings == nil {
			mc := config.GetMysqlConfig()
			settings = &upper.ConnectionURL{
				User:     mc.Usr,
				Password: mc.Pwd,
				Database: mc.Database,
				Host:     mc.Host + ":" + strconv.Itoa(mc.Port),
				Socket:   "",
				Options: map[string]string{
					"parseTime": "true",
					"loc":       "Asia/Shanghai",
				},
			}
		}
	}

	sess, err := upper.Open(settings)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

type Domain interface {
	TableName() string
}

func Close(conn sqlbuilder.Database, tx sqlbuilder.Tx) {
	if err := conn.Close(); err != nil {
		logger.GetDefaultLogger().Error(strs.ObjectToString(err))
	}
	if err := tx.Close(); err != nil {
		logger.GetDefaultLogger().Error(strs.ObjectToString(err))
	}
}

func Rollback(tx sqlbuilder.Tx) {
	err := tx.Rollback()
	if err != nil {
		logger.GetDefaultLogger().Error("Rollback error " + strs.ObjectToString(err))
	}
}

type Upd map[string]interface{}
