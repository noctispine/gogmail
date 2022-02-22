/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"

	helper "github.com/noctispine/go-email-app/cmd/helpers"
	"github.com/noctispine/go-email-app/db"
	"github.com/spf13/cobra"
)

// update password using command line args
func updateFromArgs(args []string) error {
	var err error
	for i := 0; i < len(args); i += 2 {
		newUserEmail := db.UserEmail{
			Email:    args[i],
			Password: args[i+1],
		}

		err = db.ChangeMailPassword(newUserEmail, newUserEmail.Password)
		if err != nil {
			return err
		}
	}

	return nil
}

// prompt emails, let the user select which user(s) update
// then update selected user
func promptAndUpdate(showPw bool, size int) error {
	var err error
	var emails []db.UserEmail
	var quit bool
	var pw string
	var isOk string
	// get emails from db and add quit option
	emails, err = db.MakeSliceFromEmailBucket()
	helper.AddQuitOptionToEmailSlice(emails)

	if err != nil {
		return err
	}

	quit = false
	i := 0
	for !quit {
		promptSelect := helper.SelectEmail(emails, size, showPw)
		i, _, err = promptSelect.Run()
		if err != nil {
			return err
		}

		// if index == 0, quit
		if i == 0 {
			return nil
		}

		// prompt for new password
		pw, err = helper.PromptPassword(1)
		if err != nil {
			return err
		}

		label := fmt.Sprintf("Are you sure you want to change your password: %s", emails[i].Email)

		// confirm password change
		isOk, err = helper.PromptConfirm(label)
		if isOk == "y" {
			err = db.ChangeMailPassword(emails[i], pw)
			if err != nil {
				return err
			}
			emails[i].Password = pw
		}

	}

	return nil
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update email address's password",
	Long:  `update email address's password`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var list, showPw bool
		var size int

		list, err = cmd.Flags().GetBool("list")
		if err != nil {
			log.Fatal(err)
		}

		showPw, err = cmd.Flags().GetBool("password")
		if err != nil {
			log.Fatal(err)
		}

		size, err = cmd.Flags().GetInt("size")
		if err != nil {
			log.Fatal(err)
		}

		if list {
			// prompt a list of emails to select and update
			err = promptAndUpdate(showPw, size)
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
	updateCmd.Flags().BoolP("password", "p", false, "show emails with passwords, flag list must be enabled")
	updateCmd.Flags().IntP("size", "s", 5, "set list size, flag list must be enabled")
	updateCmd.SetUsageTemplate(rootCmd.Name() + "Usage: user update [email1] [newPassword1] [email2] [newPassword2]... or user update [flags]\n" +
		"\nFlags:\n" + updateCmd.Flags().FlagUsages())

}
