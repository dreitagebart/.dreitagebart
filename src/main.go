package main

import (
	"log"
	"os"
	"path"

	"github.com/charmbracelet/huh"
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

	installHomebrewFont("font-meslo-lg-nerd-font")
	installHomebrewPackage("lazygit")
	installHomebrewPackage("fzf")
	installHomebrewPackage("ripgrep", "rg")
	installHomebrewPackage("bat")
	installHomebrewPackage("eza")
	installHomebrewPackage("zoxide")
	installHomebrewPackage("thefuck")

	copyTemplateFiles(templateFiles, path.Join(formValues.homeDir, ".dreitagebart"))
	configureGit()
	createBackupPath()

	stowFile(".gitconfig", "git")
	stowFile(".zshrc", "zsh")
	stowFile(".p10k.zsh", "p10k")

	if formValues.installNeovimAddons {
		stowFile(".config/nvim", "nvim")
	}

	if formValues.installTmuxAddons {
		stowFile(".tmux.conf", "tmux")
	}

	setDefaultShell()
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
