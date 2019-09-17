package dingtalk

import (
	"gitea.bjx.cloud/allstar/common/core/config"
	"gitea.bjx.cloud/allstar/common/core/consts"
	"gitea.bjx.cloud/allstar/common/library/cache"
	"github.com/pkg/errors"
	"github.com/polaris-team/dingtalk-sdk-golang/sdk"
	"os"
	"strconv"
	"sync"
)

var sdkMutex sync.Mutex

var dingTalkSDK *sdk.DingTalkSDK

//var log = logger.GetDefaultLogger()

func initSDK() {
	if config.GetDingTalkSdkConfig() == nil {
		panic(errors.New("DingTalk SDK Configuration is missing!"))
	}
	if dingTalkSDK == nil {
		sdkMutex.Lock()
		defer sdkMutex.Unlock()
		if dingTalkSDK == nil {
			dtc := config.GetDingTalkSdkConfig()
			os.Setenv("SUITE_KEY", dtc.SuiteKey)
			os.Setenv("SUITE_SECRET", dtc.SuiteSecret)
			os.Setenv("SUITE_TOKEN", dtc.Token)
			os.Setenv("SUITE_AES_KEY", dtc.AesKey)
			os.Setenv("APP_ID", strconv.FormatInt(dtc.AppId, 10))
			dingTalkSDK = sdk.NewSDK()
		}
	}
}

func GetSDKProxy() *sdk.DingTalkSDK {
	initSDK()
	return dingTalkSDK
}

func GetCrypto() *sdk.Crypto {
	initSDK()
	return dingTalkSDK.CreateCrypto()
}

func GetDingTalkClientRest(corpId string) (*sdk.DingTalkClient, error) {
	suiteTicket, err := GetSuiteTicket()
	if err != nil {
		return nil, err
	}
	return GetDingTalkClient(corpId, suiteTicket)
}

func GetDingTalkClient(corpId string, suiteTicket string) (*sdk.DingTalkClient, error) {
	initSDK()
	corp := dingTalkSDK.CreateCorp(corpId, suiteTicket)

	tokenInfo, err := corp.GetCorpToken()
	if err != nil {
		return nil, err
	}
	if tokenInfo.ErrCode != 0 {
		return nil, errors.New(tokenInfo.ErrMsg)
	}
	authInfo, err := corp.GetAuthInfo()
	if err != nil {
		return nil, err
	}
	if authInfo.ErrCode != 0 {
		return nil, errors.New(authInfo.ErrMsg)
	}

	agents := authInfo.AuthInfo.Agent
	var targetAgent *sdk.Agent = nil
	for _, agent := range agents {
		if agent.AppId == config.GetDingTalkSdkConfig().AppId {
			targetAgent = &agent
			break
		}
	}
	if targetAgent == nil {
		return nil, errors.New("当前应用不在该企业授权应用列表中!")
	}
	client := NewDingTalkClient(tokenInfo.AccessToken, targetAgent.AgentId)
	return client, nil
}

func NewDingTalkClient(accessToken string, agentId int64) *sdk.DingTalkClient {
	return &sdk.DingTalkClient{
		AccessToken: accessToken,
		AgentId:     agentId,
	}
}

func GetSuiteTicket() (string, error) {
	return cache.Get(consts.CacheSuiteTicketKey)
}
