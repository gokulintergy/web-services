// Package notification handles sending notifications. At present it just wraps
// the 8o8/email package so that calles are agnostic as to how the
// notificfations are implemented.
package notification

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/8o8/email"
)

// Email is is a copy of email.Email
type Email email.Email

// Send sends an email using the default mx service specified in the .env
func (e Email) Send() error {

	// get the preferred mx service from the env
	mx := os.Getenv("MAPPCPD_MX_SERVICE")
	if mx == "" {
		return errors.New("notification.Send() could not get the preferred MX service from env var MAPPCPD_MX_SERVICE")
	}

	switch strings.ToLower(mx) {
	case "mailgun":
		fmt.Println("Sending email via Mailgun")
		return e.sendMailgun(
			os.Getenv("MAILGUN_API_KEY"),
			os.Getenv("MAILGUN_DOMAIN"),
		)
	case "sendgrid":
		fmt.Println("Sending email via Sendgrid")
		return e.sendSendgrid(
			os.Getenv("SENDGRID_API_KEY"),
		)
	case "ses":
		fmt.Println("Sending email via SES")
		return e.sendSES(
			os.Getenv("AWS_SES_REGION"),
			os.Getenv("AWS_SES_ACCESS_KEY_ID"),
			os.Getenv("AWS_SES_SECRET_ACCESS_KEY"),
		)
	}

	return fmt.Errorf("notification.Send() unknown value for MAPPCPD_MX_SERVICE %q", mx)
}

// sendSES sends the email with Amazon SES
func (e Email) sendSES(awsRegion, awsAccessKeyID, awsSecretAccessKey string) error {
	cfg := email.SESCfg{
		AWSRegion:          awsRegion,
		AWSAccessKeyID:     awsAccessKeyID,
		AWSSecretAccessKey: awsSecretAccessKey,
	}
	sndr, err := email.NewSES(cfg)
	if err != nil {
		return err
	}

	return sndr.Send(email.Email(e)) // cast local Email type to email.Email
}

// sendMailgun sends the email using Mailgun
func (e Email) sendMailgun(apiKey, domain string) error {
	cfg := email.MailgunCfg{
		APIKey: apiKey,
		Domain: domain,
	}
	sndr := email.NewMailgun(cfg)

	return sndr.Send(email.Email(e))
}

// sendSendgrid sends the email using Sendgrid
func (e Email) sendSendgrid(apiKey string) error {
	cfg := email.SendgridCfg{
		APIKey: apiKey,
	}
	sndr := email.NewSendgrid(cfg)

	return sndr.Send(email.Email(e))
}
