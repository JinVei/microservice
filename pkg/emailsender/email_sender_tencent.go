package emailsender

import (
	"context"

	"github.com/jinvei/microservice/base/framework/log"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
)

var flog = log.New()

type TEmailSender struct {
	client *ses.Client
	conf   Config
}

// tencent platform implement
func NewTEmailSender(conf Config) (EmailSender, error) {
	credential := common.NewCredential(
		conf.SID,
		conf.SK,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ses.tencentcloudapi.com"

	client, err := ses.NewClient(credential, "ap-hongkong", cpf)
	if err != nil {
		return nil, err
	}
	return &TEmailSender{
		conf:   conf,
		client: client,
	}, nil

}

func (v *TEmailSender) Send(ctx context.Context, email, parameter string) error {
	request := ses.NewSendEmailRequest()

	request.FromEmailAddress = common.StringPtr(v.conf.FromEmail)
	request.ReplyToAddresses = common.StringPtr(v.conf.ReplyTo)
	request.Destination = common.StringPtrs([]string{email})

	request.Template = &ses.Template{
		TemplateID:   common.Uint64Ptr(v.conf.TemplateID),
		TemplateData: common.StringPtr(parameter),
	}

	request.Subject = common.StringPtr(v.conf.Subject)
	request.TriggerType = common.Uint64Ptr(1)
	request.SetContext(ctx)

	resp, err := v.client.SendEmail(request)
	if err, ok := err.(*errors.TencentCloudSDKError); ok {
		flog.Errorf("An API error has returned: %s", err)
		return err
	}
	// TODO: check resp.Response.MessageId
	//resp.Response.MessageId
	flog.Debug(resp.ToJsonString())

	return nil
}
