package config

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var conf Config = Config{
	Viper:         viper.New(),
	Mysql:         nil,
	Redis:         nil,
	Mail:          nil,
	Server:        nil,
	DingTalk:      nil,
	Logs:          nil,
	Application:   nil,
	Parameters:    nil,
	Mq:            nil,
	OSS:           nil,
	ElasticSearch: nil,
	Sentry:		   nil,
}

type Config struct {
	Viper         *viper.Viper
	Mysql         *MysqlConfig
	Redis         *RedisConfig
	Mail          *MailConfig
	Server        *ServerConfig
	DingTalk      *DingTalkSDKConfig
	Logs          *map[string]LogConfig
	Application   *ApplicationConfig
	Parameters    *ParameterConfig
	Mq            *MQConfig
	OSS           *OSSConfig
	ElasticSearch *ElasticSearchConfig
	Sentry		  *SentryConfig
	SkyWalking 	  *SkyWalkingConfig
}

type MysqlConfig struct {
	Host     string
	Port     int
	Usr      string
	Pwd      string
	Database string
}

type RedisConfig struct {
	Host           string
	Port           int
	Pwd            string
	Database       int
	MaxIdle        int
	MaxActive      int
	MaxIdleTimeout int
}

type OSSConfig struct {
	BucketName      string
	EndPoint        string
	AccessKeyId     string
	AccessKeySecret string
	Policies        OSSPolicyConfig
}

type OSSPolicyConfig struct {
	ProjectCover  OSSPolicyInfo
	IssueResource OSSPolicyInfo
	IssueInputFile OSSPolicyInfo
}

type OSSPolicyInfo struct {
	Expire      int64
	Dir         string
	MaxFileSize int64
}

type MailConfig struct {
	Usr  string
	Pwd  string
	Host string
	Port int
}

type ServerConfig struct {
	Port int
	Name string
	Host string
}

type SentryConfig struct {
	Dsn string
}

type SkyWalkingConfig struct {
	ReportAddress string
}

type DingTalkSDKConfig struct {
	SuiteKey    string
	SuiteSecret string
	Token       string
	AesKey      string
	AppId       int64
}

type LogConfig struct {
	LogPath          string
	Level            string
	FileSize         int64
	FileNum          int
	IsConsoleOut     bool
	EnableKafka      bool
	KafkaTopic       string
	KafkaNameServers string
}

type ApplicationConfig struct {
	RunMode   int
	CacheMode string
	Name      string
	AppCode   string
	AppKey    string
}

type MQConfig struct {
	Mode   string
	Rocket *RocketMQConfig
	Kafka  *KafkaMQConfig
	Topics TopicConfig
}

type RocketMQConfig struct {
	GroupID    string
	NameServer string
	Log        *LogConfig
}

//KafKa MQ Config
type KafkaMQConfig struct {
	NameServers string
	ReconsumeTimes int
	RePushTimes int
}

type TopicConfig struct {
	IssueTrends   TopicConfigInfo
	ProjectTrends TopicConfigInfo
}

type TopicConfigInfo struct {
	Topic   string
	GroupId string
}

//ElasticSearch 配置
type ElasticSearchConfig struct {
	ServerUrls []string
	Sniff      bool
	Timeout    int64
	Auth       *ElasticSearchAuthConfig
}
type ElasticSearchAuthConfig struct {
	UserName string
	Password string
}
type ParameterConfig struct {
	CodeConnector     string
	IdBufferThreshold float64
	MaxPageSize       int
	EsIndex           *EsIndexConfig
	PreUrl            map[string]string
}

type EsIndexConfig struct {
	Issue   string
	Project string
	Trends  string
}

func GetMqIssueTrendsTopicConfig() TopicConfigInfo {
	return conf.Mq.Topics.IssueTrends
}

func GetMqProjectTrendsTopicConfig() TopicConfigInfo {
	return conf.Mq.Topics.ProjectTrends
}

func GetProjectCoverPolicyConfig() OSSPolicyInfo {
	return conf.OSS.Policies.ProjectCover
}

func GetIssueResourcePolicyConfig() OSSPolicyInfo {
	return conf.OSS.Policies.IssueResource
}

func GetIssueInputFilePolicyConfig() OSSPolicyInfo {
	return conf.OSS.Policies.IssueInputFile
}

func GetSentryConfig() *SentryConfig{
	return conf.Sentry
}

func GetSkyWalkingConfig() *SkyWalkingConfig{
	return conf.SkyWalking
}

func GetMysqlConfig() *MysqlConfig {
	return conf.Mysql
}

func GetEnv() string {
	env := os.Getenv("POL_ENV")
	if "" == env {
		env = "local"
	}
	return env
}

func GetKafkaConfig() *KafkaMQConfig {
	if conf.Mq == nil {
		panic(errors.New("mq configuration is nil!"))
	}
	if conf.Mq.Kafka == nil {
		panic(errors.New("kafka configuration is nil!"))
	}
	return conf.Mq.Kafka
}

func GetOSSConfig() *OSSConfig {
	return conf.OSS
}

func GetRedisConfig() *RedisConfig {
	return conf.Redis
}

func GetConfig() Config {
	return conf
}

func GetMailConfig() *MailConfig {
	return conf.Mail
}

func GetServerConfig() *ServerConfig {
	return conf.Server
}

func GetDingTalkSdkConfig() *DingTalkSDKConfig {
	return conf.DingTalk
}

func GetApplication() *ApplicationConfig {
	return conf.Application
}

func GetLogConfig(name string) *LogConfig {
	c := (*conf.Logs)[name]
	return &c
}

func GetMQ() *MQConfig {
	return conf.Mq
}

func GetElasticSearch() *ElasticSearchConfig {
	return conf.ElasticSearch
}

func GetParameters() *ParameterConfig {
	return conf.Parameters
}

func GetPreUrl(name string) string {
	if v, ok := conf.Parameters.PreUrl[name]; ok {
		return v
	}
	return ""
}

func LoadConfig(dir string, config string) error {
	return LoadEnvConfig(dir, config, "")
}

func LoadExtraConfig(dir string, config string, extraConfig interface{}) error {
	return LoadExtraEnvConfig(dir, config, "", extraConfig)
}

func LoadUnitTestConfig() {
	configPath := "F:\\workspace-golang-polaris\\polaris-backend\\config"
	configName := "application.common"
	env := "local"
	for _, arg := range flag.Args() {
		argList := strings.Split(arg, "=")
		if len(argList) != 2 {
			argList = strings.Split(arg, " ")
		}
		if len(argList) != 2 {
			fmt.Printf(" unknown arg:%v\n", arg)
			continue
		}
		arg0 := strings.TrimSpace(argList[0])
		if arg0 == "p" || arg0 == "P" {
			configPath = argList[1]
		}
		if arg0 == "n" || arg0 == "N" {
			configName = argList[1]
		}
		if arg0 == "e" || arg0 == "E" {
			env = argList[1]
		}
	}
	LoadEnvConfig(configPath, configName, env)
}

func LoadEnvConfig(dir string, config string, env string) error {
	err := loadConfig(dir, config, "")
	if err != nil {
		return err
	}
	if env != "" {
		err = loadConfig(dir, config, env)
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadExtraEnvConfig(dir string, config string, env string, extraConfig interface{}) error {
	err := loadExtraConfig(dir, config, "", extraConfig)
	if err != nil {
		return err
	}
	if env != "" {
		err = loadExtraConfig(dir, config, env, extraConfig)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadExtraConfig(dir string, config string, env string, extraConfig interface{}) error {
	err := loadConfig(dir, config, env)
	if err != nil{
		return err
	}
	if err := conf.Viper.Unmarshal(&extraConfig); err != nil {
		return err
	}
	return nil
}


func loadConfig(dir string, config string, env string) error {
	configName := config
	if env != "" {
		configName += "." + env
	}
	if conf.Viper == nil{
		conf.Viper = viper.New()
	}
	conf.Viper.SetConfigName(configName)
	conf.Viper.AddConfigPath(dir)
	conf.Viper.SetConfigType("yaml")
	if err := conf.Viper.MergeInConfig(); err != nil {
		fmt.Println(err)
		return err
	}
	if err := conf.Viper.Unmarshal(&conf); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
