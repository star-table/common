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
	FeiShu:        nil,
	ScheduleTime:  nil,
	Logs:          nil,
	Application:   nil,
	Parameters:    nil,
	Mq:            nil,
	OSS:           nil,
	ElasticSearch: nil,
	Sentry:        nil,
	SMS:           nil,
}

type Config struct {
	Viper         *viper.Viper
	Mysql         *MysqlConfig        //数据库配置
	Redis         *RedisConfig        //redis配置
	Mail          *MailConfig         //邮件配置
	Server        *ServerConfig       //服务配置
	DingTalk      *DingTalkSDKConfig  //钉钉配置
	FeiShu        *FeiShuSdkConfig    //飞书配置
	ScheduleTime  *ScheduleTimeConfig //定时时间配置
	Logs          *map[string]LogConfig
	Application   *ApplicationConfig   //应用配置
	Parameters    *ParameterConfig     //参数配置
	Mq            *MQConfig            //mq配置
	OSS           *OSSConfig           //oss配置
	ElasticSearch *ElasticSearchConfig //es配置
	Sentry        *SentryConfig        //sentry配置
	SkyWalking    *SkyWalkingConfig    //skywalking配置
	SMS           *SMSConfig           //消息配置
}

type ScheduleTimeConfig struct {
	ScheduleDailyProjectReportSecondInterval int    //期间区间秒字段的值,提供给time.ParseDuration函数使用
	ScheduleDailyProjectReportTriggerCron    string //配置项目日报每日触发时间

	ScheduleIssueRemindMaxCompensationTime int //任务提醒最大补偿时长，单位分钟（当服务重启时，会有一段时间范围导致丢失，这时要对其补偿，该字段定义最大补偿时长）
}

//mq配置
type MysqlConfig struct {
	Host     string
	Port     int
	Usr      string
	Pwd      string
	Database string
}

//redis配置
type RedisConfig struct {
	Host           string
	Port           int
	Pwd            string
	Database       int
	MaxIdle        int
	MaxActive      int
	MaxIdleTimeout int
}

//oss配置
type OSSConfig struct {
	BucketName      string
	EndPoint        string
	AccessKeyId     string
	AccessKeySecret string
	Policies        OSSPolicyConfig
	RootPath        string
	LocalDomain     string
}

type OSSPolicyConfig struct {
	ProjectCover   OSSPolicyInfo
	IssueResource  OSSPolicyInfo
	IssueInputFile OSSPolicyInfo
	ProjectResource OSSPolicyInfo
	CompatTest		OSSPolicyInfo
}

type OSSPolicyInfo struct {
	//有效期
	Expire      int64
	//目录
	Dir         string
	//最大文件大小
	MaxFileSize int64
	//回调地址
	CallbackUrl string
}

//邮件配置
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

type FeiShuSdkConfig struct {
	AppId            string
	AppSecret        string
	EventEncryptKey  string
	EventVerifyToken string
}

type LogConfig struct {
	LogPath          string
	Level            string
	FileSize         int64
	FileNum          int
	IsConsoleOut     bool
	Tag              string
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
	NameServers    string
	ReconsumeTimes int
	RePushTimes    int
}

type TopicConfig struct {
	//任务动态
	IssueTrends TopicConfigInfo
	//任务提醒
	IssueRemind TopicConfigInfo
	//项目动态
	ProjectTrends TopicConfigInfo
	//项目日报
	DailyProjectReportProject TopicConfigInfo
	//项目日报Msg
	DailyProjectReportMsg TopicConfigInfo
	//导入任务
	ImportIssue TopicConfigInfo
	//组织成员变动
	OrgMemberChange TopicConfigInfo
	//飞书回调消费
	FeiShuCallBack TopicConfigInfo
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

type SMSConfig struct {
	AccessKeyId     string
	AccessKeySecret string
	Region          string
}

func GetSMSConfig() *SMSConfig {
	return conf.SMS
}

func GetScheduleTime() *ScheduleTimeConfig {
	return conf.ScheduleTime
}

func GetMqDailyProjectReportProjectTopicConfig() TopicConfigInfo {
	return conf.Mq.Topics.DailyProjectReportProject
}

func GetMqDailyProjectReportMsgTopicConfig() TopicConfigInfo {
	return conf.Mq.Topics.DailyProjectReportMsg
}

func GetMqIssueTrendsTopicConfig() TopicConfigInfo {
	return conf.Mq.Topics.IssueTrends
}

func GetMqProjectTrendsTopicConfig() TopicConfigInfo {
	return conf.Mq.Topics.ProjectTrends
}

func GetMqFeiShuCallBackTopicConfig() TopicConfigInfo {
	return conf.Mq.Topics.FeiShuCallBack
}

func GetMqImportIssueTopicConfig() TopicConfigInfo {
	return conf.Mq.Topics.ImportIssue
}

func GetMqOrgMemberChangeConfig() TopicConfigInfo{
	return conf.Mq.Topics.OrgMemberChange
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

func GetProjectResourcePolicyConfig() OSSPolicyInfo {
	return conf.OSS.Policies.ProjectResource
}

func GetCompatTestPolicyConfig() OSSPolicyInfo{
	return conf.OSS.Policies.CompatTest
}

func GetSentryConfig() *SentryConfig {
	return conf.Sentry
}

func GetSkyWalkingConfig() *SkyWalkingConfig {
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
	//configPath := "/Users/tree/work/08_all_star/01_src/go/polaris-backend/config"
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
	if err != nil {
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
	if conf.Viper == nil {
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
