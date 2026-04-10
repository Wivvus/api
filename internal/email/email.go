package email

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendPasswordReset(toEmail, link string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	if from == "" {
		from = user
	}

	subject := "Reset your Wivvus password"
	body := fmt.Sprintf(`Hi,

You requested a password reset for your Wivvus account. Click the link below to set a new password:

%s

This link expires in 24 hours.

If you didn't request this, you can ignore this email.`, link)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, toEmail, subject, body)

	addr := host + ":" + port
	var auth smtp.Auth
	if user != "" && pass != "" {
		auth = smtp.PlainAuth("", user, pass, host)
	}

	return smtp.SendMail(addr, auth, from, []string{toEmail}, []byte(msg))
}

func SendVerification(toEmail, name, verificationLink string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	if from == "" {
		from = user
	}

	subject := "Verify your email for Wivvus"
	body := fmt.Sprintf(`Hi %s,

Thanks for signing up for Wivvus. Please verify your email and set your password by clicking the link below:

%s

This link expires in 24 hours.

If you didn't create an account, you can ignore this email.`, name, verificationLink)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, toEmail, subject, body)

	addr := host + ":" + port
	var auth smtp.Auth
	if user != "" && pass != "" {
		auth = smtp.PlainAuth("", user, pass, host)
	}

	return smtp.SendMail(addr, auth, from, []string{toEmail}, []byte(msg))
}
