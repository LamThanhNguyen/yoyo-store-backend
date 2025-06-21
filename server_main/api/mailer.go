package api

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"time"

	"github.com/rs/zerolog/log"

	mail "github.com/xhit/go-simple-mail/v2"
)

//go:embed templates
var emailTemplateFS embed.FS

func (server *Server) SendMail(
	from,
	to,
	subject,
	tmpl string,
	data interface{},
) error {
	templateToRender := fmt.Sprintf("templates/%s.html.tmpl", tmpl)

	t, err := template.New("email-html").ParseFS(emailTemplateFS, templateToRender)
	if err != nil {
		log.Error().Err(err).Msg("SendMail")
		return err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		log.Error().Err(err).Msg("SendMail")
		return err
	}

	formattedMessage := tpl.String()

	templateToRender = fmt.Sprintf("templates/%s.plain.tmpl", tmpl)
	t, err = template.New("email-plain").ParseFS(emailTemplateFS, templateToRender)
	if err != nil {
		log.Error().Err(err).Msg("SendMail")
		return err
	}

	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		log.Error().Err(err).Msg("SendMail")
		return err
	}

	plainMessage := tpl.String()

	// send the mail
	mail_server := mail.NewSMTPClient()
	mail_server.Host = server.config.SmtpHost
	mail_server.Port = server.config.SmtpPort
	mail_server.Username = server.config.SmtpUsername
	mail_server.Password = server.config.SmtpPassword
	mail_server.Encryption = mail.EncryptionTLS
	mail_server.KeepAlive = false
	mail_server.ConnectTimeout = 10 * time.Second
	mail_server.SendTimeout = 10 * time.Second

	smtpClient, err := mail_server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(from).
		AddTo(to).
		SetSubject(subject)

	email.SetBody(mail.TextHTML, formattedMessage)
	email.AddAlternative(mail.TextPlain, plainMessage)

	err = email.Send(smtpClient)
	if err != nil {
		log.Error().Err(err).Msg("SendMail")
		return err
	}

	log.Info().Msg("send email")

	return nil
}
