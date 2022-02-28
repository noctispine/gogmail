package helper

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/noctispine/gogmail/db"
)

type PromptUserDetails struct {
	EmailAddress string
	ClientID     string
}

// func ConstructEmail() (gservice.Email, error) {
// 	var err error
// 	var to, subject, body string
// 	var email gservice.Email

// 	// to, err = PromptField("To", 3)
// 	if err != nil {
// 		return email, err
// 	}

// 	// subject, err = PromptField("subject", 0)
// 	if err != nil {
// 		return email, err
// 	}

// 	// body, err = PromptField("body", 0)
// 	if err != nil {
// 		return email, err
// 	}

// 	email.To = to
// 	email.Subject = subject
// 	email.Body = body

// 	return email, nil
// }

// isDirectory determines if a file represented
// by `path` is a directory or not
func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func getWorkingDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return wd, nil
}

func PromptSelectDir() (string, error) {
	var files []string
	var err error
	var prompt promptui.Select
	var isDir bool
	var path string

	path, err = getWorkingDir()
	if err != nil {
		return "", err
	}

	for {
		files, err = MakeSliceFromDir(path)

		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "\U000027A4  ({{ . | green }})",
			Inactive: "  ({{ . | cyan }}) ",
		}

		prompt = promptui.Select{
			Label:        "Files",
			Items:        files,
			Size:         10,
			HideSelected: true,
			Templates:    templates,
		}
		if err != nil {
			return path, err
		}
		var i int
		i, _, err = prompt.Run()
		if err != nil {
			return path, err
		}

		if i == 0 {
			return "", nil
		} else if i == 1 {
			index := strings.LastIndex(path, "/")
			path = path[:index]
		} else {
			path = path + "/" + files[i]
			isDir, err = isDirectory(path)
			if err != nil {
				return path, err
			}
			if !isDir {
				return path, nil
			}
		}
	}
}

func SelectUser(emails []db.User, size int, showPw bool, cursorPos_optional ...int) promptui.Select {
	cursorPos := 0

	if len(cursorPos_optional) > 0 {
		cursorPos = cursorPos_optional[0]
	}

	searcher := func(input string, index int) bool {
		email := emails[index]
		emailAddress := strings.Replace(strings.ToLower(email.EmailAddress), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(emailAddress, input)
	}

	// set custom template's active and inactive field
	// accordingly flag password (which enables to print passwords)
	// when user want to list
	var activeField, inActiveField string
	if showPw {
		activeField = "\U000027A4  {{ .EmailAddress | cyan }} ({{ .Infos.ClientID | red }})"
		inActiveField = "  {{ .EmailAddress | cyan }} ({{ .Infos.ClientID | red }})"
	} else {
		activeField = "\U000027A4  ({{ .EmailAddress | cyan }})"
		inActiveField = "  ({{ .EmailAddress | cyan }}) "
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   activeField,
		Inactive: inActiveField,
	}

	prompt := promptui.Select{
		Label:        "Emails",
		Items:        emails,
		Searcher:     searcher,
		Templates:    templates,
		Size:         size,
		CursorPos:    cursorPos,
		HideSelected: true,
	}

	return prompt
}

func AddQuitOptionToEmailSlice(emails []db.User) {
	quitOpt := db.User{
		EmailAddress: "Quit",
	}

	emails = append(emails, quitOpt)

	replaceFirstAndLastElementOfSlice(emails)

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

// prompt wrt given options
func PromptField(prompt promptui.Prompt) (string, error) {

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	return result, nil
}

// it returns y if it is ok otherwise n
func PromptConfirm(label string) (string, error) {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	return result, err
}

func replaceFirstAndLastElementOfSlice(data []db.User) {
	temp := data[0]
	data[0] = data[len(data)-1]
	data[len(data)-1] = temp

}
