package sms

import (
	"errors"
	"fmt"
	"gitea.bjx.cloud/allstar/common/core/config"
	"gitea.bjx.cloud/allstar/common/core/util/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

func SendSMS(phone string, signName string, templateCode string, params map[string]string) (*dysmsapi.SendSmsResponse, error){
	smsConfig := config.GetSMSConfig()
	if smsConfig == nil{
		return nil, errors.New("missing SMS configuration")
	}

	client, err := dysmsapi.NewClientWithAccessKey(smsConfig.Region, smsConfig.AccessKeyId, smsConfig.AccessKeySecret)

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = phone
	request.SignName = signName
	request.TemplateCode = templateCode
	request.TemplateParam = json.ToJsonIgnoreError(params)

	response, err := client.SendSms(request)
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	return response, nil
}


