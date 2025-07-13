package util

import (
	"fmt"
	"net/smtp"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func SendResetPasswordEmail(cfg EmailConfig, to, link string) error {
	host := cfg.Host
	port := cfg.Port
	auth := smtp.PlainAuth(
		"",
		cfg.Username,
		cfg.Password,
		host,
	)

	subject := "Password Reset Request"
	body := fmt.Sprintf("Hello,\n\nTo reset your password, please click the link below:\n\n%s\n\nIf you did not request this, please ignore.\n", link)

	msg := []byte(
		fmt.Sprintf("From: %s\r\n", cfg.From) +
			fmt.Sprintf("To: %s\r\n", to) +
			fmt.Sprintf("Subject: %s\r\n\r\n", subject) +
			body,
	)

	addr := fmt.Sprintf("%s:%d", host, port)
	return smtp.SendMail(addr, auth, cfg.From, []string{to}, msg)
}
