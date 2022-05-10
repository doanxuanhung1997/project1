package mail

import (
	"fmt"
	"net/smtp"
	"houze_ops_backend/config"
	"strconv"
)

const (
	TypeForgotPassword      = "forgot_password"
	TypeCreateListener      = "create_listener"
	TypeNotifyBlockListener = "notify_block_listener"

	ListenerRequest = "Listener"
	ExpertsRequest  = "Experts"
	UserRequest     = "User"
)

//Handle error(unencrypted connection) in case server mail not set up TLS => Source code from customer
type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}

func SendEmail(to string, title string, content string) (bool, error) {
	// Get config info from file config dev.env
	env := config.GetEnvValue()
	from := env.Mail.Email
	pass := env.Mail.Password
	host := env.Mail.Smtp
	port := env.Mail.Port

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
	if err := smtp.SendMail(host+":"+strconv.Itoa(port), auth, from, []string{to}, msg); err != nil {
		fmt.Println("Error SendMail: ", err)
		return false, err
	}
	fmt.Println("Email Sent!")
	return true, nil
}

