package email

import (
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

// Email interface
type EmailSender interface {
	SendOTP(to string, otp string) error
}

// Sender struct
type Sender struct {
	dialer *gomail.Dialer
}

var _ EmailSender = (*Sender)(nil)

// NewEmailSender creates a new email sender
func NewEmailSender(host string, port int, username string, password string) *Sender {
	return &Sender{
		dialer: gomail.NewDialer(host, port, username, password),
	}
}

// SendOTP sends an OTP to the user
func (s *Sender) SendOTP(to string, otp string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "noreply@localhost.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Your OTP")
	m.SetBody("text/html", fmt.Sprintf("Your OTP is: %s", otp))

	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("Failed to send OTP email to %s: %v", to, err)
		return err
	}

	log.Printf("OTP email sent to %s successfully", to)

	return nil
}
