/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"log"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	helper "github.com/noctispine/gogmail/cmd/helpers"
	"github.com/noctispine/gogmail/db"
	"github.com/noctispine/gogmail/gservice"
	"github.com/spf13/cobra"
)

// add emails from command line args (when flag secure == false)
func addEmailsFromArgs(args []string) error {
	for i := 0; i < len(args); i += 5 {
		infos := gservice.OAuthInfos{
			ClientID:     args[i+1],
			ClientSecret: args[i+2],
			RefreshToken: args[i+3],
			AccessToken:  args[i+4],
		}
		user := db.User{
			EmailAddress: args[i],
			Infos:        infos,
		}

		err := db.AddUser(user)
		if err != nil {
			return err
		}
	}

	return nil
}

// add emails when flag securce is enabled
// it conceales email's password whiles user writing
func addEmailsWithSecureFlag() error {
	var user db.User
	var infos gservice.OAuthInfos
	var err error

	// set prompt options
	validate := func(input string) error {
		if len(input) < 1 {
			return errors.New("Must have more than 0 characters")
		}
		return nil
	}

	promptClientID := promptui.Prompt{
		Label:    "Client ID",
		Validate: validate,
	}

	promptClientSecret := promptui.Prompt{
		Label:       "Client Secret",
		Mask:        '*',
		Validate:    validate,
		HideEntered: true,
	}

	promptRefreshToken := promptui.Prompt{
		Label:       "Refresh Token",
		Mask:        '*',
		HideEntered: true,
	}
	promptAccessToken := promptui.Prompt{
		Label:       "Access Token",
		Mask:        '*',
		HideEntered: true,
	}

	user.EmailAddress, err = helper.PromptEmail()
	if err != nil {
		return err
	}

	infos.ClientID, err = helper.PromptField(promptClientID)
	if err != nil {
		return err
	}

	infos.ClientSecret, err = helper.PromptField(promptClientSecret)
	if err != nil {
		return err
	}

	infos.RefreshToken, err = helper.PromptField(promptRefreshToken)
	if err != nil {
		return err
	}

	infos.AccessToken, err = helper.PromptField(promptAccessToken)
	if err != nil {
		return err
	}

	user.Infos = infos

	err = db.AddUser(user)
	if err != nil {
		return err
	}

	return nil
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [email] [client_id] [client_secret] [refresh_token] [access_token]... or add [flags]",
	Short: "Add email name and client credentials ",
	Long: `Add emails and client credentials. It is possible to
	add multiple combinations at once
	Also it is possible to conceal client_id with enabling
	secure flag (but user can add only one email).`,
	Args: func(cmd *cobra.Command, args []string) error {
		sec, err := cmd.Flags().GetBool("secure")
		log.Println(len(args))
		if err != nil {
			return errors.New("flag secure cannot parsed")
		}
		if !sec && len(args) == 0 {
			return errors.New("please provide email and password")
		} else if (!sec && (len(args) != 0)) && (len(args)%5 != 0) {
			return errors.New("arguments length should be even number")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		// check if user enabled secure flag
		sec, err := cmd.Flags().GetBool("secure")
		if err != nil {
			log.Fatal(err)
		}

		if sec {
			err = addEmailsWithSecureFlag()
		} else {
			err = addEmailsFromArgs(args)
		}

		if err != nil {
			log.Fatal("Error occured while trying to add email(s) to db")
		}

		color.Green("Email(s) Successfully Added")

	},
}

func init() {
	userCmd.AddCommand(addCmd)
	addCmd.Flags().BoolP("secure", "s", false, "hide your password when adding")
	addCmd.SetUsageTemplate(rootCmd.Name() + " add [email] [client_id] [client_secret] [refresh_token] [access_token]... or add [flags]\n" +
		"\nFlags:\n" + addCmd.Flags().FlagUsages())
}
