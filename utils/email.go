package utils

import (
	"net/smtp"
)

func SendEmail(to, subject, body string) error {
	from := "pesan.reski@gmail.com"
	password := "password"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte("Subject: "+subject+"\r\n\r\n"+body))
	if err != nil {
		return err
	}

	return nil
}
