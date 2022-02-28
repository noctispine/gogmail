package gservice

// GmailService : Gmail client for sending email

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// https://medium.com/wesionary-team/sending-emails-with-go-golang-using-smtp-gmail-and-oauth2-185ee12ab306

var GmailService *gmail.Service

type OAuthInfos struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
	AccessToken  string
}

type Email struct {
	To      string
	Cc      []string
	Subject string
	Body    string
}

func OAuthGmailService(infos OAuthInfos) {
	config := oauth2.Config{
		ClientID:     infos.ClientID,
		ClientSecret: infos.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
	}

	token := oauth2.Token{
		AccessToken:  infos.AccessToken,
		RefreshToken: infos.RefreshToken,
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Printf("Unable to retrieve Gmail client: %v", err)
	}

	GmailService = srv
}

func SendEmailOAUTH2(email Email) (bool, error) {

	var message gmail.Message
	var err error

	cc := ""

	emailTo := "To: " + email.To + "\r\n"
	subject := "Subject: " + email.Subject + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	if len(email.Cc) > 0 {
		cc = "cc: " + strings.Join(email.Cc, ", ") + "\n"
	}
	msg := []byte(emailTo + cc + subject + mime + "\n" + email.Body)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	// Send the message
	_, err = GmailService.Users.Messages.Send("me", &message).Do()

	if err != nil {
		return false, err
	}

	return true, nil
}

func SendEmailWithAttachmentOAUTH2(email Email, fileDir string, fileName string) (bool, error) {
	var message gmail.Message

	fileBytes, err := ioutil.ReadFile(fileDir + fileName)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fileMIMEType := http.DetectContentType(fileBytes)

	fileData := base64.StdEncoding.EncodeToString(fileBytes)

	boundary := randStr(32, "alphanum")

	messageBody := []byte("Content-Type: multipart/mixed; boundary=" + boundary + " \n" +
		"MIME-Version: 1.0\n" +
		"to: " + email.To + "\n" +
		"subject: " + email.Subject + "\n" +
		"cc: " + strings.Join(email.Cc, ", ") + "\n\n" +

		"--" + boundary + "\n" +
		"Content-Type: text/plain; charset=" + string('"') + "UTF-8" + string('"') + "\n" +
		"MIME-Version: 1.0\n" +
		"Content-Transfer-Encoding: 7bit\n\n" +
		email.Body + "\n\n" +
		"--" + boundary + "\n" +

		"Content-Type: " + fileMIMEType + "; name=" + string('"') + fileName + string('"') + " \n" +
		"MIME-Version: 1.0\n" +
		"Content-Transfer-Encoding: base64\n" +
		"Content-Disposition: attachment; filename=" + string('"') + fileName + string('"') + " \n\n" +
		chunkSplit(fileData, 76, "\n") +
		"--" + boundary + "--")

	message.Raw = base64.URLEncoding.EncodeToString(messageBody)

	// Send the message
	_, err = GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Printf("Error: %v", err)
	}
	return true, nil
}
func randStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	var strBytes = make([]byte, strSize)
	_, _ = rand.Read(strBytes)
	for k, v := range strBytes {
		strBytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(strBytes)
}

func chunkSplit(body string, limit int, end string) string {
	var charSlice []rune

	// push characters to slice
	for _, char := range body {
		charSlice = append(charSlice, char)
	}

	var result = ""

	for len(charSlice) >= 1 {
		// convert slice/array back to string
		// but insert end at specified limit
		result = result + string(charSlice[:limit]) + end

		// discard the elements that were copied over to result
		charSlice = charSlice[limit:]

		// change the limit
		// to cater for the last few words in
		if len(charSlice) < limit {
			limit = len(charSlice)
		}
	}
	return result
}

// func parseTemplate(templateFileName string, data interface{}) (string, error) {
// 	templatePath, err := filepath.Abs(fmt.Sprintf("gomail/email_templates/%s", templateFileName))
// 	if err != nil {
// 		return "", errors.New("invalid template name")
// 	}
// 	t, err := template.ParseFiles(templatePath)
// 	if err != nil {
// 		return "", err
// 	}
// 	buf := new(bytes.Buffer)
// 	if err = t.Execute(buf, data); err != nil {
// 		return "", err
// 	}
// 	body := buf.String()
// 	return body, nil
// }
