package ses

import (
	"context"
	"fmt"
	"mlock/shared"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func GetClient(ctx context.Context) (*ses.SES, error) {
	cd, err := shared.GetContextData(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting context data: %s", err.Error())
	}

	if cd.SES != nil {
		return cd.SES, nil
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)
	if err != nil {
		return nil, fmt.Errorf("error getting aws session: %s", err.Error())
	}

	cd.SES = ses.New(sess)
	return cd.SES, nil
}

func SendEamil(ctx context.Context, subject string, body string) error {
	s, err := GetClient(ctx)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err.Error())
	}

	from := shared.GetConfig("EMAIL_FROM_ADDRESS")
	if from == "" {
		return fmt.Errorf("empty from address")
	}

	tos := []*string{}
	for _, a := range strings.Split(shared.GetConfig("EMAIL_TO_ADDRESSES"), ";") {
		tos = append(tos, &a)
	}

	characterSet := "UTF-8"

	_, err = s.SendEmailWithContext(ctx, &ses.SendEmailInput{
		Source: &from,
		Destination: &ses.Destination{
			ToAddresses: tos,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: &characterSet,
					Data:    &body,
				},
			},
			Subject: &ses.Content{
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
