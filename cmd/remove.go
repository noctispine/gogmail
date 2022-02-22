/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"log"

	helper "github.com/noctispine/go-email-app/cmd/helpers"
	"github.com/noctispine/go-email-app/db"
	"github.com/spf13/cobra"
)

// list

func listAndRemoveEmails(showPw bool, size int) error {
	// var all, tail, head bool
	var err error
	var emails []db.UserEmail

	emails, err = db.MakeSliceFromEmailBucket()
	if err != nil {
		return err
	}
	helper.AddQuitOptionToEmailSlice(emails)

	quit := false

	i := 0
	for !quit {
		selectPrompt := helper.SelectEmail(emails, size, showPw, i)

		i, _, err = selectPrompt.Run()
		if err != nil {
			return err
		}

		// if index == 0, quit
		// otherwise try to remove email
		if i == 0 {
			quit = true
		} else {
			err = db.RemoveUserEmail(emails[i].Email)
			if err != nil {
				return err
			}
			emails[i] = db.UserEmail{
				Email:    "DELETED",
				Password: "DELETED",
			}
		}
	}

	return nil

}

func removeEmailsFromArgs(args []string) error {
	for _, arg := range args {
		return db.RemoveUserEmail(arg)
	}

	return nil
}

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove email",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		// list, err := removeCmd.Flags().GetBool("list")
		// if err != nil {
		// 	return err
		// }

		list, err := cmd.Flags().GetBool("list")
		if err != nil {
			return err
		}

		if !list && (len(args) == 0) {
			return errors.New("Arguments length must not equal to 0")
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		var (
			size         int
			showPw, list bool
			err          error
		)

		// get flags: size, password, list
		size, err = cmd.Flags().GetInt("size")
		if err != nil {
			log.Fatal(err)
		}

		list, err = cmd.Flags().GetBool("list")
		if err != nil {
			log.Fatal(err)
		}

		showPw, err = cmd.Flags().GetBool("password")
		if err != nil {
			log.Fatal(err)
		}

		// if user enabled flag list, prompt a list to
		// allow to user select a email
		// otherwise if use directly wrote emails
		// to command line args, delete those directly
		if list {
			err = listAndRemoveEmails(showPw, size)
			if err != nil {
				log.Fatal(err)
			}

		} else {

			err = removeEmailsFromArgs(args)
			if err != nil {
				log.Fatal(err)
			}

		}
	},
}

func init() {
	userCmd.AddCommand(removeCmd)
	// removeCmd.Flags().Int64P("head", "e", 10, "list head, it should be used with list flag")
	// removeCmd.Flags().Int64P("tail", "t", 10, "list tail, it should be used with list flag")
	removeCmd.Flags().IntP("size", "s", 10, "page size, it must be used with list flag")
	removeCmd.Flags().BoolP("password", "p", false, "list emails with passwords, it must be used with list flag")
	removeCmd.Flags().BoolP("list", "l", false, "prompt option to select emails that will be removed")
	removeCmd.SetUsageTemplate(rootCmd.Name() + "Usage: user remove [email1] [email2] [email3...] or user remove [flags]\n" +
		"\nFlags:\n" + removeCmd.Flags().FlagUsages())

}
