package notificationutil

import (
	"fmt"
	"net/smtp"
)

func SendMail(to string, body string) {
	fmt.Println("Sending email")
	from := "chingtoe@gmail.com"
	password := "hpmx pize tnvr msew"

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
