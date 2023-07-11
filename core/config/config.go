package config

import (
	"flag"
	"fmt"
	nacosconfig "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"

	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
)

var conf = &Config{
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
	MQTT:          nil,
}

type Config struct {
	Viper         *viper.Viper
	Mysql         *MysqlConfig             // 数据库配置
	Redis         *RedisConfig             // redis配置
	Mail          *MailConfig              // 邮件配置
	Server        *ServerConfig            // 服务配置
	DingTalk      *DingTalkSDKConfig       // 钉钉配置
	FeiShu        *FeiShuSdkConfig         // 飞书配置
	WeiXin        *WeiXinSdkConfig         // 微信配置
	PersonWeiXin  *PersonWeiXinLoginConfig // 个人微信
	ServerCommon  *ServerCommonConfig      // 一些共有配置
	ScheduleTime  *ScheduleTimeConfig      // 定时时间配置
	Logs          *map[string]LogConfig    // log配置
	Application   *ApplicationConfig       // 应用配置
	Parameters    *ParameterConfig         // 参数配置
	Mq            *MQConfig                // mq配置
	OSS           *OSSConfig               // oss配置
	ElasticSearch *ElasticSearchConfig     // es配置
	Sentry        *SentryConfig            // sentry配置
	SkyWalking    *SkyWalkingConfig        // skywalking配置
	SMS           *SMSConfig               // 消息配置
	MQTT          *MQTTConfig              // mqtt配置
	Jaeger        *JaegerConfig            // Jaeger
	WeLink        *WeLinkConfig            // welink配置
	Wechat        *WechatConfig            // wechat
	Nacos         *NacosConfig             // nacos
	PA            *PrivatizationAuthority  // 私有化授权
	YiDun         *YiDunConfig             // 易盾
}

type ScheduleTimeConfig struct {
	ScheduleDailyProjectReportSecondInterval int    //期间区间秒字段的值,提供给time.ParseDuration函数使用
	ScheduleDailyProjectReportTriggerCron    string //配置项目日报每日触发时间

	ScheduleIssueRemindMaxCompensationTime int //任务提醒最大补偿时长，单位分钟（当服务重启时，会有一段时间范围导致丢失，这时要对其补偿，该字段定义最大补偿时长）
}

//mq配置
type MysqlConfig struct {
	Host         string
	Port         int
	Usr          string
	Pwd          string
	Database     string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  int
	Driver       string
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
	IsSentinel     bool
	MasterName     string
	SentinelPwd    string
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
	ProjectCover    OSSPolicyInfo
	IssueResource   OSSPolicyInfo
	IssueInputFile  OSSPolicyInfo
	ProjectResource OSSPolicyInfo
	CompatTest      OSSPolicyInfo
	UserAvatar      OSSPolicyInfo
	Feedback        OSSPolicyInfo
}

type OSSPolicyInfo struct {
	//容器
	BucketName string
	//有效期
	Expire int64
	//目录
	Dir string
	//最大文件大小
	MaxFileSize int64
	//回调地址
	CallbackUrl string
}

//邮件配置
type MailConfig struct {
	//别名
	Alias string
	Usr   string
	Pwd   string
	Host  string
	Port  int
}

type ServerConfig struct {
	Port              int
	Name              string
	Host              string
	Domain            string
	NewRedirectDomain string
	Api               string
}

type SentryConfig struct {
	Dsn string
}

type SkyWalkingConfig struct {
	ReportAddress string
}

type JaegerConfig struct {
	UdpAddress   string
	TraceService string
	SamplerType  string
	SamplerParam float64
}

type DingTalkSDKConfig struct {
	SuiteKey       string
	SuiteSecret    string
	Token          string
	AesKey         string
	AppId          int64
	OauthAppId     string
	OauthAppSecret string
	RunType        int
	AgentId        int64
	TemplateId     string // 使用消息模板发送个人消息的模板id
	CallBackUrl    string
	FrontUrl       string
	CoolAppCode    string
}

type FeiShuSdkConfig struct {
	AppId            string
	AppSecret        string
	EventEncryptKey  string
	EventVerifyToken string
	CardJumpLink     CardJumpLink
}

type WeiXinSdkConfig struct {
	CorpId         string
	ProviderSecret string
	SuiteId        string
	Secret         string
	SuiteToken     string
	SuiteAesKey    string
	CorpDSN        string
	BuyerUserId    string
	CustomInfos    []struct {
		CorpId  string
		AgentId int
		Secret  string
	}
}

type PersonWeiXinLoginConfig struct {
	AppId     string
	AppSecret string
	Token     string
	AesKey    string
}

type ServerCommonConfig struct {
	Host                        string
	WeiXinHost                  string
	FsNewbieGuideTemplateId     int64
	DingNewbieGuideTemplateId   int64
	WeiXinNewbieGuideTemplateId int64
}

//飞书卡片跳转链接
type CardJumpLink struct {
	//项目日报pc跳转url
	ProjectDailyPcUrl string
	//个人日报pc跳转url
	PersonalDailyPcUrl string
	//欢迎页面进入PC项目列表页
	ProjectWelcomeUrl string
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
	// 0， 公共， 1：私有
	RunType int
	// public:集成部署；private:私有化部署；
	DeployType string
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

// 私有化授权配置
type PrivatizationAuthority struct {
	// 数据上报入口
	EndPoint string
	// 密钥
	Secret string
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
	// 个人日报Msg
	DailyIssueReportMsg TopicConfigInfo
	//导入任务
	ImportIssue TopicConfigInfo
	//组织成员变动
	OrgMemberChange TopicConfigInfo
	//同步成员、部门信息
	SyncMemberDept TopicConfigInfo
	// 提醒新用户使用北极星
	RemindUsingToNewUser TopicConfigInfo
	//飞书回调消费
	FeiShuCallBack TopicConfigInfo
	//飞书任务日报推送
	DailyTaskPushToFeishu TopicConfigInfo
	//迭代燃尽图
	StatisticIterationBurnDownChart TopicConfigInfo
	//项目燃尽图
	StatisticProjectIssueBurnDownChart TopicConfigInfo
	//飞书欢迎语推送
	HelpMessagePushToFeiShu TopicConfigInfo
	// 飞书提醒新版本进入方式
	FeiShuEntranceRemind TopicConfigInfo
	//第三方订单处理
	ThirdSourceChannelOrder TopicConfigInfo
	//飞书回调消息处理
	FeishuCallBackMsg TopicConfigInfo
	//飞书回调订单消息处理
	FeishuCallBackOrder TopicConfigInfo
	//飞书捷径triger处理
	FeishuShortcut TopicConfigInfo
	//飞书捷径时间推送
	FeishuShortcutPush TopicConfigInfo
	//飞书统一消息推送
	FeishuCommonMsgPush TopicConfigInfo
	// 飞书后台动态增加推送卡片任务
	FeishuDynCommonMsgPush TopicConfigInfo
	// bjx 后台动态增加脚本任务
	PolarisDynScriptTask TopicConfigInfo
	//向飞书发送北极星的调查问卷
	FeishuSurveyNotice1123Push TopicConfigInfo
	//飞书回调通讯录范围变更
	FeishuCallBackContactScopeChange TopicConfigInfo
	//飞书回调员工变更
	FeishuCallBackUserChange TopicConfigInfo
	//钉钉回调消息推送
	DingTalkCallBackMsg TopicConfigInfo
	//日志上报-事件
	LogEvent TopicConfigInfo
	//发送消息卡片
	CardPush TopicConfigInfo
	//批量创建任务
	BatchCreateIssue TopicConfigInfo
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
	Resource          *ResourceConfig
}

type ResourceConfig struct {
	Office *OfficeConfig
}

type OfficeConfig struct {
	OfficeType string // collabora microsoft
	Collabora  *CollaboraOffice
	Microsoft  *MicrosoftOffice
}

type CollaboraOffice struct {
	Url  string
	Wopi string
	Ext  string
}
type MicrosoftOffice struct {
	Url  string
	Wopi string
	Ext  string
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

type MQTTConfig struct {
	//地址
	Address string
	//Host
	Host string
	//Port
	Port int
	//key
	SecretKey string
	//channel
	Channel string
	//连接数
	ConnectPoolSize int
	//是否开启
	Enable bool
}

type WeLinkConfig struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type WechatConfig struct {
	ProviderCorpId string `json:"providerCorpId"`
	SuiteId        string `json:"suiteId"`
	Secret         string `json:"secret"`
	Token          string `json:"token"`
	AesKey         string `json:"aesKey"`
}

type NacosBaseConfig struct {
	AppName   string
	Host      string
	Port      string
	NameSpace string
	UserName  string
	Password  string
	Group     string
}

func GetYiDunConfig() *YiDunConfig {
	return conf.YiDun
}

func GetWeLinkConfig() *WeLinkConfig {
	return conf.WeLink
}

func GetWechatConfig() *WechatConfig {
	return conf.Wechat
}

func GetSMSConfig() *SMSConfig {
	return conf.SMS
}

func GetScheduleTime() *ScheduleTimeConfig {
	return conf.ScheduleTime
}

func GetMQTTConfig() *MQTTConfig {
	return conf.MQTT
}

func GetDailyTaskPushToFeishu() TopicConfigInfo {
	return conf.Mq.Topics.DailyTaskPushToFeishu
}

func GetStatisticIterationBurnDownChart() TopicConfigInfo {
	return conf.Mq.Topics.StatisticIterationBurnDownChart
}

func GetStatisticProjectIssueBurnDownChart() TopicConfigInfo {
	return conf.Mq.Topics.StatisticProjectIssueBurnDownChart
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

// 同步飞书用户、部门信息 topic
func GetMqSyncMemberDeptTopicConfig() TopicConfigInfo {
	return conf.Mq.Topics.SyncMemberDept
}

// 提醒新用户使用北极星的推送配置
func GetMqRemindUsingToNewUserPushTopicConfig() TopicConfigInfo {
	return conf.Mq.Topics.RemindUsingToNewUser
}

func GetMqOrgMemberChangeConfig() TopicConfigInfo {
	return conf.Mq.Topics.OrgMemberChange
}

func GetHelpMessagePushToFeiShu() TopicConfigInfo {
	return conf.Mq.Topics.HelpMessagePushToFeiShu
}

func GetFeiShuEntranceRemind() TopicConfigInfo {
	return conf.Mq.Topics.FeiShuEntranceRemind
}

func GetThirdSourceChannelOrder() TopicConfigInfo {
	return conf.Mq.Topics.ThirdSourceChannelOrder
}

func GetFeishuCallBackMsg() TopicConfigInfo {
	return conf.Mq.Topics.FeishuCallBackMsg
}

func GetFeishuCallBackOrder() TopicConfigInfo {
	return conf.Mq.Topics.FeishuCallBackOrder
}

func GetFeishuShortcutTopic() TopicConfigInfo {
	return conf.Mq.Topics.FeishuShortcut
}

func GetFeishuShortcutPush() TopicConfigInfo {
	return conf.Mq.Topics.FeishuShortcutPush
}

func GetFeishuCommonMsgPushTopic() TopicConfigInfo {
	return conf.Mq.Topics.FeishuCommonMsgPush
}

func GetFeishuDynCommonMsgPushTopic() TopicConfigInfo {
	return conf.Mq.Topics.FeishuDynCommonMsgPush
}

func GetPolarisDynScriptTaskTopic() TopicConfigInfo {
	return conf.Mq.Topics.PolarisDynScriptTask
}

func GetFeishuSurveyNotice1123PushTopic() TopicConfigInfo {
	return conf.Mq.Topics.FeishuSurveyNotice1123Push
}

func GetFeishuCallBackContactScopeChange() TopicConfigInfo {
	return conf.Mq.Topics.FeishuCallBackContactScopeChange
}

func GetFeishuCallBackUserChange() TopicConfigInfo {
	return conf.Mq.Topics.FeishuCallBackUserChange
}

func GetDingTalkCallBackTopic() TopicConfigInfo {
	return conf.Mq.Topics.DingTalkCallBackMsg
}

func GetCardPushConfig() TopicConfigInfo {
	return conf.Mq.Topics.CardPush
}

func GetBatchCreateIssueTopicConfig() TopicConfigInfo {
	return conf.Mq.Topics.BatchCreateIssue
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

func GetCompatTestPolicyConfig() OSSPolicyInfo {
	return conf.OSS.Policies.CompatTest
}

func GetUserAvatarPolicyConfig() OSSPolicyInfo {
	return conf.OSS.Policies.UserAvatar
}

func GetFeedbackPolicyConfig() OSSPolicyInfo {
	return conf.OSS.Policies.Feedback
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

type NacosConfig struct {
	Client    NacosClientConfig            `json:"client"`
	Server    map[string]NacosServerConfig `json:"server"`
	Discovery *DiscoveryConfig             `json:"discovery"`
}

// 配置
type YiDunConfig struct {
	// API
	Api string `json:"api"`
	// 版本
	Version string `json:"version"`
	// CaptchaId
	CaptchaId string `json:"captchaId"`
	// SecretId
	SecretId string `json:"secretId"`
	// SecretKey
	SecretKey string `json:"secretKey"`
}

type NacosClientConfig struct {
	TimeoutMs            uint64 `json:"timeout_ms"`              // TimeoutMs http请求超时时间，单位毫秒
	ListenInterval       uint64 `json:"listen_interval"`         // ListenInterval 监听间隔时间，单位毫秒（仅在ConfigClient中有效）
	BeatInterval         int64  `json:"beat_interval"`           // BeatInterval         心跳间隔时间，单位毫秒（仅在ServiceClient中有效）
	NamespaceId          string `json:"namespace_id"`            // NamespaceId          nacos命名空间
	Endpoint             string `json:"endpoint"`                // Endpoint             获取nacos节点ip的服务地址
	CacheDir             string `json:"cacheDir"`                // CacheDir             缓存目录
	LogDir               string `json:"logDir"`                  // LogDir               日志目录
	UpdateThreadNum      int    `json:"update_thread_num"`       // UpdateThreadNum      更新服务的线程数
	NotLoadCacheAtStart  bool   `json:"not_load_cache_at_start"` // NotLoadCacheAtStart  在启动时不读取本地缓存数据，true--不读取，false--读取
	UpdateCacheWhenEmpty bool   `json:"update_cache_when_empty"` // UpdateCacheWhenEmpty 当服务列表为空时是否更新本地缓存，true--更新,false--不更新
	Username             string `json:"username"`
	Password             string `json:"password"`
}

type NacosServerConfig struct {
	IpAddr      string `json:"ip_addr"`      // IpAddr      nacos命名空间
	ContextPath string `json:"context_path"` // ContextPath 获取nacos节点ip的服务地址
	Port        uint64 `json:"port"`         // Port        缓存目录
}

// DiscoveryConfig is nacos service config.
type DiscoveryConfig struct {
	GroupName   string            `json:"group_name"`
	ClusterName string            `json:"cluster_name"`
	Weight      float64           `json:"weight"`
	Enable      bool              `json:"enable"`
	Healthy     bool              `json:"healthy"`
	Ephemeral   bool              `json:"ephemeral"`
	MetaData    map[string]string `json:"metadata"`
}

func GetEnv() string {
	env := os.Getenv("POL_ENV")
	if "" == env {
		env = "local"
	}
	return env
}

func GetJaegerConfig() *JaegerConfig {
	return conf.Jaeger
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

func GetConfig() *Config {
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

/**
 * 获取office配置
 **/
func GetOffice() *OfficeConfig {

	if conf.Parameters.Resource == nil {
		return nil
	}
	return conf.Parameters.Resource.Office
}

/**
 * 获取officeUrl
 **/
func GetOfficeUrl() string {

	var officeUrl string = ""

	if conf.Parameters.Resource == nil {
		return officeUrl
	}

	officeType := conf.Parameters.Resource.Office.OfficeType
	if officeType == "collabora" {
		officeUrl = conf.Parameters.Resource.Office.Collabora.Url + conf.Parameters.Resource.Office.Collabora.Wopi
	} else if officeType == "microsoft" {
		officeUrl = conf.Parameters.Resource.Office.Microsoft.Url + conf.Parameters.Resource.Office.Microsoft.Wopi
	}
	fmt.Printf("officeUrl:%s", officeUrl)
	return officeUrl
}

func LoadConfig(flagconf, nacosHost, nacosPort, nacosNamespace, appName string) error {
	return loadConfig(flagconf, nacosHost, nacosPort, nacosNamespace, appName)
}

func loadConfig(flagconf, nacosHost, nacosPort, nacosNamespace, appName string) error {
	flag.Parse()
	if flagconf != "" {
		c := config.New(
			config.WithSource(
				file.NewSource(flagconf),
			),
		)
		defer c.Close()

		if err := c.Load(); err != nil {
			panic(err)
		}

		if err := c.Scan(&conf); err != nil {
			panic(err)
		}

		return nil
	}

	return loadNacosConfig(nacosHost, nacosPort, nacosNamespace, appName)
}

func loadNacosConfig(nacosHost, nacosPort, nacosNamespace, appName string) error {
	client, err := getNacosConfigClient(nacosHost, nacosPort, nacosNamespace, appName)
	if err != nil {
		return err
	}

	configSource := nacosconfig.NewConfigSource(client, nacosconfig.WithGroup("DEFAULT_GROUP"), nacosconfig.WithDataID(appName))

	c := config.New(
		config.WithSource(
			configSource,
		),
		config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)

	if err := c.Load(); err != nil {
		return err
	}

	if err := c.Scan(&conf); err != nil {
		return err
	}

	return nil
}

func getNacosServerAndClientConfig(nacosHost, nacosPort, nacosNamespace, appName string) ([]constant.ServerConfig, constant.ClientConfig) {
	return []constant.ServerConfig{
			*constant.NewServerConfig(nacosHost, cast.ToUint64(nacosPort)),
		},
		constant.ClientConfig{
			AppName:             appName,
			NamespaceId:         nacosNamespace, //namespace id
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogLevel:            "error",
		}
}

func getNacosConfigClient(nacosHost, nacosPort, nacosNamespace, appName string) (config_client.IConfigClient, error) {
	sc, cc := getNacosServerAndClientConfig(nacosHost, nacosPort, nacosNamespace, appName)
	// a more graceful way to create naming client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		return nil, err
	}
	return client, nil
}
