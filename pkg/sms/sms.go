package sms

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711"
)

type SendStatus struct {

	// 发送流水号。
	SerialNo *string `json:"SerialNo,omitempty" name:"SerialNo"`

	// 手机号码,e.164标准，+[国家或地区码][手机号] ，示例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号。
	PhoneNumber *string `json:"PhoneNumber,omitempty" name:"PhoneNumber"`

	// 计费条数，计费规则请查询 [计费策略](https://cloud.tencent.com/document/product/382/36135)。
	Fee *uint64 `json:"Fee,omitempty" name:"Fee"`

	// 用户Session内容。
	SessionContext *string `json:"SessionContext,omitempty" name:"SessionContext"`

	// 短信请求错误码，具体含义请参考错误码。
	Code *string `json:"Code,omitempty" name:"Code"`

	// 短信请求错误码描述。
	Message *string `json:"Message,omitempty" name:"Message"`

	// 国家码或地区码，例如CN,US等，对于未识别出国家码或者地区码，默认返回DEF,具体支持列表请参考国际/港澳台计费总览。
	IsoCode *string `json:"IsoCode,omitempty" name:"IsoCode"`
}
type SendSmsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 短信发送状态。
		SendStatusSet []*SendStatus `json:"SendStatusSet,omitempty" name:"SendStatusSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

var SMS *common.Credential

func CodeRandom() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	return code
}

func NewSNS(SecretId, SecretKey string) {
	SMS = common.NewCredential(
		SecretId,
		SecretKey,
	)
}

/*
sms发送短信服务

phone-手机号码

template_id-模板 ID

sms_sdk_id-短信SdkAppid

sign-短信签名内容

template_param-模板参数
*/

func SmsSend(phone, template_id, sms_sdk_id, sign string, template_param []string) string {

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, _ := sms.NewClient(SMS, "", cpf)

	request := sms.NewSendSmsRequest()

	request.PhoneNumberSet = common.StringPtrs([]string{phone})
	request.TemplateID = common.StringPtr(template_id)
	request.SmsSdkAppid = common.StringPtr(sms_sdk_id)
	request.Sign = common.StringPtr(sign)
	request.TemplateParamSet = common.StringPtrs(template_param)

	response, err := client.SendSms(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Sprintf("An API error has returned: %s", err)
	}
	if err != nil {
		return "error"
	}
	return *response.Response.SendStatusSet[0].Message
}
