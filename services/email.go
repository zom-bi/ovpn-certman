package services

import (
	"errors"
	"log"
	"time"

	"github.com/go-mail/mail"
)

var (
	ErrMailUninitializedConfig = errors.New("Mail: uninitialized config")
)

type EmailConfig struct {
	From         string
	SMTPServer   string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
}

type Email struct {
	config *EmailConfig

	mailChan chan *mail.Message
}

func NewEmail(conf *EmailConfig) *Email {
	if conf == nil {
		log.Println(ErrMailUninitializedConfig)
	}

	return &Email{
		config:   conf,
		mailChan: make(chan *mail.Message, 0),
	}
}

// Send sends an email to the receiver
func (email *Email) Send(to, subject, text, html string) error {
	if email.config == nil {
		log.Print("Error: trying to send mail with uninitialized config.")
		return ErrMailUninitializedConfig
	}

	m := mail.NewMessage()
	m.SetHeader("From", email.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", text)
	m.AddAlternative("text/html", html)

	// put email in chan
	email.mailChan <- m
	return nil
}

// Daemon is a function that takes Mail and sends it without blocking.
// WIP
func (email *Email) Daemon() {
	if email.config == nil {
		log.Print("Error: trying to set up mail deamon with uninitialized config.")
		return
	}

	d := mail.NewDialer(
		email.config.SMTPServer,
		email.config.SMTPPort,
		email.config.SMTPUsername,
		email.config.SMTPPassword)

	var s mail.SendCloser
	var err error
	open := false
	for {
		select {
		case m, ok := <-email.mailChan:
			if !ok {
				// channel is closed
				return
			}
			if !open {
				if s, err = d.Dial(); err != nil {
					log.Print(err)
					return
				}
				open = true
			}
			if err := mail.Send(s, m); err != nil {
				log.Print(err)
			}
		// Close the connection if no email was sent in the last 30 seconds.
		case <-time.After(30 * time.Second):
			if open {
				if err := s.Close(); err != nil {
					panic(err)
				}
				open = false
			}
		}
	}
}
