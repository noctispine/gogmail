/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/noctispine/go-email-app/db"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all registered emails",
	Long:  `list all registered emails`,
	Run: func(cmd *cobra.Command, args []string) {
		var emails []db.UserEmail
		var err error
		var showPw bool
		showPw, err = cmd.Flags().GetBool("password")
		if err != nil {
			log.Fatal(err)
		}

		emails, err = db.MakeSliceFromEmailBucket()
		if err != nil {
			log.Fatal("Error occured while getting emails from the bucket", err)
		}

		for i, email := range emails {
			gr := color.New(color.FgGreen)
			re := color.New(color.FgRed)
			if i%2 == 1 {
				gr.Printf("%d: %s", i, email.Email)

			} else {
				fmt.Printf("%d: %s", i, email.Email)
			}

			if showPw {
				re.Printf("%20s", email.Password)
			}

			fmt.Println()
		}

	},
}

func init() {
	userCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("password", "p", false, "list emails with passwords")
	listCmd.SetUsageTemplate(rootCmd.Name() + "Usage: user list [flags]\n" +
		"\nFlags:\n" + listCmd.Flags().FlagUsages())

}
