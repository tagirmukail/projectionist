package apps

import (
	"fmt"
	"net/smtp"
	"strconv"

	"projectionist/config"
)

const (
	SubjectTemplate = "Service %s notification"
	BodyHtmlPtrn    = `

<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>%s</title>
</head>
<body>
    <h1>%s</h1>
    Service %s notification: %s.
</body>
</html>
`
	MsgTmpl = "To: %s\r\n" +
		"Subject: %s\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n" +
		"%s\r\n"
)

type Notifier struct {
	cfg  *config.NotifierConfig
	auth smtp.Auth
	Body string
}

func NewNotifier(cfg *config.Config) *Notifier {
	auth := smtp.PlainAuth("", cfg.NotifierConfig.From, cfg.NotifierConfig.Password, cfg.NotifierConfig.Address)
	notifier := &Notifier{
		cfg:  &cfg.NotifierConfig,
		auth: auth,
		Body: BodyHtmlPtrn,
	}

	return notifier
}

// TODO: test this
func (n *Notifier) Send(to []string, serviceName, message string) error {
	addr := n.cfg.Address + strconv.Itoa(n.cfg.Port)

	for _, toEmail := range to {
		subject := fmt.Sprintf(SubjectTemplate, serviceName)
		htmlBody := fmt.Sprintf(BodyHtmlPtrn, subject, subject, serviceName, message)

		msg := fmt.Sprintf(MsgTmpl, toEmail, subject, htmlBody)

		err := smtp.SendMail(addr, n.auth, n.cfg.From, to, []byte(msg))
		if err != nil {
			return err
		}
	}

	return nil
}
