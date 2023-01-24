package email

import (
	"context"
	"errors"
)

const (
	AWSSesProvider provider = "aws-ses"
	SMTPProvider   provider = "smtp"

	DefaultCharset = "UTF-8"
)

type provider string

func NewProvider(ctx context.Context, o Opts) (Descriptor, error) {
	var d Descriptor
	var err error
	switch o.Provider {
	case SMTPProvider:
		if o.SmtpConfig == nil {
			return nil, errors.New("missing smtp config")
		}
		d = newSmtpProvider(*o.SmtpConfig)
	default:
		d, err = newAwsSesProvider(ctx)
	}
	if err != nil {
		return nil, err
	}

	return d, nil
}

type Descriptor interface {
	Send(ctx context.Context, data Data) (Response, error)
}

type Opts struct {
	Provider   provider
	SmtpConfig *SmtpConfig
}

type Data struct {
	Sender   string
	Subject  string
	BodyHTML string
	Body     string
	// Charset of the email, default UTF-8
	Charset     string
	ToAddresses []string
	CCAddresses []string
}

func (d *Data) Prepare() {
	if d.Charset == "" {
		d.Charset = DefaultCharset
	}
}

func (d Data) Validate() error {
	if d.Sender == "" {
		return errors.New("sender address should be present")
	}
	if d.Subject == "" {
		return errors.New("subject should be present")
	}
	if d.BodyHTML == "" && d.Body == "" {
		return errors.New("body should be present")
	}
	if len(d.ToAddresses) < 1 {
		return errors.New("at least one to address should be present")
	}

	return nil
}

type Response struct {
	MessageID string
}
