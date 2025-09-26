package email

import (
	"fmt"
	"net/smtp"
)

// SendEmailSMTP sends an email using basic SMTP credentials.
//
// Parameters:
//   - smtpHost: SMTP server host (e.g., "smtp.gmail.com").
//   - smtpPort: SMTP port (e.g., 587).
//   - username: SMTP username (usually your email address).
//   - password: SMTP password or app password.
//   - from: Email sender address (e.g., "you@example.com").
//   - to: Recipient address (e.g., "user@example.com").
//   - subject: Email subject.
//   - body: Plaintext body of the email.
//
// Example:
//
//	err := SendEmailSMTP(
//		"smtp.gmail.com", 587,
//		"your-email@gmail.com", "app-password",
//		"your-email@gmail.com", "user@example.com",
//		"Welcome", "Hello, welcome to our service!")
func SendEmailSMTP(
	smtpHost string,
	smtpPort int,
	username, password string,
	from, to, subject, body string,
) error {
	auth := smtp.PlainAuth("", username, password, smtpHost)
	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		from, to, subject, body,
	))

	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	return smtp.SendMail(addr, auth, from, []string{to}, msg)
}
