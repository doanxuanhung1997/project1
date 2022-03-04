package mail

import (
	"fmt"
	"net/smtp"
	"sandexcare_backend/helpers/config"
	"sandexcare_backend/helpers/constant"
	"strconv"
	"strings"
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

func GetHtmlContentCreateListener(roleId int, phoneNumber string, password string) string {
	receiver := UserRequest
	switch roleId {
	case constant.RoleListener:
		receiver = ListenerRequest
		break
	case constant.RoleExperts:
		receiver = ExpertsRequest
	default:
		break
	}
	emailTemplate := GetMailFromDb(TypeCreateListener)
	replacer := strings.NewReplacer("{receiver}", receiver, "{phone_number}", phoneNumber, "{password}", password)
	emailTemplate = replacer.Replace(emailTemplate)
	return emailTemplate
}

func GetHtmlContentResetPassword(roleId int, phoneNumber string, code string) string {
	receiver := UserRequest
	switch roleId {
	case constant.RoleListener:
		receiver = ListenerRequest
		break
	case constant.RoleExperts:
		receiver = ExpertsRequest
	default:
		break
	}
	emailTemplate := GetMailFromDb(TypeForgotPassword)
	replacer := strings.NewReplacer("{receiver}", receiver, "{phone_number}", phoneNumber, "{code}", code)
	emailTemplate = replacer.Replace(emailTemplate)
	return emailTemplate
}

func GetHtmlContentBlockListener() string {
	emailTemplate := GetMailFromDb(TypeNotifyBlockListener)
	return emailTemplate
}

func GetMailFromDb(typeEmail string) string {
	/*Get template from db*/
	emailDb := NewResource().GetEmailTemplate(typeEmail)
	return emailDb.Template
}

/*Get email subject with type corresponding*/
func GetSubject(typeTemplate string) (subject string) {
	switch typeTemplate {
	case TypeForgotPassword:
		subject = "Request a password reset"
		break
	case TypeCreateListener:
		subject = "Create account successful"
		break
	case TypeNotifyBlockListener:
		subject = "Blocked account"
		break
	}
	return subject
}
