package main

import (
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func main() {
	// TODO Move to an email package and allow for params
	log.Println("Composing email")

	from := mail.NewEmail("Sabaac Dealer", "sabaac@jessemillar.com")
	subject := "Your Sabacc Game"
	to := mail.NewEmail("Sabacc Player", "hellojessemillar@gmail.com")

	plainTextContent := "Click here to take your turn."
	htmlContent := `It's hellojessemillar@gmail.com's turn!
<br>
<br>
<a href="https://jessemillar.github.io/sabacc?player=test@test.com">Click here if that's you!</a>`

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	}

	// Catch issues with the email API
	if response.StatusCode != 202 {
		log.Println("There was an issue sending the email message:")
		log.Println(response.StatusCode)
		log.Println(response.Headers)
		log.Println(response.Body)
		os.Exit(1)
	}

	log.Println("Email sent")
}
