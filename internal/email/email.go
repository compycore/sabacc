package email

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	mailjet "github.com/mailjet/mailjet-apiv3-go"
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

func SendHandDiscardNotice(emailAddress string, codename string, hand string, score string) error {
	message := "Someone rolled matching dice so everyone's hands were discarded. Your new hand is " + hand + " with a score of " + score + ". Please wait patiently for the next email alerting you that it's your turn."
	return SendMessage(emailAddress, codename, message, message)
}

func SendMessage(toEmailAddress string, codename string, messagePlain string, messageHTML string) error {
	if len(os.Getenv("SABACC_DEBUG")) == 0 {
		mailjetClient := mailjet.NewMailjetClient(os.Getenv("MAILJET_API_KEY_PUBLIC"), os.Getenv("MAILJET_API_KEY_PRIVATE"))
		messagesInfo := []mailjet.InfoMessagesV31{
			mailjet.InfoMessagesV31{
				From: &mailjet.RecipientV31{
					Email: getFromEmailAddress(codename),
					Name:  "Sabacc Dealer",
				},
				To: &mailjet.RecipientsV31{
					mailjet.RecipientV31{
						Email: toEmailAddress,
						Name:  "Sabacc Player",
					},
				},
				Subject:  "Your Sabacc Game (Codename: " + codename + ")",
				TextPart: messagePlain,
				HTMLPart: messageHTML,
			},
		}
		messages := mailjet.MessagesV31{Info: messagesInfo}
		res, err := mailjetClient.SendMailV31(&messages)
		if err != nil {
			return errors.New("Email error: " + err.Error())
		}

		fmt.Printf("Data: %+v\n", res)
	}

	return nil
}

func getFromEmailAddress(codename string) string {
	return strings.ToLower(strings.ReplaceAll(codename, " ", "-")) + "@compycore.com"
}
