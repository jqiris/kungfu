package cloud

import (
	"fmt"
	"testing"

	"github.com/jqiris/kungfu/v2/config"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func TestSmsSend(t *testing.T) {
	client := NewSmsClient(config.TecentSms{
		SecretId:  "AKIDb1C7l9cz83jALGvKTMyOZS5SZS6UcT5R",
		SecretKey: "80cy4tSloc9QSFm23pKZYHjvrYz97oIY",
		EndPoint:  "sms.tencentcloudapi.com",
		Region:    "ap-nanjing",
		SdkAppid:  "1400739632",
		SignName:  "寻光长青",
	})
	response, err := client.SendMsg("18516536416", "1552628", []string{"123025", "5"})
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}
	// 输出json格式的字符串回包
	fmt.Printf("%s", response.ToJsonString())
}
