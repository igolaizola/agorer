package mail

import (
	"context"
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Dry      bool
}

func Send(ctx context.Context, cfg *Config, sender, recipient, subject, body, file string) error {
	// Create a new message.
	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	m.Attach(file)

	// Send the email
	log.Println("Sending email...")
	fmt.Println("From:", sender)
	fmt.Println("To:", recipient)
	fmt.Println("Subject:", subject)
	fmt.Println("Body:", body)
	fmt.Println("File:", file)

	if cfg.Dry {
		log.Println("Dry run, not sending email")
		return nil
	}

	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	errC := make(chan error)
	defer close(errC)
	go func() {
		errC <- d.DialAndSend(m)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errC:
		if err != nil {
			return fmt.Errorf("couldn't send email: %w", err)
		}
	}
	return nil
}
