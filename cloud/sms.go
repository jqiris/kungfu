package cloud

import (
	"github.com/jqiris/kungfu/v2/config"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type SmsClient struct {
	client   *sms.Client
	SdkAppid string
	SignName string
}

func NewSmsClient(cfg config.TecentSms) *SmsClient {
	credential := common.NewCredential(
		cfg.SecretId,
		cfg.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = cfg.EndPoint
	client, err := sms.NewClient(credential, cfg.Region, cpf)
	if err != nil {
		logger.Error(err)
	}
	return &SmsClient{
		client:   client,
		SdkAppid: cfg.SdkAppid,
		SignName: cfg.SignName,
	}
}

func (s *SmsClient) SendMsg(phone, tplId string, args []string) (*sms.SendSmsResponse, error) {
	request := sms.NewSendSmsRequest()
	request.PhoneNumberSet = common.StringPtrs([]string{phone})
	request.SmsSdkAppId = common.StringPtr(s.SdkAppid)
	request.SignName = common.StringPtr(s.SignName)
	request.TemplateId = common.StringPtr(tplId)
	request.TemplateParamSet = common.StringPtrs(args)
	return s.client.SendSms(request)
}
