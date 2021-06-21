package email

import "strings"

/*
  @Author : Mustang Kong
*/

// SendEmail 发送纯文本email
func SendEmail(service string, textContent PlainTextContent) {
	emailClient := NewMailClient(strings.Split(service, "@")[0])
	emailClient.SendEmail(MailKindPlainText, textContent)
}
