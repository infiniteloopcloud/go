package email

import (
	"context"
	"crypto/tls"

	gomail "gopkg.in/mail.v2"
)

var _ Descriptor = smtpProvider{}

type SmtpConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type smtpProvider struct {
	d *gomail.Dialer
}

func newSmtpProvider(config SmtpConfig) smtpProvider {
	d := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return smtpProvider{
		d: d,
	}
}

func (s smtpProvider) Send(_ context.Context, data Data) (Response, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", data.Sender)
	m.SetHeader("To", data.ToAddresses...)
	m.SetHeader("Subject", data.Subject)
	for _, ccAddr := range data.CCAddresses {
		m.SetAddressHeader("Cc", ccAddr, ccAddr)
	}
	if data.Body != "" {
		m.SetBody("plain/text", data.Body)
	} else {
		m.SetBody("text/html", data.BodyHTML)
	}
	if err := s.d.DialAndSend(m); err != nil {
		return Response{}, err
	}
	return Response{}, nil
}
