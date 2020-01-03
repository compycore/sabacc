package email

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendLink(emailAddress string, linkString string, round int) error {
	plainTextContent := "It's round " + strconv.Itoa(round) + "! Click here to take your turn, " + emailAddress + "!" + linkString
	htmlContent := `It's round ` + strconv.Itoa(round) + `! <a href="` + linkString + `">Click here to take your turn, ` + emailAddress + `!</a>`
	return SendMessage(emailAddress, plainTextContent, htmlContent)
}

func SendMessage(emailAddress string, messagePlain string, messageHTML string) error {
	if len(os.Getenv("SABACC_DEBUG")) == 0 {
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
	}

	return nil
}
