package emt

import (
	"errors"
	"gitea.bjx.cloud/allstar/common/core/config"
	"gitea.bjx.cloud/allstar/common/core/lock"
	"gitea.bjx.cloud/allstar/common/core/logger"
	"gitea.bjx.cloud/allstar/common/core/util/json"
	emitter "github.com/emitter-io/go/v2"
	"strconv"
	"sync"
	"time"
)

var (
	mqttMutex sync.Mutex
	log = logger.GetDefaultLogger()
	disConnectErr = errors.New("mqtt disconnected!")
	clients []*emitter.Client = nil
	)

func GetClient() (*emitter.Client, error){
	client, selector, err := Connect(nil)
	if err != nil{
		log.Error(err)
		return nil, err
	}
	if client == nil{
		return nil, disConnectErr
	}

	//fmt.Printf("selector %d\n", selector)

	//断开连接，重试一次
	if ! client.IsConnected(){
		log.Infof("连接断开，开始重连...")

		key := strconv.Itoa(selector)
		lock.Lock(key)
		defer lock.Unlock(key)
		client = clients[selector]
		if client == nil || ! client.IsConnected(){
			clients[selector] = nil
			client, _, err = Connect(nil)
			if err != nil{
				log.Error(err)
				return nil, err
			}
			if ! client.IsConnected(){
				return nil, disConnectErr
			}
		}
	}
	return client, nil
}

func Connect(handler emitter.MessageHandler, options ...func(*emitter.Client)) (*emitter.Client, int, error){
	initPool()

	selector := int(time.Now().Unix() & int64(len(clients) - 1))
	selectClient := clients[selector]

	if selectClient == nil {
		mqttMutex.Lock()
		defer mqttMutex.Unlock()

		selectClient = clients[selector]
		if selectClient == nil {
			mqttConfig := config.GetMQTTConfig()
			if mqttConfig == nil{
				panic("missing mqtt config.")
			}

			log.Infof("mqtt config %s", json.ToJsonIgnoreError(mqttConfig))

			if mqttConfig.Enable{
				var err error
				clients[selector], err = emitter.Connect(mqttConfig.Host, handler, options...)
				if err != nil{
					log.Error(err)
					return nil, selector, err
				}
			}else{
				log.Error("mqtt已被禁用")
			}
		}

	}
	return clients[selector], selector, nil
}

func initPool(){
	mqttConfig := config.GetMQTTConfig()
	if mqttConfig == nil{
		panic("missing mqtt config.")
	}

	if clients == nil{
		mqttMutex.Lock()
		defer mqttMutex.Unlock()
		if clients == nil{
			poolSize := mqttConfig.ConnectPoolSize
			if poolSize <= 0{
				poolSize = 10
			}
			clients = make([]*emitter.Client, poolSize)
		}
	}
}

func GetNativeConnect(handler emitter.MessageHandler, options ...func(*emitter.Client)) (*emitter.Client, error){
	mqttConfig := config.GetMQTTConfig()
	if mqttConfig == nil{
		return nil, errors.New("missing mqtt config.")
	}

	log.Infof("mqtt config %s", json.ToJsonIgnoreError(mqttConfig))

	mc, err := emitter.Connect(mqttConfig.Host, handler, options...)
	if err != nil{
		log.Error(err)
		return nil, err
	}
	return mc, nil
}

func GenerateKey(channel string, permissions string, ttl int) (string, error){
	client, err := GetClient()
	if err != nil{
		log.Error(err)
		return "", err
	}
	mqttConfig := config.GetMQTTConfig()
	if mqttConfig == nil{
		panic("missing mqtt config.")
	}
	secretKey := mqttConfig.SecretKey

	key, err := client.GenerateKey(secretKey, channel, permissions, ttl)
	if err != nil{
		log.Error(err)
		return "", err
	}
	return key, nil
}

func Publish(key, channel string, payload interface{}, handler emitter.ErrorHandler) error{
	client, err := GetClient()
	if err != nil{
		log.Error(err)
		return err
	}
	if handler != nil{
		client.OnError(handler)
	}
	return client.Publish(key, channel, payload, emitter.WithAtLeastOnce())
}