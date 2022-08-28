package ses

import (
	"context"
	"fmt"
	"mlock/lambdas/shared"
	mshared "mlock/shared"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type EmailService struct {
	c *ses.Client
}

func NewEmailService(ctx context.Context) (*EmailService, error) {
	c, err := getClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %s", err.Error())
	}

	return &EmailService{
		c: c,
	}, nil
}

func getClient(ctx context.Context) (*ses.Client, error) {
	// Might want to refactor this more to store the `EmailService` in the context (if we need to store anything in there).
	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting context data: %s", err.Error())
	}

	if cd.SES != nil {
		return cd.SES, nil
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-1"))
	if err != nil {
		return nil, fmt.Errorf("error getting aws config: %s", err.Error())
	}

	cd.SES = ses.NewFromConfig(cfg)

	return cd.SES, nil
}

func (s *EmailService) SendEmailToAdmins(ctx context.Context, subject string, body string) error {
	return s.SendEmail(
		ctx,
		subject,
		body,
		strings.Split(mshared.GetConfigUnsafe("EMAIL_TO_ADMINS"), ";"),
	)
}

func (s *EmailService) SendEmailToDevelopers(ctx context.Context, subject string, body string) error {
	return s.SendEmail(
		ctx,
		subject,
		body,
		strings.Split(mshared.GetConfigUnsafe("EMAIL_TO_DEVELOPERS"), ";"),
	)
}

func (s *EmailService) SendEmail(
	ctx context.Context,
	subject string,
	body string,
	tos []string,
) error {
	from, err := mshared.GetConfig("EMAIL_FROM_ADDRESS")
	if err != nil {
		return fmt.Errorf("empty from address")
	}

	characterSet := "UTF-8"

	_, err = s.c.SendEmail(ctx, &ses.SendEmailInput{
		Source: &from,
		Destination: &types.Destination{
			ToAddresses: tos,
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: &characterSet,
					Data:    &body,
				},
			},
			Subject: &types.Content{
				Charset: &characterSet,
				Data:    &subject,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error sending email: %s", err.Error())
	}

	return nil
}
