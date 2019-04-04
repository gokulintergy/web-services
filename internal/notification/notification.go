// Package notification handles sending notifications. At present it just wraps
// the 8o8/email package so that calles are agnostic as to how the
// notificfations are implemented.
package notification

import (
	"os"

	"github.com/8o8/email"
)

// Email is is a copy of email.Email
type Email email.Email

// Send sends an email using a particular sender
func (e Email) Send() error {
	return e.sendSES(
		os.Getenv("AWS_SES_REGION"),
		os.Getenv("AWS_SES_ACCESS_KEY_ID"),
		os.Getenv("AWS_SES_SECRET_ACCESS_KEY"),
	)
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
