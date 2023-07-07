package mysql

import (
	"context"
	"database/sql"
	"sync"
	"time"

	upperPg "upper.io/db.v3/postgresql"

	"github.com/go-sql-driver/mysql"
	"github.com/opentracing/opentracing-go"
	"github.com/qustavo/sqlhooks/v2"
	"github.com/star-table/common/core/consts"
	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/library/tracing"

	"errors"
	"strconv"

	"github.com/star-table/common/core/config"
	"upper.io/db.v3/lib/sqlbuilder"
	upper "upper.io/db.v3/mysql"
)

var mysqlMutex sync.Mutex
var sess sqlbuilder.Database

// Hooks satisfies the sqlhook.Hooks interface
type Hooks struct{}

// Before hook will print the query with it's args and return the context with the timestamp
func (h *Hooks) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	if tracing.EnableTracing() {
		if v, ok := threadlocal.Mgr.GetValue(consts.JaegerContextSpanKey); ok {
			if parentSpan, ok := v.(opentracing.Span); ok {
				spanCtx := parentSpan.Context()
				span := tracing.StartSpan("mysql opt", opentracing.ChildOf(spanCtx))
				span.SetTag("sql", query)
				span.SetTag("args", args)
				span.SetTag("operation", "mysql opt")
				return context.WithValue(ctx, "traceSpan", span), nil
			}
		}
	}
	return ctx, nil
}

// After hook will get the timestamp registered on the Before hook and print the elapsed time
func (h *Hooks) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	if tracing.EnableTracing() {
		if v := ctx.Value("traceSpan"); v != nil {
			if span, ok := v.(opentracing.Span); ok {
				span.Finish()
			}
		}
	}
	return ctx, nil
}

func init() {
	sql.Register("mysql-hooks", sqlhooks.Wrap(&mysql.MySQLDriver{}, &Hooks{}))
}

func GetConnect() (sqlbuilder.Database, error) {
	if config.GetMysqlConfig() == nil {
		panic(errors.New("Mysql Datasource Configuration is missing!"))
	}

	if sess == nil {
		mysqlMutex.Lock()
		defer mysqlMutex.Unlock()
		if sess == nil {
			var err error
			sess, err = InitSess()
			if err != nil {
				return nil, err
			}
		}
	}
	if err := sess.Ping(); err != nil {
		sess, err = InitSess()
		if err != nil {
			return nil, err
		}
	}
	return sess, nil
}

func InitSess() (sqlbuilder.Database, error) {
	mc := config.GetMysqlConfig()
	if mc.Driver == "pgsql" {
		settings := &upperPg.ConnectionURL{
			User:     mc.Usr,
			Password: mc.Pwd,
			Database: mc.Database,
			Host:     mc.Host + ":" + strconv.Itoa(mc.Port),
			Socket:   "",
			Options: map[string]string{
				"TimeZone": "Asia/Shanghai",
				"sslmode":  "disable",
			},
		}
		var err error
		sess, err = upperPg.Open(settings)
		if err != nil {
			return nil, err
		}
	} else {
		settings := &upper.ConnectionURL{
			User:     mc.Usr,
			Password: mc.Pwd,
			Database: mc.Database,
			Host:     mc.Host + ":" + strconv.Itoa(mc.Port),
			Socket:   "",
			Options: map[string]string{
				"parseTime": "true",
				"loc":       "Asia/Shanghai",
				"charset":   "utf8mb4",
				"collation": "utf8mb4_unicode_ci",
			},
		}
		var err error
		sess, err = upper.Open(settings)
		if err != nil {
			return nil, err
		}
	}

	maxOpenConns := 50
	maxIdleConns := 10
	maxLifetime := 300
	if mc.MaxOpenConns > 0 {
		maxOpenConns = mc.MaxOpenConns
	}
	if mc.MaxIdleConns > 0 {
		maxIdleConns = mc.MaxIdleConns
	}
	if mc.MaxLifetime > 0 {
		maxLifetime = mc.MaxLifetime
	}
	sess.SetMaxOpenConns(maxOpenConns)
	sess.SetMaxIdleConns(maxIdleConns)
	sess.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)
	return sess, nil
}

type Domain interface {
	TableName() string
}

func Close(conn sqlbuilder.Database, tx sqlbuilder.Tx) {
	//if conn != nil{
	//	if err := conn.Close(); err != nil {
	//		log.Error(strs.ObjectToString(err))
	//	}
	//}
	if tx != nil {
		if err := tx.Close(); err != nil {
			log.Error(strs.ObjectToString(err))
		}
	}
}

func Rollback(tx sqlbuilder.Tx) {
	err := tx.Rollback()
	if err != nil {
		log.Error("Rollback error " + strs.ObjectToString(err))
	}
}

type Upd map[string]interface{}
