package emt

import (
	"sync/atomic"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/errors"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/json"
	emitter "github.com/star-table/emitter-go-client"
)

var log = logger.GetDefaultLogger()
var _clients []*emitter.Client = nil
var _next uint32 = 0

func Init() error {
	conf := config.GetMQTTConfig()
	if conf == nil {
		log.Error("[MQTT] cannot get config.")
		return errors.MQTTNoConfigError
	}
	if !conf.Enable {
		log.Info("[MQTT] disabled.")
		return nil
	}

	log.Infof("[MQTT] init config: %q", json.ToJsonBytesIgnoreError(conf))

	clientCount := conf.ConnectPoolSize
	if clientCount <= 0 {
		clientCount = 10
	}
	_clients = make([]*emitter.Client, clientCount, clientCount)

	// init each client
	for i := 0; i < clientCount; i++ {
		c, err := emitter.Connect(conf.Address, nil)
		if err != nil || !c.IsConnected() {
			log.Errorf("[MQTT] client %v init failed: %v, err: %v", i, conf.Address, err)
			return err
		}
		_clients[i] = c
	}
	return nil
}

// RoundRobin get next index
func nextI() int {
	total := len(_clients)
	if total == 0 {
		return -1
	}

	n := atomic.AddUint32(&_next, 1)
	return (int(n) - 1) % total
}

func GetClient() (*emitter.Client, error) {
	var client *emitter.Client
	var i int
	tryTimes := 0
	conf := config.GetMQTTConfig()
	if conf == nil {
		log.Error("[MQTT] get client failed, cannot get config.")
		return nil, errors.MQTTNoConfigError
	}

RETRY:
	tryTimes += 1
	if tryTimes > len(_clients) {
		log.Errorf("[MQTT] get client failed, retry %v times!!!", tryTimes)
		return nil, errors.MQTTNoClientError
	}

	i = nextI()
	if i == -1 {
		// 理论上不可能走进的分支
		log.Errorf("[MQTT] get client failed, no client.")
		return nil, errors.MQTTNoClientError
	}
	client = _clients[i]

	// 理论上不可能走进的分支
	if client == nil {
		log.Errorf("[MQTT] get client failed, client %v is nil.", i)
		return nil, errors.MQTTNoClientError
	}

	// 理论上不可能走进的分支
	if !client.IsConnected() {
		log.Errorf("[MQTT] client [%v] is not connected, use next.", i)
		goto RETRY
	}

	// 默认自动重连是开启的，如果IsConnectionOpen=False则意味着底层连接正在自动重新连接状态
	// 此时跳过此连接，直接使用下一个连接
	if !client.IsConnectionOpen() {
		log.Infof("[MQTT] client [%v] connection is not open, use next.", i)
		goto RETRY
	}

	return client, nil
}

func GetNativeConnect(handler emitter.MessageHandler, options ...func(*emitter.Client)) (*emitter.Client, error) {
	conf := config.GetMQTTConfig()
	if conf == nil {
		return nil, errors.MQTTNoConfigError
	}

	log.Infof("[MQTT] GetNativeConnect config %q", json.ToJsonBytesIgnoreError(conf))

	c, err := emitter.Connect(conf.Address, handler, options...)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return c, nil
}

// 先搞明白客户端目前请求GenerateKey的频率, KEY的TTL设置是否必要？REDIS缓存要加上
func GenerateKey(channel string, permissions string, ttl int) (string, error) {
	client, err := GetClient()
	if err != nil {
		log.Errorf("[MQTT] GenerateKey channel: %v, permissions: %s, ttl: %d, cannot get client: %v",
			channel, permissions, ttl, err)
		return "", err
	}
	conf := config.GetMQTTConfig()
	if conf == nil {
		log.Errorf("[MQTT] GenerateKey channel: %v, permissions: %s, ttl: %d, no config.",
			channel, permissions, ttl)
		return "", errors.MQTTNoConfigError
	}

	key, err := client.GenerateKey(conf.SecretKey, channel, permissions, ttl)
	if err != nil {
		log.Errorf("[MQTT] GenerateKey channel: %v, permissions: %s, ttl: %d, failed: %v.",
			channel, permissions, ttl, err)
		return "", err
	}
	return key, nil
}

func Publish(key, channel string, payload interface{}, ttl int, handler emitter.ErrorHandler) error {
	client, err := GetClient()
	if err != nil {
		log.Errorf("[MQTT] Publish key: %v, channel: %v, payload: %v, cannot get client: %v.",
			key, channel, payload, err)
		return err
	}
	if handler != nil {
		client.OnError(handler)
	}
	if ttl > 0 {
		go func() {
			if err = client.Publish(key, channel, payload, emitter.WithAtLeastOnce(), emitter.WithTTL(ttl)); err != nil {
				log.Errorf("[MQTT] Publish key: %v, channel: %v, payload: %v, failed: %v.",
					key, channel, payload, err)
			}
		}()
	} else {
		go func() {
			if err = client.Publish(key, channel, payload, emitter.WithAtLeastOnce()); err != nil {
				log.Errorf("[MQTT] Publish key: %v, channel: %v, payload: %v, failed: %v.",
					key, channel, payload, err)
			}
		}()
	}
	return nil
}
