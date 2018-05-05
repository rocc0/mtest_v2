package main

import (
	"bytes"
	"html/template"
	"net/smtp"
	"log"
)

var auth smtp.Auth

func doSendEmail(usr User, h Hash, template string) error {
	auth = smtp.PlainAuth("M-TEST", "noreply@mtest.com.ua", "Qzwxec1425", "mail.adm.tools")
	templateData := struct {
		Name string
		URL  string
	}{
		Name: usr.Name,
		URL:  h.Hash,
	}
	r := NewRequest([]string{h.Email}, "Активація аккаунту", "Активація аккаунту")
	if err := r.ParseTemplate("templates/" + template + ".html", templateData); err == nil {
		ok, err := r.SendEmail()
		if err != nil {
			return err
		}
		log.Println(ok)
	} else {
		return err
	}
	return nil
}


type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewRequest(to []string, subject, body string) *Request {
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

