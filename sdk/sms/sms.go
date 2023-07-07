package sms

import (
	"errors"
	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/strs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

var log = logger.GetDefaultLogger()

func SendSMS(phone string, signName string, templateCode string, params map[string]string) (*dysmsapi.SendSmsResponse, error) {
	smsConfig := config.GetSMSConfig()
	if smsConfig == nil {
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
		log.Error("send sms error:" + strs.ObjectToString(err))
		return nil, err
	}
	return response, nil
}
