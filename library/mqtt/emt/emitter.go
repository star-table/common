package emt

import (
	"errors"
	"gitea.bjx.cloud/allstar/common/core/config"
	"gitea.bjx.cloud/allstar/common/core/logger"
	"gitea.bjx.cloud/allstar/common/core/util/json"
	emitter "github.com/emitter-io/go/v2"
	"sync"
)

var (
	mqttClient *emitter.Client
	mqttMutex sync.Mutex
	log = logger.GetDefaultLogger()
	disConnectErr = errors.New("mqtt disconnected!")
	)

func GetClient() (*emitter.Client, error){
	client, err := Connect(nil)
	if err != nil{
		log.Error(err)
		return nil, err
	}
	//断开连接，重试一次
	if ! client.IsConnected(){
		mqttClient = nil
		client, err = Connect(nil)
		if err != nil{
			log.Error(err)
			return nil, err
		}
		if ! client.IsConnected(){
			return nil, disConnectErr
		}
	}
	return client, nil
}

func Connect(handler emitter.MessageHandler, options ...func(*emitter.Client)) (*emitter.Client, error){
	if mqttClient == nil {
		mqttMutex.Lock()
		defer mqttMutex.Unlock()
		if mqttClient == nil {
			mqttConfig := config.GetMQTTConfig()
			if mqttConfig == nil{
				panic("missing mqtt config.")
			}

			log.Infof("mqtt config %s", json.ToJsonIgnoreError(mqttConfig))

			if mqttConfig.Enable{
				var err error
				mqttClient, err = emitter.Connect(mqttConfig.Host, handler, options...)
				if err != nil{
					log.Error(err)
					return nil, err
				}
			}
		}

	}
	return mqttClient, nil
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