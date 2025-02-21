package main

import (
	"errors"
	"fmt"
	"net/mail"
	"os/user"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

var (
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

type SelectItem struct {
	Value string
	Label string
}

type GitConfig struct {
	Name  string
	Email string
}

func main() {
	printDreitagebart()
	runQuestionnaire()
}

func runQuestionnaire() {
	gitConfig := questionGitConfig()

	useTmux := questionUseTmux()
	useNeoVim := questionUseNeoVim()

	// if gitConfig {
	fmt.Println("Your name is " + gitConfig.Name)
	fmt.Println("Your email is " + gitConfig.Email)
	// }

	if useTmux {
		fmt.Println("You want to use tmux")
	}

	if useNeoVim {
		fmt.Println("You want to use neovim")
	}

}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)

	return err == nil
}

func questionGitConfig() GitConfig {
	var result string
	var gitConfig GitConfig

	validateName := func(input string) error {
		if len(input) < 3 {
			return errors.New("Username must have more than 3 characters")
		}

		return nil
	}

	userInfo, err := user.Current()

	if err == nil {
		gitConfig.Name = userInfo.Username
	}

	namePrompt := promptui.Prompt{
		Label:     "What is your name in .gitconfig?",
		Validate:  validateName,
		Default:   gitConfig.Name,
		AllowEdit: true,
	}

	for {
		result, err = namePrompt.Run()

		if err == nil {
			gitConfig.Name = result

			break
		}
	}

	validateEmail := func(input string) error {
		if len(input) == 0 {
			return errors.New("This is not a valid email address")
		}

		if !isValidEmail(input) {
			return errors.New("This is not a valid email address")
		}

		return nil
	}

	emailPrompt := promptui.Prompt{
		AllowEdit: true,
		Label:     "What is your email in .gitconfig",
		Validate:  validateEmail,
		Default:   userInfo.Username + "@example.com",
	}

	for {
		result, err = emailPrompt.Run()

		if err == nil {
			gitConfig.Email = result

			break
		}
	}

	return gitConfig
}

func questionUseNeoVim() bool {
	prompt := promptui.Prompt{
		Label:     "Do you want to use neovim?",
		IsConfirm: true,
	}

	_, err := prompt.Run()

	if err != nil {
		return false
	}

	return true
}

func questionUseTmux() bool {
	prompt := promptui.Prompt{
		Label:     "Do you want to use tmux?",
		IsConfirm: true,
	}

	_, err := prompt.Run()

	if err != nil {
		return false
	}

	// fmt.Printf("You choose %q\n", result)

	return true
}

func printDreitagebart() {
	fmt.Println(`

                                ▓▓▓▓▓▓▓▓▓▓▓▓▓█
                             ▓▓▓▓▓▓▓▓█▓▓▓▓▓▓▓▓▓▓▓▓
                         ▓▓▓▓▓▓▓▓▓▓▓▓█▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
                        ▓▓▓▓▓▓▓▓▓▓▓▓▓▓█▓▓▓▓▓▓▓▓▓▓▓▓█▓▓▓█
                     █▓▓▓▓▓▓▓▓▓▓▓▓▓▓█████▓▓▓▓▓▓▓▓▓▓▓▓▓█▓▓
                    █▓▓▓▓▓▓▓▓▓▓▓▓██▓▓████████▓▓▓▓▓▓▓▓▓▓█▓▓▓
                   ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓███▓▓▓▓▓▓▓▓▓█▓▓
                  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓█████████▓▓▓▓▓███▓██▓▓▓▓▓▓
                  ▓▓▓▓▓▓█▓▓▓▓▓█████████████████▓▓▓▓████████▓▓
                 ▓▓█████▓▓▓███▒░░░░░░░░░░░░░░░░▓██▓▓▓▓██████▓
                 ▓████▓▓▓██▒░░░░░░░░░░░░░░░░░░░░░░░▓█▓▓█████▓▓
                ▓████▓▓█▒░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▓▓▓████▓
                ▓███▓▓▒░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▓▓▓██▓
                ▓██▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▒▓▓▓
                █▓░░░░░░▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░▒▓▓▓▓▓▓▓▓▓░░░░░░▓▓
                ░░▓█████▓▓▓▓▓▓▓▓███▓▓▒▒▒▒▒▓████▓▓▓▓▓▓▓█████▓░░░
               ░░░██▓░░░░░░▓▓▓▓▒░░░▒██████▓░░░▒▓▓▓▓▒░░░░░▒██▒░░░
               ░░░░█▒░░░░▒ ░▓█▓▒ ░░░██▒▒██▒░░ ░▓█▓░ ░▒░░░░▓▒░░░░
               ░░▒░▓▒░░░░░░░░░░░░░░░█░░░░█▓░░░░░░░░░░░░░░░█░░▒░░
               ░░▓░▒█░░░░░░░░░░░░░░▓▓░░░░░█░░░░░░░░░░░░░░▒▓░░▒░░
               ░░▓▒░▓▓░░░░░░░░░░░░▓▓░░░░░░▒█░░░░░░░░░░░░░█░░▒▒░░
                ░▓▒░░░▒▓▓▓▓▓▓▓▓▓▓▒░░░░░░░░░░▒▒▓▓▓▓▓▓▓▓▓▓▒░░░▒▒░
                ░▓▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▓▒░
                ▒▓▓▒░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▒▓▓
                 ▓▓▒░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▓▓▓
                 ▓▓▓▒░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▒▓▓▓
                 ▓▓▓▓░░░░░░░░░░░░░░░░▒▓▓▒░░░░░░░░░░░░░░░░░▓▓▓▓
                  ▓▓▓▒░░░░░░░░░▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▒░░░░░░░░░░▓▓▓▓
                  ▓▓▓▓▒░░░░░░▒▓▓▓▓▓▓▓▓▒▒▓▒▓▓▓▓▓▓▓░░░░░░░▒▓▓▓▓
                   ▓▓▓▓▓▒░░░▒▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▓▓░░░░▓▓▓▓▓▓
                   ▓▓▓▓▓▓▓▓▓▓▓▒░░░░░░░░░░░░░░░░░▓▓▒▒▓▓▓▓▓▓▓
                    ▓▓▓▓▓▓▓▓▓▓▓▒░░░░▓▓▓▓▓▓▒░░░░▒▓▓▓▓▓▓▓▓▓▓▓
                      ▓▓▓▓▓▓▓▓▓▓░░░░░▒▓▓▓▒░░░░▒▓▓▓▓▓▓▓▓▓▓
                       ▓▓▓▓▓▓▓▓▓▓▒░░░▒▒▒▒░░░░▒▓▓▓▓▓▓▓▓▓▓
                         ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
                          ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
                             ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
                                ▓▓▓▓▓▓▓▓▓▓▓▓▓▓
                                    ▓▓▓▓▓▓▓

	`)
	fmt.Println("Hit ENTER to start the installer...")
	s := spinner.New(spinner.CharSets[27], 100*time.Millisecond)
	s.Start()
	time.Sleep(4 * time.Second)
	s.Stop()
	fmt.Scanln()
}
