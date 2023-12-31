package notificationutil

import (
	"delta-core/bootstrap"
	"fmt"
	"net/smtp"
)

func SendMail(env *bootstrap.Env, to string, body string) {
	fmt.Println("Sending email")
	from := env.GmailSender
	password := env.GmailSenderAppPassword
	fmt.Println(from)
	fmt.Println(password)

	msg := "From: " + from + "\n" + "To: " + to + "\n" + "Subject: Delta signal detected\n\n" + body
	fmt.Println("Ready to send email")
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{to}, []byte(msg),
	)
	fmt.Println("email just sent")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email sent succesfully!")

}
