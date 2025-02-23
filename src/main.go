package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

func main() {
	detectOS()
	detectUserInfo()
	printDreitagebart()
	runQuestionnaire()
	runInstallation()
}

func runInstallation() {
	installHomebrew()
	installNativePackage("zsh")
	installNativePackage("stow")
	if formValues.installNeovim {
		installNativePackage("neovim", "nvim")
	}
	if formValues.installTmux {
		installNativePackage("tmux")
	}

	installHomebrewPackage("fzf")
	installHomebrewPackage("ripgrep", "rg")
	installHomebrewPackage("bat")
	installHomebrewPackage("eza")
	installHomebrewPackage("zoxide")
	installHomebrewPackage("thefuck")

	copyTemplateFiles(templateFiles, path.Join(formValues.homeDir, ".dreitagebart"))

	stowFile(".gitconfig", "git")
}

func runQuestionnaire() {
	formValues.run = true
	formValues.installNeovimAddons = true
	formValues.installTmuxAddons = true
	formValues.installTmux = !isPackageInstalled("tmux")
	formValues.installNeovim = !isPackageInstalled("nvim")

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("It seems that you are using "+formValues.osName+" as your linux distribution"),
			huh.NewSelect[string]().
				Title("I will use "+formValues.packageManager+" for installing your software packages. Is this okay?").
				Options(
					huh.NewOption("apt", "apt"),
					huh.NewOption("dnf", "dnf"),
					huh.NewOption("pacman", "pacman"),
				).
				Value(&formValues.packageManager),
		),
		huh.NewGroup(
			huh.NewNote().
				Title("What's your name and email?\nI will put this information in your .gitconfig file"),
			huh.NewInput().
				Title("Your name").
				Value(&formValues.gitConfig.Name),
			huh.NewInput().
				Title("Your email").
				Value(&formValues.gitConfig.Email),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want to install Tmux?").
				Value(&formValues.installTmux),
		).WithHideFunc(func() bool {
			return !formValues.installTmux
		}),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Tmux is already installed - install tmux theme and plugins?").
				Value(&formValues.installTmuxAddons),
		).WithHideFunc(func() bool {
			return formValues.installTmux
		}),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want to install tmux themes and plugins as well?").
				Value(&formValues.installTmuxAddons),
		).WithHideFunc(func() bool {
			return !formValues.installTmux
		}),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want to install NeoVim?").
				Value(&formValues.installNeovim),
		).WithHideFunc(func() bool {
			return !formValues.installNeovim
		}),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Neovim is already installed - install neovim theme and plugins?").
				Value(&formValues.installNeovimAddons),
		).WithHideFunc(func() bool {
			return formValues.installNeovim
		}),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want to install neovim themes and plugins as well?").
				Value(&formValues.installNeovimAddons),
		).WithHideFunc(func() bool {
			return !formValues.installNeovim
		}),

		// huh.NewGroup(
		// 	huh.NewConfirm().
		// 		TitleFunc(func() string {
		// 			if isNeovimInstalled {
		// 				return "Neovim is already installed - install neovim theme and plugins?"
		// 			}

		// 			return "Do you want to install neovim themes and plugins as well?"
		// 		}, &isNeovimInstalled).
		// 		Value(&formValues.installNeovimAddons),
		// ).WithHideFunc(func() bool {
		// 	return !formValues.installNeovim
		// }),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Are you sure you want to run the installer?").
				Value(&formValues.run),
		),
	).WithTheme(huh.ThemeCatppuccin())

	err := form.Run()

	if err != nil {
		log.Fatal(err)
	}

	if !formValues.run {
		os.Exit(0)
	}
}

// func isValidEmail(email string) bool {
// 	_, err := mail.ParseAddress(email)

// 	return err == nil
// }

func stowFile(filename string, template string) {
	var command *exec.Cmd

	stowPath := path.Join(formValues.homeDir, filename)
	templatePath := path.Join(formValues.homeDir, ".dreitagebart", template)

	if fileExists(stowPath) {
		err := os.Rename(stowPath, path.Join(formValues.backupPath, filename))

		if err != nil {
			fmt.Printf("failed to move file: %s", err)
			os.Exit(1)
		}
	}

	command = exec.Command("stow", templatePath, "--dotfiles")

	err := spinner.New().Type(spinner.MiniDot).ActionWithErr(func(context.Context) error {
		_, err := command.CombinedOutput()

		return err
	}).Run()

	if err != nil {
		fmt.Println(red(fmt.Sprintf("Failed to stow %s", filename)))
		os.Exit(1)
	}
}
