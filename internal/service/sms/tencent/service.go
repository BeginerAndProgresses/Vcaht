package tencent

import (
	"Vchat/pkg/utils"
	"context"
	"fmt"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	client   *sms.Client
	appId    *string
	signName *string
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, phone ...string) error {
	request := sms.NewSendSmsRequest()
	request.SetContext(ctx)
	request.SmsSdkAppId = s.appId
	request.SignName = s.signName
	request.TemplateId = utils.ToPtr[string](tplId)
	request.TemplateParamSet = utils.ToPtrSlice[string](args)
	request.PhoneNumberSet = utils.ToPtrSlice[string](phone)
	response, err := s.client.SendSms(request)
	// 处理异常
	if err != nil {
		return err
	}
	for _, statusPtr := range response.Response.SendStatusSet {
		if statusPtr == nil {
			continue
		}
		status := *statusPtr
		if status.Code != nil && *status.Code != "Ok" {
			return fmt.Errorf("send sms failed, code: %s, message: %s", *status.Code, *status.Message)
		}
	}
	return nil
}

func NewService(client *sms.Client, appId, secretId string) *Service {
	return &Service{
		client:   client,
		appId:    &appId,
		signName: &secretId,
	}
}
