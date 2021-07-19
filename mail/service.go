package mail

import (
	"bytes"
	"crypto/tls"
	"errors"
	"html/template"
	"net"
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

	return true, send(r.from, r.to[0], "smtp.ukr.net:465", r.auth.Email, r.auth.Password, msg)
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

func send(from, to, servername, username, password string, msg []byte) error {
	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", username, password, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(from); err != nil {
		return err
	}

	if err = c.Rcpt(to); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}
