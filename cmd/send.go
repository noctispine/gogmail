/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/manifoldco/promptui"
	helper "github.com/noctispine/gogmail/cmd/helpers"
	"github.com/noctispine/gogmail/db"
	"github.com/noctispine/gogmail/gservice"
	"github.com/spf13/cobra"
)

// if there is authenticated user first try to email
// if token has been expired then try to relogin
// after try to send email again
func sendEmail(email gservice.Email, userInfos gservice.OAuthInfos) error {
	var err error

	if gservice.GmailService != nil {
		_, err = gservice.SendEmailOAUTH2(email)
		if err != nil {
			gservice.OAuthGmailService(userInfos)
			_, err = gservice.SendEmailOAUTH2(email)
			if err != nil {
				return err
			}
		}
	} else {
		gservice.OAuthGmailService(userInfos)
		_, err = gservice.SendEmailOAUTH2(email)
		if err != nil {
			return err
		}

	}

	return nil
}

func getEmailBody() (string, error) {
	var text string
	scn := bufio.NewScanner(os.Stdin)
	fmt.Println("Body: (press q in empty line to exit)")
	var lines []string
	for scn.Scan() {
		line := scn.Text()
		if len(line) == 1 {
			// Group Separator (GS ^]): ctrl-]
			// control-enter
			if line[0] == 'q' {
				break
			}
		}
		lines = append(lines, line)
	}

	if len(lines) > 0 {
		for i, line := range lines {
			if i != 0 {
				text += "\n"
			}
			text += line
		}
	}

	if err := scn.Err(); err != nil {
		return "", err
	}
	// if len(lines) == 0 {
	// 	break
	// }

	return text, nil
}

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send email with using Gmail API",
	Long:  `Send email with using Gmail API`,
	Run: func(cmd *cobra.Command, args []string) {
		var users []db.User
		var err error
		var to, subject, body string
		var doesSend string
		var size int
		var textFilePath string
		var showCredentials bool

		showCredentials, err = cmd.Flags().GetBool("credentials")
		if err != nil {
			log.Fatal(err)
		}

		size, err = cmd.Flags().GetInt("size")
		if err != nil {
			log.Fatal(err)
		}

		textFilePath, err = cmd.Flags().GetString("file")
		if err != nil {
			log.Fatal(err)
		}

		users, err = db.MakeSliceFromUser()
		if err != nil {
			log.Fatal(err)
		}

		helper.AddQuitOptionToEmailSlice(users)

		prompt := helper.SelectUser(users, size, showCredentials)
		var userIndex int

		// set prompt options

		validate := func(input string) error {
			if len(input) < 1 {
				return errors.New("Must have more than 0 characters")
			}
			return nil
		}

		promptTo := promptui.Prompt{
			Label:    "To",
			Validate: validate,
		}

		promptSubject := promptui.Prompt{
			Label: "Subject",
		}

		// promptBody := promptui.Prompt{
		// 	Label:     "Body (VIM MODE)",
		// 	IsVimMode: true,
		// }

		for {
			userIndex, _, err = prompt.Run()
			if userIndex == 0 {
				return
			}

			to, err = helper.PromptField(promptTo)
			if err != nil {
				log.Fatal(err)
			}

			subject, err = helper.PromptField(promptSubject)
			if err != nil {
				log.Fatal(err)
			}

			// if textFilePath is defined, email body should be text file's content
			// otherwise take body via stdin
			if textFilePath != "" {
				var b []byte
				b, err = ioutil.ReadFile(textFilePath)
				if err != nil {
					log.Fatal(err)
				}
				body = string(b)

			} else {
				// body, err = helper.PromptField(promptBody)
				body, err = getEmailBody()
				if err != nil {
					log.Fatal(err)
				}

			}

			doesSend, err = helper.PromptConfirm("are you sure to send this email?")

			if doesSend == "y" {
				var email gservice.Email

				email.To = to
				email.Subject = subject
				email.Body = body
				fmt.Println(body)

				sendEmail(email, users[userIndex].Infos)

			}
		}

		// fmt.Println("send called")
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	// add size, showCredentials, text_file flags
	sendCmd.Flags().StringP("file", "f", "", "text file for the body of the email")
	sendCmd.Flags().BoolP("credentials", "c", false, "show credentials")
	sendCmd.Flags().IntP("size", "s", 5, "prompt email list size")

}
