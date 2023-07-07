package nacos

import (
	"fmt"
	"os"
	"strconv"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/http"
	"github.com/star-table/common/core/util/network"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/pkg/errors"
)

var (
	nacosClient *NacosClient
	log         = logger.GetDefaultLogger()
)

func Init() {
	if config.GetConfig().Nacos == nil {
		log.Info("nacos not config.")
		return
	}
	client, err := NewNacosClient(config.GetConfig().Nacos)
	if err != nil {
		log.Fatal(err)
		return
	}
	nacosClient = client
	host := config.GetConfig().Server.Host
	if host == "" {
		host = network.GetIntranetIp()
	}

	metaData := config.GetConfig().Nacos.Discovery.MetaData
	if metaData == nil {
		metaData = map[string]string{
			"kind":    "http",
			"version": "",
		}
	}

	if config.GetConfig().Nacos.Discovery.Enable {
		suc, err := nacosClient.RegisterInstance(vo.RegisterInstanceParam{
			Ip:          host,
			Port:        uint64(config.GetConfig().Server.Port),
			ServiceName: config.GetConfig().Application.Name,
			GroupName:   config.GetConfig().Nacos.Discovery.GroupName,
			ClusterName: config.GetConfig().Nacos.Discovery.ClusterName,
			Weight:      config.GetConfig().Nacos.Discovery.Weight,
			Enable:      config.GetConfig().Nacos.Discovery.Enable,
			Healthy:     config.GetConfig().Nacos.Discovery.Healthy,
			Ephemeral:   config.GetConfig().Nacos.Discovery.Ephemeral,
			Metadata:    metaData,
		})
		if err != nil {
			log.Fatal(err)
			return
		}
		if !suc {
			log.Fatal("服务注册失败")
			return
		}
	}
}

type NacosClient struct {
	namingClient naming_client.INamingClient
	configClient config_client.IConfigClient
}

var (
	emptyClientConf = config.NacosClientConfig{}
	defaultClient   = constant.ClientConfig{
		TimeoutMs:      10 * 1000,
		ListenInterval: 30 * 1000,
		BeatInterval:   5 * 1000,
	}
)

func GetNacosClient() *NacosClient {
	return nacosClient
}

func NewNacosClient(conf *config.NacosConfig) (*NacosClient, error) {
	if len(conf.Server) == 0 {
		return nil, errors.Errorf("nacos server config cannot empty\n")
	}
	serverConfig := newServerConfigs(conf.Server)
	clientConfig := newClientConfig(conf.Client)

	naming, err := newNamingClient(serverConfig, clientConfig)
	if err != nil {
		return nil, err
	}
	config, err := newConfigClient(serverConfig, clientConfig)
	if err != nil {
		return nil, err
	}
	return &NacosClient{
		namingClient: naming,
		configClient: config,
	}, nil
}

func newClientConfig(conf config.NacosClientConfig) constant.ClientConfig {
	if conf == emptyClientConf {
		return defaultClient
	}
	return constant.ClientConfig{
		NamespaceId:    conf.NamespaceId,
		TimeoutMs:      conf.TimeoutMs,
		ListenInterval: conf.ListenInterval,
		BeatInterval:   conf.BeatInterval,
		LogDir:         conf.LogDir,
		CacheDir:       conf.CacheDir,
		Username:       conf.Username,
		Password:       conf.Password,
	}
}

func newServerConfigs(confs map[string]config.NacosServerConfig) []constant.ServerConfig {
	host := os.Getenv("REGISTER_HOST")
	if host == "" {
		log.Fatal("nacos host is empty")
		return nil
	}
	portStr := os.Getenv("REGISTER_PORT")
	if portStr == "" {
		log.Fatal("nacos port is empty")
		return nil
	}
	port, err := strconv.ParseUint(portStr, 10, 64)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	var ss []constant.ServerConfig
	for _, conf := range confs {
		ss = append(ss, constant.ServerConfig{
			IpAddr:      host,
			Port:        port,
			ContextPath: conf.ContextPath,
		})
	}
	return ss
}

func newNamingClient(ss []constant.ServerConfig, c constant.ClientConfig) (naming_client.INamingClient, error) {
	return clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": ss,
		"clientConfig":  c,
	})
}

func newConfigClient(ss []constant.ServerConfig, c constant.ClientConfig) (config_client.IConfigClient, error) {
	return clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": ss,
		"clientConfig":  c,
	})
}

func (n *NacosClient) RegisterInstance(param vo.RegisterInstanceParam) (bool, error) {
	return n.namingClient.RegisterInstance(param)
}

func (n *NacosClient) DeregisterInstance(param vo.DeregisterInstanceParam) (bool, error) {
	return n.namingClient.DeregisterInstance(param)
}

func (n *NacosClient) GetService(param vo.GetServiceParam) (model.Service, error) {
	return n.namingClient.GetService(param)
}

func (n *NacosClient) GetAllServicesInfo(param vo.GetAllServiceInfoParam) (model.ServiceList, error) {
	return n.namingClient.GetAllServicesInfo(param)
}

func (n *NacosClient) SelectAllInstances(param vo.SelectAllInstancesParam) ([]model.Instance, error) {
	return n.namingClient.SelectAllInstances(param)
}

func (n *NacosClient) SelectInstances(param vo.SelectInstancesParam) ([]model.Instance, error) {
	return n.namingClient.SelectInstances(param)
}

func (n *NacosClient) SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) (*model.Instance, error) {
	return n.namingClient.SelectOneHealthyInstance(param)
}

func (n *NacosClient) Subscribe(param *vo.SubscribeParam) error {
	return n.namingClient.Subscribe(param)
}

func (n *NacosClient) Unsubscribe(param *vo.SubscribeParam) error {
	return n.namingClient.Unsubscribe(param)
}

func (n *NacosClient) DoGet(serviceName, api string, params map[string]interface{}, headerOptions ...http.HeaderOption) ([]byte, int, error) {
	instance, err := n.namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
	})
	if err != nil {
		return nil, 501, err
	}
	if instance == nil {
		return nil, 501, errors.New(fmt.Sprintf("service [%s] not found one healthy instance! "))
	}
	url := fmt.Sprintf("http://%s:%d/%s", instance.Ip, instance.Port, api)
	return http.Get(url, params, headerOptions...)
}

func (n *NacosClient) DoPost(serviceName, api string, params map[string]interface{}, body []byte, headerOptions ...http.HeaderOption) ([]byte, int, error) {
	instance, err := n.namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
	})
	if err != nil {
		return nil, 500, err
	}
	if instance == nil {
		return nil, 500, errors.New(fmt.Sprintf("service [%s] not found one healthy instance! "))
	}
	url := fmt.Sprintf("http://%s:%d/%s", instance.Ip, instance.Port, api)
	return http.Post(url, params, body, headerOptions...)
}
