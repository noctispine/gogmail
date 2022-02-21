/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/noctispine/go-email-app/db"
	"github.com/spf13/cobra"
)

// add emails from command line args (when flag secure == false)
func addEmailsFromArgs(args []string) error {
	for i := 0; i < len(args); i += 2 {
		userEmail := db.UserEmail{
			Email:    args[i],
			Password: args[i+1],
		}

		err := db.AddUserEmail(userEmail)
		if err != nil {
			return err
		}
	}

	return nil
}

// add emails when flag securce is enabled
// it conceales email's password whiles user writing
func addEmailsWithSecureFlag() error {
	var userEmail db.UserEmail
	fmt.Printf("Email:")
	fmt.Scanln(&userEmail.Email)

	fmt.Printf("Password:")
	color.Set(color.Concealed)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	line := scanner.Text()
	userEmail.Password = line

	userEmail.Password = line
	color.Unset()

	err := db.AddUserEmail(userEmail)
	if err != nil {
		return err
	}

	return nil
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [email] [password]",
	Short: "Add email and password",
	Long: `Add emails and passwords. It is possible to
	add multiple email and password combinations at once
	eg: add email1 password1 email2 password2...
	Also it is possible to conceal password with enabling
	secure flag (but user can add only one email).`,
	Args: func(cmd *cobra.Command, args []string) error {
		sec, err := cmd.Flags().GetBool("secure")
		if err != nil {
			return errors.New("flag secure cannot parsed")
		}

		if (!sec && len(args) != 0) && (len(args)%2 == 1) {
			return errors.New("arguments length should be even number and not equal to 0")
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

		color.Green("Emails Successfully Added")

	},
}

func init() {
	userCmd.AddCommand(addCmd)
	addCmd.Flags().BoolP("secure", "s", false, "hide your password when adding")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
