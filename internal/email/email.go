package email

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"text/template"

	"github.com/compycore/sabacc/internal/models"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

func SendLink(emailAddress string, allEmailAddreses string, codename string, linkString string, round int) error {
	plainTextContent := "It's round " + strconv.Itoa(round) + " in your Sabacc game against " + allEmailAddreses + "! Click here to take your turn, " + emailAddress + "!" + linkString
	htmlContent := `It's round ` + strconv.Itoa(round) + ` in your game against ` + allEmailAddreses + `!<br><br><a href="` + linkString + `">Click here to take your turn, ` + emailAddress + `!</a>`
	return SendMessage(emailAddress, codename, plainTextContent, htmlContent)
}

func SendGameStartNotice(database models.Database) error {
	message, err := executeTemplate(database, "game-start.html")
	if err != nil {
		return err
	}

	for _, player := range database.AllPlayers {
		err = SendMessage(player.Email, database.Codename, message, message)
		if err != nil {
			return err
		}
	}

	return nil
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

// ExecuteTemplate takes data passed and uses it to execute the specified Go template and returns the result as a string
func executeTemplate(database models.Database, templateFile string) (string, error) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	parsedTemplate, err := template.ParseFiles(filepath.Join(basepath, "./templates/"+templateFile))
	if err != nil {
		return "", err
	}

	var content bytes.Buffer
	parsedTemplate.Execute(&content, database)

	return content.String(), nil
}
