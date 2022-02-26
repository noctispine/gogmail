package gservice

// GmailService : Gmail client for sending email

import (
	"context"
	"encoding/base64"
	"log"
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

	emailTo := "To: " + email.To + "\r\n"
	subject := "Subject: " + email.Subject + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + subject + mime + "\n" + email.Body)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	// Send the message
	_, err = GmailService.Users.Messages.Send("me", &message).Do()

	if err != nil {
		return false, err
	}

	return true, nil
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
