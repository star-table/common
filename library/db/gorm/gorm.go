package gorm

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	gormopentracing "github.com/star-table/go-common/pkg/gorm/opentracing"
	"github.com/star-table/common/core/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var gormMutex sync.Mutex
var gormDB *gorm.DB

func GetDB() (*gorm.DB, error) {
	if config.GetMysqlConfig() == nil {
		panic(errors.New("Mysql Datasource Configuration is missing!"))
	}

	if gormDB == nil {
		gormMutex.Lock()
		defer gormMutex.Unlock()
		if gormDB == nil {
			var err error
			gormDB, err = initDB()
			if err != nil {
				return nil, err
			}
		}
	}
	if db, err := gormDB.DB(); err != nil || db.Ping() != nil {
		gormDB, err = initDB()
		if err != nil {
			return nil, err
		}
	}
	return gormDB, nil
}

func initDB() (*gorm.DB, error) {
	mc := config.GetMysqlConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Asia/Shanghai&collation=utf8mb4_unicode_ci",
		mc.Usr, mc.Pwd, mc.Host, mc.Port, mc.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
		return nil, err
	}

	err = db.Use(gormopentracing.New())
	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
		return nil, err
	}

	maxOpenConns := 100
	maxIdleConns := 32
	maxLifetime := 3600
	if mc.MaxOpenConns > 0 {
		maxOpenConns = mc.MaxOpenConns
	}
	if mc.MaxIdleConns > 0 {
		maxIdleConns = mc.MaxIdleConns
	}
	if mc.MaxLifetime > 0 {
		maxLifetime = mc.MaxLifetime
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)

	return db, nil
}
