package sender

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
)

type EmailData struct {
	Username   string
	Email      string
	ActiveCode string
}

func SendActiveEmail(username, email, code string) {
	from := "golangshoptest.gmail.com"
	password := "19045522"
	to := []string{email}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	data := EmailData{
		Username:   username,
		Email:      email,
		ActiveCode: code,
	}

	auth := smtp.PlainAuth("", from, password, smtpHost)

	tmp, _ := template.ParseFiles("./templates/emailActive.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("فعالسازی حساب کاربری \n%s\n\n", mimeHeaders)))

	tmp.Execute(&body, data)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		return
	}
	return
}
