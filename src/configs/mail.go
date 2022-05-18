package configs

import (
	"fmt"
	"net/smtp"
)

// handle error(unencrypted connection) in case server mail not set up TLS => Source code from customer
type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}

func SendEmail(to string, title string, content string) (bool, error) {
	// Get configs info from file configs dev.env
	env := GetEnvConfig()
	from := env.MailAdmin
	pass := env.MailPassword
	host := env.MailSmtp
	port := env.MailPort

	// Setting content
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject:" + title + "\n"
	msg := []byte("To: " + to + "\r\n" + subject + mime + content)

	// Authentication
	auth := unencryptedAuth{
		smtp.PlainAuth(
			"",
			from,
			pass,
			host,
		),
	}
	// Sending mail
	if err := smtp.SendMail(host+":"+port, auth, from, []string{to}, msg); err != nil {
		fmt.Println("Error SendMail: ", err)
		return false, err
	}
	fmt.Println("Email Sent!")
	return true, nil
}
