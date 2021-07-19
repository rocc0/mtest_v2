package mail

import (
	"bytes"
	"errors"
	"html/template"
	"net/smtp"
)

type (
	templateData struct {
		Name string
		URL  string
	}
	Request struct {
		from    string
		to      []string
		subject string
		body    string
		auth    Auth
	}
)

type Auth struct {
	Email    string
	Password string
}

func SendEmail(name, email, hash, template string, auth Auth) error {
	tmpl := templateData{Name: name, URL: hash}
	r := newRequest([]string{email}, "Активація аккаунту", "Активація аккаунту", auth)
	if err := r.parseTemplate("templates/"+template+".html", tmpl); err != nil {
		return err
	}
	ok, err := r.send()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("email not sent")
	}
	return nil
}

func newRequest(to []string, subject, body string, auth Auth) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
		auth:    auth,
	}
}

func (r *Request) send() (bool, error) {
	r.from = "noreply@mtest.com.ua"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	from := "From:" + r.from + "\n"
	msg := []byte(from + subject + mime + "\n" + r.body)
	addr := "smtp.ukr.net:2525"

	auth := smtp.PlainAuth("M-TEST", r.auth.Email, r.auth.Password, "smtp.ukr.net")
	if err := smtp.SendMail(addr, auth, r.from, r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Request) parseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
