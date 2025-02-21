package main

import (
	"context"
	"fmt"
	"log"
	"net/mail"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/fatih/color"
)

var (
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

type GitConfig struct {
	Name  string
	Email string
}

type FormValues struct {
	homeDir        string
	osName         string
	packageManager string
	gitConfig      GitConfig
	run            bool
}

var formValues FormValues

func main() {
	detectOS()
	detectUserInfo()
	printDreitagebart()
	runQuestionnaire()
	runInstallation()
}

func runInstallation() {
	installPackage("zsh")
	installPackage("stow")
	installPackage("neovim", "nvim")
	// installPackage("tmux")
	// installPackage("fzf")
	// installPackage("ripgrep")
	// installPackage("bat")
	// installPackage("eza")
	// installPackage("zoxide")
	// installPackage("thefuck")
}

func runQuestionnaire() {
	formValues.run = true

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("I detected "+formValues.osName+" as your linux distribution - so I will use "+formValues.packageManager+" for installing your software packages. Is this okay?").
				Options(
					huh.NewOption("apt", "apt"),
					huh.NewOption("dnf", "dnf"),
					huh.NewOption("pacman", "pacman"),
				).
				Value(&formValues.packageManager),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("What's your name and email?\nI need this information for your .gitconfig file").
				Value(&formValues.gitConfig.Name),
			huh.NewInput().
				Title("This email will be used for git commits").
				Value(&formValues.gitConfig.Email),
		),
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

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)

	return err == nil
}

func isPackageInstalled(packageName string) bool {
	command := exec.Command("which", packageName)

	_, err := command.CombinedOutput()

	if err != nil {
		return false
	}

	return true
}

func printDreitagebart() {
	fmt.Println(`
  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▒▓▓▓
  ▓▓▓▓▓▓░░░░░░░░░░▒▓▓▓▓▓▓▓▓▓░░░░░░▓▓
  ▓▓▓▓███▓▓▒▒▒▒▒▓████▓▓▓▓▓▓▓█████▓░░░
  ▓▓▓▒░░░▒██████▓░░░▒▓▓▓▓▒░░░░░▒██▒░░░
  ▓█▓▒ ░░░██▒▒██▒░░ ░▓█▓░ ░▒░░░░▓▒░░░░
  ░░░░░░░░█░░░░█▓░░░░░░░░░░░░░░░█░░▒░░
  ░░░░░░░▓▓░░░░░█░░░░░░░░░░░░░░▒▓░░▒░░
  ░░░░░░▓▓░░░░░░▒█░░░░░░░░░░░░░█░░▒▒░░
  ▓▓▓▓▓▒░░░░░░░░░░▒▒▓▓▓▓▓▓▓▓▓▓▒░░░▒▒░
  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▓▒░
  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▒▓▓
  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▓▓▓
  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▒▓▓▓
  ░░░░░░░░░▒▓▓▒░░░░░░░░░░░░░░░░░▓▓▓▓         _       _    __ _ _
  ░░░▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▒░░░░░░░░░░▓▓▓▓         | |     | |  / _(_) |
  ░▒▓▓▓▓▓▓▓▓▒▒▓▒▓▓▓▓▓▓▓░░░░░░░▒▓▓▓▓       __| | ___ | |_| |_ _| | ___  ___
  ▒▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▓▓░░░░▓▓▓▓▓▓       / _' |/ _ \| __|  _| | |/ _ \/ __|
  ▓▓▒░░░░░░░░░░░░░░░░░▓▓▒▒▓▓▓▓▓▓▓       | (_| | (_) | |_| | | | |  __/\__ \
  ▓▓▓▒░░░░▓▓▓▓▓▓▒░░░░▒▓▓▓▓▓▓▓▓▓▓▓      (_)__,_|\___/ \__|_| |_|_|\___||___/
  ▓▓▓▓░░░░░▒▓▓▓▒░░░░▒▓▓▓▓▓▓▓▓▓▓
  ▓▓▓▓▓▒░░░▒▒▒▒░░░░▒▓▓▓▓▓▓▓▓▓▓
  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
	`)
	fmt.Println("Hit ENTER to start the installer...")
	fmt.Scanln()
}

func detectUserInfo() {
	currentUser, err := user.Current()

	if err == nil {
		formValues.gitConfig.Name = currentUser.Username
		formValues.gitConfig.Email = fmt.Sprintf("%s@example.com", currentUser.Username)

		formValues.homeDir = currentUser.HomeDir
	}
}

func detectOS() {
	_, err := os.Stat("/etc/os-release")

	if err == nil {
		content, err := os.ReadFile("/etc/os-release")

		if err != nil {
			return
		}

		lines := strings.Split(string(content), "\n")

		for _, line := range lines {
			if strings.HasPrefix(line, "ID=") {
				osID := strings.Trim(strings.TrimPrefix(line, "ID="), "\"")

				formValues.osName = osID

				switch osID {
				case "arch", "manjaro":
					formValues.packageManager = "pacman"
				case "debian", "ubuntu", "raspbian", "linuxmint", "pop":
					formValues.packageManager = "apt"
					return
				case "fedora", "rocky", "almalinux":
					formValues.packageManager = "dnf"
					return
				}
			}
		}
	}
}

func installNixInstaller() {
	var command *exec.Cmd

	command = exec.Command("sh <(curl -L https://nixos.org/nix/install) --daemon")

	err := spinner.New().Type(spinner.MiniDot).ActionWithErr(func(context.Context) error {
		_, err := command.CombinedOutput()

		return err
	}).Run()

	if err != nil {
		fmt.Println(red(fmt.Sprintf("Failed to install nix installer: %v", err)))
		os.Exit(1)
	}
}

func installHomebrew() {
	var command *exec.Cmd

	command = exec.Command("/bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"")

	err := spinner.New().Type(spinner.MiniDot).ActionWithErr(func(context.Context) error {
		_, err := command.CombinedOutput()

		return err
	}).Run()

	if err != nil {
		fmt.Println(red(fmt.Sprintf("Failed to install homebrew: %v", err)))
		os.Exit(1)
	}
}

func installPackage(packageName string, alias ...string) {
	var pkg string
	var command *exec.Cmd

	if len(alias) > 0 {
		pkg = alias[0]
	} else {
		pkg = packageName
	}

	if isPackageInstalled(pkg) {
		fmt.Println(green(pkg) + " is already installed... skipped")

		return
	}

	switch formValues.packageManager {
	case "apt":
		command = exec.Command("sudo", "apt", "install", "-y", packageName)
	case "dnf":
		command = exec.Command("sudo", "dnf", "install", "-y", packageName)
	case "pacman":
		command = exec.Command("sudo", "pacman", "-Suy", packageName)
	}

	err := spinner.New().Type(spinner.MiniDot).
		Title(" Installing package...").ActionWithErr(func(context.Context) error {
		_, err := command.CombinedOutput()

		return err
	}).Run()

	if err != nil {
		fmt.Println(red(fmt.Sprintf("Failed to install %s: %v", packageName, err)))

		os.Exit(1)
	}

	fmt.Println(green(fmt.Sprintf("%s installed successfully.", packageName)))
}
