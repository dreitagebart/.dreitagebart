package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
	"time"

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
	homeDir             string
	osName              string
	packageManager      string
	installTmux         bool
	installTmuxAddons   bool
	installNeovim       bool
	installNeovimAddons bool
	backupPath          string
	gitConfig           GitConfig
	run                 bool
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

func isPackageInstalled(packageName string) bool {
	command := exec.Command("which", packageName)

	_, err := command.CombinedOutput()

	return err == nil
}

func createBackupPath() {
	backupPath := path.Join(formValues.homeDir, ".dreitagebart", "backups")
	currentTime := time.Now()
	dateString := currentTime.Format("19851127-143015")

	formValues.backupPath = path.Join(backupPath, dateString)

	err := os.MkdirAll(formValues.backupPath, 0755)

	if err != nil {
		log.Fatal(err)
	}
}

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

	command = exec.Command("stow", templatePath)

	err := spinner.New().Type(spinner.MiniDot).ActionWithErr(func(context.Context) error {
		_, err := command.CombinedOutput()

		return err
	}).Run()

	if err != nil {
		fmt.Println(red(fmt.Sprintf("Failed to stow %s", filename)))
		os.Exit(1)
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)

	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	log.Fatal(err)

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
	// fmt.Scanln()
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

	if isPackageInstalled("brew") {
		fmt.Println(green("homebrew") + " is already installed... skipped")

		return
	}

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

func installHomebrewPackage(packageName string, alias ...string) {
	var pkg string
	var command *exec.Cmd

	if len(alias) > 0 {
		pkg = alias[0]
	} else {
		pkg = packageName
	}

	if isPackageInstalled(pkg) {
		fmt.Println(green(packageName) + " is already installed... skipped")

		return
	}

	command = exec.Command("brew", "install", pkg)

	err := spinner.New().Type(spinner.MiniDot).
		Title(" Installing package " + packageName + "...").ActionWithErr(func(context.Context) error {
		_, err := command.CombinedOutput()

		return err
	}).Run()

	if err != nil {
		fmt.Println(red(fmt.Sprintf("Failed to install %s: %v", packageName, err)))

		os.Exit(1)
	}

	fmt.Println(green(fmt.Sprintf("%s installed successfully.", packageName)))
}

func installNativePackage(packageName string, alias ...string) {
	var pkg string
	var command *exec.Cmd

	if len(alias) > 0 {
		pkg = alias[0]
	} else {
		pkg = packageName
	}

	if isPackageInstalled(pkg) {
		fmt.Println(green(packageName) + " is already installed... skipped")

		return
	}

	switch formValues.packageManager {
	case "apt":
		command = exec.Command("sudo", "apt", "install", "-y", pkg)
	case "dnf":
		command = exec.Command("sudo", "dnf", "install", "-y", pkg)
	case "pacman":
		command = exec.Command("sudo", "pacman", "-Suy", pkg)
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
