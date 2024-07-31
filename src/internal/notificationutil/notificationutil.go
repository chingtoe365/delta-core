package notificationutil

import (
	"delta-core/bootstrap"
	"log/slog"
	"net/smtp"
)

func SendMail(env *bootstrap.Env, to string, body string) {
	slog.Info("Sending email")
	from := env.GmailSender
	password := env.GmailSenderAppPassword
	slog.Info(from)
	slog.Info(password)

	msg := "From: " + from + "\n" + "To: " + to + "\n" + "Subject: Delta signal detected\n\n" + body
	slog.Info("Ready to send email")
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{to}, []byte(msg),
	)
	slog.Info("email just sent")
	if err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info("Email sent succesfully!")

}
