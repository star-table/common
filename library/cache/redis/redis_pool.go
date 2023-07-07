package redis

import (
	"errors"
	"sync"
	"time"

	"github.com/star-table/common/core/config"
	"github.com/gomodule/redigo/redis"

	"strconv"
)

var redisClient *redis.Pool

var poolMutex sync.Mutex

func initRedisPool() {
	if config.GetRedisConfig() == nil {
		panic(errors.New("Redis Configuration is missing!"))
	}

	rc := config.GetRedisConfig()

	if &rc == nil {
		return
	}

	maxIdle := 10
	if rc.MaxIdle > 0 {
		maxIdle = rc.MaxIdle
	}

	maxActive := 10
	if rc.MaxActive > 0 {
		maxActive = rc.MaxActive
	}

	maxIdleTimeout := 60
	if rc.MaxIdleTimeout > 0 {
		maxIdleTimeout = rc.MaxIdleTimeout
	}

	timeout := time.Duration(5)

	if rc.IsSentinel {
		sntnl := &Sentinel{
			Addrs:      []string{rc.Host+":"+strconv.Itoa(rc.Port)},
			MasterName: rc.MasterName,
			Dial: func(addr string) (redis.Conn, error) {
				c, err := redis.Dial("tcp", addr,
					redis.DialPassword(rc.SentinelPwd),
					redis.DialDatabase(rc.Database),
					redis.DialConnectTimeout(timeout*time.Second),
					redis.DialReadTimeout(timeout*time.Second),
					redis.DialWriteTimeout(timeout*time.Second))
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		}
		redisClient = &redis.Pool{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			IdleTimeout: time.Duration(maxIdleTimeout) * time.Second,
			Wait:        true,
			Dial: func() (redis.Conn, error) {
				masterAddr, err := sntnl.MasterAddr()
				if err != nil {
					return nil, err
				}
				c, err := redis.Dial("tcp", masterAddr,
					redis.DialPassword(rc.Pwd),
					redis.DialDatabase(rc.Database),
					redis.DialConnectTimeout(timeout*time.Second),
					redis.DialReadTimeout(timeout*time.Second),
					redis.DialWriteTimeout(timeout*time.Second))
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		}
		return
	}

	// 建立连接池
	redisClient = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(maxIdleTimeout) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", rc.Host+":"+strconv.Itoa(rc.Port),
				redis.DialPassword(rc.Pwd),
				redis.DialDatabase(rc.Database),
				redis.DialConnectTimeout(timeout*time.Second),
				redis.DialReadTimeout(timeout*time.Second),
				redis.DialWriteTimeout(timeout*time.Second))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}
}

func GetRedisConn() redis.Conn {
	if redisClient == nil {

		poolMutex.Lock()
		defer poolMutex.Unlock()
		if redisClient == nil {
			initRedisPool()
		}
	}
	return redisClient.Get()
}
