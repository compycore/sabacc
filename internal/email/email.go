package email

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/compycore/sabacc/internal/models"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

func SendLink(database models.Database) error {
	message, err := executeTemplate(database, "your-turn")
	if err != nil {
		return err
	}

	return SendMessageToOnePlayer(database, message)
}

func SendGameStartNotice(database models.Database) error {
	message, err := executeTemplate(database, "game-start")
	if err != nil {
		return err
	}

	err = SendMessageToAllPlayers(database, message)
	if err != nil {
		return err
	}

	return nil
}

func SendGameOver(database models.Database) error {
	message, err := executeTemplate(database, "game-over")
	if err != nil {
		return err
	}

	err = SendMessageToAllPlayers(database, message)
	if err != nil {
		return err
	}

	return nil
}

func SendConfirmation(database models.Database) error {
	message, err := executeTemplate(database, "turn-confirmation")
	if err != nil {
		return err
	}

	return SendMessageToOnePlayer(database, message)
}

func SendHandDiscardNotice(database models.Database) error {
	message, err := executeTemplate(database, "discard-notice")
	if err != nil {
		return err
	}

	return SendMessageToOnePlayer(database, message)
}

func SendMessageToAllPlayers(database models.Database, message string) error {
	for _, player := range database.AllPlayers {
		err := sendMessage(player.Email, database.Codename, message)
		if err != nil {
			return err
		}
	}

	return nil
}

func SendMessageToOnePlayer(database models.Database, message string) error {
	return sendMessage(database.AllPlayers[database.Turn].Email, database.Codename, message)
}

func sendMessage(toEmailAddress string, codename string, message string) error {
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
				TextPart: "Please use a viewer capable of rendering HTML.",
				HTMLPart: message,
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

func sliceToSentence(words []string, andOrOr string) string {
	l := len(words)
	wordsForSentence := make([]string, l)
	copy(wordsForSentence, words)
	wordsForSentence[l-1] = andOrOr + " " + wordsForSentence[l-1]
	return strings.Join(wordsForSentence, ", ")
}

func getAllPlayerNamesAsSentence(database models.Database) string {
	names := []string{}

	for _, player := range database.AllPlayers {
		names = append(names, player.Name)
	}

	return sliceToSentence(names, "and")
}

func EncodeDatabase(database models.Database) (string, error) {
	encodedDatabase, err := json.Marshal(database)
	if err != nil {
		return "", err
	}

	return url.QueryEscape(string(encodedDatabase)), nil
}

func createGameLink(database models.Database) (models.Database, error) {
	// Generate a rematch link if the game is over
	// if controllers.IsGameOver(database) {
	// stringified, err = createRematchLink(database)
	// if err != nil {
	// return models.Database{}, err
	// }
	// } else {
	stringified, err := EncodeDatabase(database)
	if err != nil {
		return models.Database{}, err
	}
	// }

	database.Template.Link = os.Getenv("SABACC_UI_HREF") + "?" + stringified
	return database, nil
}

func createRematchLink(database models.Database) (string, error) {
	rematchDatabase := models.Database{}
	rematchDatabase.Rematch = database.AllPlayers

	rematchDatabaseString, err := EncodeDatabase(rematchDatabase)
	if err != nil {
		return "", err
	}

	return `<a href="` + os.Getenv("SABACC_UI_HREF") + "?" + rematchDatabaseString + `">Click here for a rematch!</a>`, nil
}

// prepDatabaseForTemplate runs some functions to populate the database struct with values useful for populating templates
func prepDatabaseForTemplate(database models.Database) (models.Database, error) {
	// Create the HTTP link representing the next player's turn (or a rematch)
	database, err := createGameLink(database)
	if err != nil {
		return models.Database{}, err
	}

	// Populate the nice sentence of all the player names
	database.Template.AllPlayersNames = getAllPlayerNamesAsSentence(database)

	return database, nil
}

// executeTemplate takes data passed and uses it to execute the specified Go template and returns the result as a string
func executeTemplate(database models.Database, templateName string) (string, error) {
	database, err := prepDatabaseForTemplate(database)
	if err != nil {
		return "", err
	}

	// Do some magic to find the path we're running from
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	// Compile all templates
	templates := template.Must(template.ParseGlob(filepath.Join(basepath, "./templates/*")))

	// Actually execute a template (by template name, not filename)
	var content bytes.Buffer
	err = templates.ExecuteTemplate(&content, templateName, database)
	if err != nil {
		return "", err
	}

	// Return the compiled template as a string
	return content.String(), nil
}
