package v2

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
	}
)

func doSendEmail(usr User, h HashData, template string) error {
	tmpl := templateData{Name: usr.Name, URL: h.Hash}
	r := newRequest([]string{h.Email}, "Активація аккаунту", "Активація аккаунту")
	if err := r.ParseTemplate("templates/"+template+".html", tmpl); err != nil {
		return err
	}
	ok, err := r.SendEmail()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("email not sent")
	}
	return nil
}

func newRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Request) SendEmail() (bool, error) {
	r.from = "noreply@mtest.com.ua"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	from := "From:" + r.from + "\n"
	msg := []byte(from + subject + mime + "\n" + r.body)
	addr := "mail.adm.tools:2525"

	auth := smtp.PlainAuth("M-TEST", "noreply@mtest.com.ua", "test", "mail.adm.tools")
	if err := smtp.SendMail(addr, auth, r.from, r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
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
