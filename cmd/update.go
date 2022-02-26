/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
	helper "github.com/noctispine/gogmail/cmd/helpers"
	"github.com/noctispine/gogmail/db"
	"github.com/noctispine/gogmail/gservice"
	"github.com/spf13/cobra"
)

// update password using command line args
func updateFromArgs(args []string) error {
	var err error
	for i := 0; i < len(args); i += 5 {
		infos := gservice.OAuthInfos{
			ClientID:     args[i+1],
			ClientSecret: args[i+2],
			RefreshToken: args[i+3],
			AccessToken:  args[i+4],
		}

		err = db.ChangeEmailInfos(args[0], infos)
		if err != nil {
			return err
		}
	}

	return nil
}

// prompt emails, let the user select which user(s) update
// then update selected user
func promptAndUpdate(showCredentials bool, size int) error {
	var err error
	var users []db.User
	var quit bool
	var isOk string
	// get emails from db and add quit option
	users, err = db.MakeSliceFromUser()
	helper.AddQuitOptionToEmailSlice(users)

	if err != nil {
		return err
	}

	quit = false
	i := 0
	for !quit {
		var infos gservice.OAuthInfos
		var clientID, clientSecret, refreshToken, AccessToken string
		promptSelect := helper.SelectUser(users, size, showCredentials)

		// set prompt options
		promptClientID := promptui.Prompt{
			Label: "Client ID",
		}

		promptClientSecret := promptui.Prompt{
			Label:       "Client Secret",
			Mask:        '*',
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

		i, _, err = promptSelect.Run()
		if err != nil {
			return err
		}

		// if index == 0, quit
		if i == 0 {
			return nil
		}

		clientID, err = helper.PromptField(promptClientID)
		if err != nil {
			return err
		}

		clientSecret, err = helper.PromptField(promptClientSecret)
		if err != nil {
			return err
		}

		refreshToken, err = helper.PromptField(promptRefreshToken)
		if err != nil {
			return err
		}

		AccessToken, err = helper.PromptField(promptAccessToken)
		if err != nil {
			return err
		}

		// if user changed something mutate otherwise keep the
		// original records
		if len(clientID) > 0 {
			infos.ClientID = clientID
		} else {
			infos.ClientID = users[i].Infos.ClientID
		}
		if len(clientSecret) > 0 {
			infos.ClientSecret = clientSecret
		} else {
			infos.ClientSecret = users[i].Infos.ClientSecret
		}
		if len(refreshToken) > 0 {
			infos.RefreshToken = refreshToken
		} else {
			infos.RefreshToken = users[i].Infos.RefreshToken
		}
		if len(AccessToken) > 0 {
			infos.AccessToken = AccessToken
		} else {
			infos.AccessToken = users[i].Infos.AccessToken
		}

		label := fmt.Sprintf("Are you sure you want to change credentials for: %s", users[i].EmailAddress)

		// confirm password change
		isOk, err = helper.PromptConfirm(label)
		if isOk == "y" {
			err = db.ChangeEmailInfos(users[i].EmailAddress, infos)
			if err != nil {
				return err
			}
			users[i].Infos = infos
		}

	}

	return nil
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update gmail oauth credentials",
	Long:  `update gmail oauth credentials`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var list, showC bool
		var size int

		list, err = cmd.Flags().GetBool("list")
		if err != nil {
			log.Fatal(err)
		}

		showC, err = cmd.Flags().GetBool("credentials")
		if err != nil {
			log.Fatal(err)
		}

		size, err = cmd.Flags().GetInt("size")
		if err != nil {
			log.Fatal(err)
		}

		if list {
			// prompt a list of emails to select and update
			err = promptAndUpdate(showC, size)
		} else {
			// update using command line args directly
			err = updateFromArgs(args)
			if err != nil {
				log.Fatal(err)
			}
		}

	},
}

func init() {
	userCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolP("list", "l", false, "prompt all registered emails")

	updateCmd.Flags().BoolP("credentials", "c", false, "list emails with credentials, it must be used with list flag")

	updateCmd.Flags().IntP("size", "s", 5, "set list size, flag list must be enabled")
	updateCmd.SetUsageTemplate(rootCmd.Name() + "Usage: user update [email1] [newPassword1] [email2] [newPassword2]... or user update [flags]\n" +
		"\nFlags:\n" + updateCmd.Flags().FlagUsages())

}
