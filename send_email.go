package main

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

func SendSimpleEmail(subject string, body string,attachment io.ReadCloser, filename string) (string, error) {
	apiKey := os.Getenv("MAILGUN_API_KEY")
	emailFrom := os.Getenv("MAILGUN_FROM")
	emailTo := os.Getenv("MAILGUN_TO")
	domain := os.Getenv("MAILGUN_DOMAIN")
	mg := mailgun.NewMailgun(domain, apiKey)

	m := mailgun.NewMessage(
		emailFrom,
		subject,
		body,
		emailTo,
	)

	// Add attachment if provided
	if attachment != nil {
		defer attachment.Close()
		m.AddReaderAttachment(filename, attachment)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)

	return id, err
}
