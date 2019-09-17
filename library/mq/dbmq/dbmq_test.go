package dbmq

import (
	"fmt"
	"gitea.bjx.cloud/allstar/common/core/config"
	"gitea.bjx.cloud/allstar/common/core/errors"
	"gitea.bjx.cloud/allstar/common/core/model"
	"strconv"
	"testing"
)

func pushMsg(message *model.MqMessageExt) errors.SystemErrorInfo {

	if "msgeeee - 4" == message.Body {
		fmt.Printf("fail: %v \n", message)
		return errors.BuildSystemErrorInfo(errors.DbMQConsumerStartedError)
	}
	fmt.Println(message)
	return nil
}

func TestGetDbMQProxy(t *testing.T) {

	config.LoadConfig("F:\\workspace-golang-polaris\\polaris-backend\\polaris-server\\config", "application")

	for i := 0; i < 10 ; i++ {

		msg := &model.MqMessage{

			Topic: "test_test_topic",
			Tags:  "a,b,c",
			Keys:  "123",
			Body:  "msgeeee - " + strconv.FormatInt(int64(i), 10),
		}

		msgExt, err := GetDbMQProxy().PushMessage(msg)

		fmt.Println(msgExt)
		fmt.Println(err)
	}

	err := GetDbMQProxy().ConsumeMessage("test_test_topic", "consumer1", pushMsg)
	fmt.Println(err)

}
