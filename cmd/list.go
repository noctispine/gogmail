/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/noctispine/gogmail/db"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all registered emails",
	Long:  `list all registered emails`,
	Run: func(cmd *cobra.Command, args []string) {
		var emails []db.User
		var err error
		var showC bool
		showC, err = cmd.Flags().GetBool("credentials")
		if err != nil {
			log.Fatal(err)
		}

		emails, err = db.MakeSliceFromUser()
		if err != nil {
			log.Fatal("Error occured while getting emails from the bucket", err)
		}

		for i, email := range emails {
			// The first index in email slice is empty
			// for QUIT option so we can skip 0 index
			if i != 0 {
				gr := color.New(color.FgGreen)
				re := color.New(color.FgRed)
				if i%2 == 1 {
					gr.Printf("%d: %s", i, email.EmailAddress)

				} else {
					fmt.Printf("%d: %s", i, email.EmailAddress)
				}

				if showC {
					re.Printf("%20s %s %s %s", email.Infos.ClientID, email.Infos.ClientSecret, email.Infos.RefreshToken, email.Infos.AccessToken)
				}

				fmt.Println()
			}

		}

	},
}

func init() {
	userCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("credentials", "c", false, "list emails with credentials")
	listCmd.SetUsageTemplate(rootCmd.Name() + "Usage: user list [flags]\n" +
		"\nFlags:\n" + listCmd.Flags().FlagUsages())

}
