package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/huh/spinner"
)

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

// func installNixInstaller() {
// 	var command *exec.Cmd

// 	command = exec.Command("sh <(curl -L https://nixos.org/nix/install) --daemon")

// 	err := spinner.New().Type(spinner.MiniDot).ActionWithErr(func(context.Context) error {
// 		_, err := command.CombinedOutput()

// 		return err
// 	}).Run()

// 	if err != nil {
// 		fmt.Println(red(fmt.Sprintf("Failed to install nix installer: %v", err)))
// 		os.Exit(1)
// 	}
// }

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

func copyTemplateFiles(fsys embed.FS, destDir string) error {
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, path)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		content, err := fsys.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, content, 0644)
	})
}
