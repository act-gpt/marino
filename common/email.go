package common

import (
	"fmt"
	"strings"

	"github.com/act-gpt/marino/config/system"

	"github.com/resendlabs/resend-go"
)

func SendEmailByResend(subject string, receiver string, content string) error {
	conf := system.Config.Mail
	// 未配置 mail 服务器
	if conf.SMTPToken == "" {
		return nil
	}
	client := resend.NewClient(conf.SMTPToken)
	to := strings.Split(receiver, ";")
	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", system.Config.Organization.Name, conf.SMTPFrom),
		To:      to,
		Html:    content,
		Subject: subject,
	}
	_, err := client.Emails.Send(params)
	return err
}
