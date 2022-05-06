package main

import (
	"bytes"
	"html/template"
	"log"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From         string
	FromName     string
	To           string
	Subject      string
	Attachements []string
	Data         any
	DataMap      map[string]any
}

func (m *Mail) SendSMTPMessage(msg Message) error {
	logSnippet := "\n[mail-service][mailer][SendSMTPMessage] =>"

	if msg.From == "" {
		msg.From = m.FromAddress
	}
	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		log.Printf("%s (ERROR-server.Connect): %s", logSnippet, err.Error())
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)
	if len(msg.Attachements) > 0 {
		for _, x := range msg.Attachements {
			email.AddAttachment(x)
		}
	}

	err = email.Send(smtpClient)
	if err != nil {
		log.Printf("%s (ERROR- email.Send): %s", logSnippet, err.Error())
		return err
	}

	return nil
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	logSnippet := "\n[mail-service][mailer][buildHTMLMessage] =>"

	templateToRender := "./templates/mail.html.gohtml"

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		log.Printf("%s (ERROR-template.New.ParseFiles): %s", logSnippet, err.Error())
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		log.Printf("%s (ERROR-t.ExecuteTemplate): %s", logSnippet, err.Error())
		return "", err
	}

	formattedMessage := tpl.String()
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		log.Printf("%s (ERROR-m.inlineCSS): %s", logSnippet, err.Error())
		return "", err
	}

	return formattedMessage, nil
}

func (m *Mail) inlineCSS(s string) (string, error) {
	logSnippet := "\n[mail-service][mailer][inlineCSS] =>"

	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		log.Printf("%s (ERROR-premailer.NewPremailerFromString): %s", logSnippet, err.Error())
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		log.Printf("%s (ERROR-prem.Transform): %s", logSnippet, err.Error())
		return "", err
	}

	return html, nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	logSnippet := "\n[mail-service][mailer][buildPlainTextMessage] =>"

	templateToRender := "./templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		log.Printf("%s (ERROR-template.New.ParseFiles): %s", logSnippet, err.Error())
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		log.Printf("%s (ERROR-t.ExecuteTemplate): %s", logSnippet, err.Error())
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
