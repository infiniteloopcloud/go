package email

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	sesTypes "github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type sesProvider struct {
	cfg aws.Config
}

func newAwsSesProvider(ctx context.Context) (*sesProvider, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &sesProvider{
		cfg: cfg,
	}, nil
}

func (s sesProvider) Send(ctx context.Context, data Data) (Response, error) {
	svc := ses.NewFromConfig(s.cfg)

	if err := data.Validate(); err != nil {
		return Response{}, err
	}

	data.Prepare()

	res, err := svc.SendEmail(ctx, createEmailInput(data))
	if err != nil {
		return Response{}, err
	}

	return Response{
		MessageID: aws.ToString(res.MessageId),
	}, nil
}

func createEmailInput(data Data) *ses.SendEmailInput {
	return &ses.SendEmailInput{
		Destination: &sesTypes.Destination{
			CcAddresses: data.CCAddresses,
			ToAddresses: data.ToAddresses,
		},
		Message: &sesTypes.Message{
			Body: &sesTypes.Body{
				Html: &sesTypes.Content{
					Data:    aws.String(data.BodyHTML),
					Charset: aws.String(data.Charset),
				},
				Text: &sesTypes.Content{
					Data:    aws.String(data.Body),
					Charset: aws.String(data.Charset),
				},
			},
			Subject: &sesTypes.Content{
				Data:    aws.String(data.Subject),
				Charset: aws.String(data.Charset),
			},
		},
		Source: aws.String(data.Sender),
	}
}
