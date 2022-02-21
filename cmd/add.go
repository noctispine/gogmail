/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/noctispine/go-email-app/db"
	"github.com/spf13/cobra"
)

// get user email with using prompt
func promptEmail() (string, error) {
	validate := func(input string) error {
		if len(input) < 3 {
			return errors.New("Email length must be at least 3")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Email",
		Validate: validate,
	}

	email, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return email, nil
}

// get user password with using prompt
func promptPassword() (string, error) {
	validate := func(input string) error {
		if len(input) < 1 {
			return errors.New("Password must have more than 0 characters")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Password",
		Validate: validate,
		Mask:     '*',
	}

	pw, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	return pw, nil
}

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
	var err error

	userEmail.Email, err = promptEmail()
	if err != nil {
		return err
	}

	userEmail.Password, err = promptPassword()
	if err != nil {
		return err
	}

	err = db.AddUserEmail(userEmail)
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
		if ((!sec && (len(args) == 0)) || !sec && (len(args) != 0)) && (len(args)%2 == 0) {
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
	addCmd.SetUsageTemplate(rootCmd.Name() + "Usage: user add [email1] [password1] [email2] [password2]...\n" +
		"\nFlags:\n" + addCmd.Flags().FlagUsages())
}
