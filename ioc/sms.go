package ioc

import (
	"Vchat/internal/service/sms"
	"Vchat/internal/service/sms/localsms"
	"Vchat/internal/service/sms/tencent"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
)

func InitSmsService() sms.Service {
	return localsms.NewService()
}

func InitTencentSmsService() sms.Service {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		panic("找不到腾讯 SMS 的 secret id")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		panic("找不到腾讯 SMS 的 secret key")
	}
	c, err := tencentSMS.NewClient(
		common.NewCredential(secretId, secretKey),
		"ap-nanjing",
		profile.NewClientProfile(),
	)
	if err != nil {
		panic(err)
	}
	return tencent.NewService(c, "XXX", "XXX")
}
