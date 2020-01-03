package email

import (
	"errors"
	"os"
	"strconv"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendLink(emailAddress string, allEmailAddreses string, linkString string, round int) error {
	plainTextContent := "It's round " + strconv.Itoa(round) + " in your Sabacc game against " + allEmailAddreses + "! Click here to take your turn, " + emailAddress + "!" + linkString
	htmlContent := `It's round ` + strconv.Itoa(round) + ` in your game against ` + allEmailAddreses + `!<br><br><a href="` + linkString + `">Click here to take your turn, ` + emailAddress + `!</a>`
	return SendMessage(emailAddress, plainTextContent, htmlContent)
}

func SendConfirmation(emailAddress string, hand string, score string) error {
	message := "Your turn has been recorded. Your hand is currently " + hand + " with a score of " + score + "."
	return SendMessage(emailAddress, message, message)
}

func SendMessage(emailAddress string, messagePlain string, messageHTML string) error {
	if len(os.Getenv("SABACC_DEBUG")) == 0 {
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
	}

	return nil
}
