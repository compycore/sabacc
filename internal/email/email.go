package email

import (
	"errors"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendLink(emailAddress string, linkString string) error {
	plainTextContent := "Click here to take your turn, " + emailAddress + "!" + linkString
	htmlContent := `<a href="` + linkString + `">Click here to take your turn, ` + emailAddress + `!</a>`
	return SendMessage(emailAddress, plainTextContent, htmlContent)
}

func SendMessage(emailAddress string, messagePlain string, messageHTML string) error {
	log.Println("Composing email")

	from := mail.NewEmail("Sabaac Dealer", "sabaac@jessemillar.com")
	subject := "Your Sabacc Game"
	to := mail.NewEmail("Sabacc Player", emailAddress)

	message := mail.NewSingleEmail(from, subject, to, messagePlain, messageHTML)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	response, err := client.Send(message)
	if err != nil {
		return errors.New("Email error: " + err.Error())
	}

	// Catch issues with the email API
	if response.StatusCode != 202 {
		return errors.New(response.Body)
	}

	log.Println("Email sent to " + emailAddress + " with message " + messagePlain)

	return nil
}
