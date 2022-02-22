package helper

import (
	"errors"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/noctispine/go-email-app/db"
)

func SelectEmail(emails []db.UserEmail, size int, showPw bool, cursorPos_optional ...int) promptui.Select {
	cursorPos := 0

	if len(cursorPos_optional) > 0 {
		cursorPos = cursorPos_optional[0]
	}

	searcher := func(input string, index int) bool {
		email := emails[index]
		emailAddress := strings.Replace(strings.ToLower(email.Email), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(emailAddress, input)
	}

	// set custom template's active and inactive field
	// accordingly flag password (which enables to print passwords)
	// when user want to list
	var activeField, inActiveField string
	if showPw {
		activeField = "\U000027A4  {{ .Email | cyan }} ({{ .Password | red }})"
		inActiveField = "  {{ .Email | cyan }} ({{ .Password | red }})"
	} else {
		activeField = "\U000027A4  ({{ .Email | cyan }})"
		inActiveField = "  ({{ .Email | cyan }}) "
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   activeField,
		Inactive: inActiveField,
	}

	prompt := promptui.Select{
		Label:     "Emails",
		Items:     emails,
		Searcher:  searcher,
		Templates: templates,
		Size:      size,
		CursorPos: cursorPos,
	}

	return prompt
}

func AddQuitOptionToEmailSlice(emails []db.UserEmail) {
	quitOpt := db.UserEmail{
		Email:    "Quit",
		Password: "Quit",
	}

	emails = append(emails, quitOpt)
	replaceFirstAndLastElementOfEmailSlice(emails)

}

// get email with using prompt
func PromptEmail() (string, error) {
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

// get password with using prompt
func PromptPassword(minumumLength int) (string, error) {
	validate := func(input string) error {
		if len(input) < minumumLength {
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

func PromptConfirm(label string) (string, error) {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	return result, err
}

func replaceFirstAndLastElementOfEmailSlice(data []db.UserEmail) {
	temp := data[0]
	data[0] = data[len(data)-1]
	data[len(data)-1] = temp

}
