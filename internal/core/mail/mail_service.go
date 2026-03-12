package mail

import (
	"context"

	brevo "github.com/getbrevo/brevo-go/lib"
	"github.com/itsLeonB/ungerr"
	"github.com/reflect-homini/stora/internal/core/config"
	"github.com/reflect-homini/stora/internal/core/otel"
)

type Service interface {
	Send(ctx context.Context, msg MailMessage) error
}

type MailMessage struct {
	RecipientMail string `validate:"required,min=3,email"`
	RecipientName string `validate:"required,min=1"`
	Subject       string `validate:"required,min=3"`
	HTMLContent   string
	TextContent   string
}

type brevoMailService struct {
	client     *brevo.APIClient
	senderMail string
	senderName string
}

func NewMailService(mailCfg config.Mail) Service {
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", mailCfg.ApiKey)
	br := brevo.NewAPIClient(cfg)
	return &brevoMailService{br, mailCfg.SenderMail, mailCfg.SenderName}
}

func (ms *brevoMailService) Send(ctx context.Context, msg MailMessage) error {
	ctx, span := otel.Tracer.Start(ctx, "brevoMailService.Send")
	defer span.End()

	mail := brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Name:  ms.senderName,
			Email: ms.senderMail,
		},
		To: []brevo.SendSmtpEmailTo{{
			Email: msg.RecipientMail,
			Name:  msg.RecipientName,
		}},
		Subject:     msg.Subject,
		HtmlContent: msg.HTMLContent,
		TextContent: msg.TextContent,
	}

	if _, _, err := ms.client.TransactionalEmailsApi.SendTransacEmail(ctx, mail); err != nil {
		return ungerr.Wrap(err, "error sending email")
	}
	return nil
}
