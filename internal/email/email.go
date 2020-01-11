package email

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendLink(emailAddress string, allEmailAddreses string, codename string, linkString string, round int) error {
	plainTextContent := "It's round " + strconv.Itoa(round) + " in your Sabacc game against " + allEmailAddreses + "! Click here to take your turn, " + emailAddress + "!" + linkString
	htmlContent := `It's round ` + strconv.Itoa(round) + ` in your game against ` + allEmailAddreses + `!<br><br><a href="` + linkString + `">Click here to take your turn, ` + emailAddress + `!</a>`
	return SendMessage(emailAddress, codename, plainTextContent, htmlContent)
}

func SendConfirmation(emailAddress string, codename string, hand string, score string) error {
	message := "Your turn has been recorded. Your hand is currently " + hand + " with a score of " + score + ". Please wait patiently for the next email alerting you that it's your turn."
	return SendMessage(emailAddress, codename, message, message)
}

func SendMessage(toEmailAddress string, codename string, messagePlain string, messageHTML string) error {
	if len(os.Getenv("SABACC_DEBUG")) == 0 {
		from := mail.NewEmail("Sabaac Dealer", getFromEmailAddress(codename))
		subject := "Your Sabacc Game (Codename: " + codename + ")"
		to := mail.NewEmail("Sabacc Player", toEmailAddress)

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

func getFromEmailAddress(codename string) string {
	return strings.ToLower(strings.ReplaceAll(codename, " ", "-")) + "@jessemillar.com"
}
