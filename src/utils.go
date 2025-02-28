package main

import (
	"bytes"
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
	"gopkg.in/ini.v1"
)

func isPackageInstalled(packageName string) bool {
	command := exec.Command("which", packageName)

	_, err := command.CombinedOutput()

	return err == nil
}

func isCaskInstalled(caskName string) (bool, error) {
	cmd := exec.Command("brew", "list", "--casks")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		return false, fmt.Errorf("failed to run brew list --casks: %w", err)
	}

	output := out.String()
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == caskName {
			return true, nil
		}
	}

	return false, nil
}

func setDefaultShell() {
	cacheSudoPassword()

	user := os.Getenv("USER")

	if user == "" {
		fmt.Println(red(fmt.Printf("USER environment variable not set")))
		return
	}

	zshPath, err := exec.LookPath("zsh")

	if err != nil {
		fmt.Println(red(fmt.Sprintf("Failed to lookup zsh file path: %v", err)))
		return
	}

	command := exec.Command("sudo", "chsh", "-s", zshPath, user)
	command.Stdin = os.Stdin
	// command.Stdout = os.Stdout
	// command.Stderr = os.Stderr
	err = command.Run()

	if err != nil {
		fmt.Println(red(fmt.Sprintf("Failed to set default shell: %v", err)))
		return
	}

	fmt.Println("ZSH is now your default shell")
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

func configureGit() {
	gitConfigPath := path.Join(formValues.homeDir, ".dreitagebart", "git", "dot-gitconfig")

	cfg, err := ini.Load(gitConfigPath)

	if err != nil {
		fmt.Println("Could not configure .gitconfig file")
		os.Exit(1)
	}

	gitUserSection, err := cfg.GetSection("user")
	if err != nil {
		gitUserSection, err = cfg.NewSection("user")

		if err != nil {
			fmt.Println("Error creating user section: ", err)
			return
		}
	}

	gitUserSection.Key("name").SetValue(formValues.gitConfig.Name)
	gitUserSection.Key("email").SetValue(formValues.gitConfig.Email)

	err = cfg.SaveTo(gitConfigPath)
	if err != nil {
		fmt.Println("Error saving .gitconfig:", err)
		return
	}
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
	// var command *exec.Cmd

	if isPackageInstalled("brew") {
		fmt.Println(green("homebrew") + " is already installed... skipped")

		return
	}

	cacheSudoPassword()

	command := exec.Command("/bin/bash", "-c", "curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh | bash")

	err := spinner.New().
		Type(spinner.MiniDot).
		Title(" Installing homebrew...").
		ActionWithErr(func(context.Context) error {
			command.Stderr = os.Stderr
			err := command.Run()

			return err
		}).
		Accessible(true).
		Run()

	if err != nil {
		fmt.Println(red(fmt.Sprintf("Failed to install homebrew: %v", err)))
		os.Exit(1)
	}

	command = exec.Command("/bin/bash", "-c", "eval \"$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)\"")
	err = command.Run()
	if err != nil {
		fmt.Println(red(fmt.Sprintf("Failed to eval homebrew: %v", err)))
		os.Exit(1)
	}

	// if err != nil {
	// 	fmt.Println(red(fmt.Sprintf("Failed to install homebrew: %v", err)))
	// 	os.Exit(1)
	// }
}

func installHomebrewFont(fontName string) {
	isInstalled, err := isCaskInstalled(fontName)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if isInstalled {
		return
	}

	command := exec.Command("brew", "install", "--cask", fontName)

	err = spinner.New().
		Type(spinner.MiniDot).
		Title(" Installing font " + fontName + "...").
		ActionWithErr(func(context.Context) error {
			command.Stderr = os.Stderr
			err := command.Run()

			return err
		}).Run()

	if err != nil {
		fmt.Println(red(fmt.Sprintf("Failed to install %s: %v", fontName, err)))

		os.Exit(1)
	}

	fmt.Println(green(fmt.Sprintf("Font %s installed successfully.", fontName)))
}

func installHomebrewPackage(packageName string, alias ...string) {
	var pkg string

	if len(alias) > 0 {
		pkg = alias[0]
	} else {
		pkg = packageName
	}

	if isPackageInstalled(pkg) {
		fmt.Println(green(packageName) + " is already installed... skipped")

		return
	}

	command := exec.Command("brew", "install", pkg)

	err := spinner.New().
		Type(spinner.MiniDot).
		Title(" Installing package " + packageName + "...").
		ActionWithErr(func(context.Context) error {
			// command.Stdin = os.Stdin
			// command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			err := command.Run()

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

	cacheSudoPassword()

	// command = exec.Command("/bin/bash", "-c", "sudo echo \"\"")
	// command.Stdin = os.Stdin
	// command.Stdout = os.Stdout
	// command.Stderr = os.Stderr
	// err := command.Run()

	// if err != nil {
	// 	fmt.Println(red(fmt.Sprintf("%v", err)))
	// 	os.Exit(1)
	// }

	switch formValues.packageManager {
	case "apt":
		command = exec.Command("sudo", "apt", "install", "-y", packageName)
	case "dnf":
		// command = exec.Command("sudo", "dnf", "install", "-y", packageName)

		command = exec.Command("/bin/bash", "-c", "sudo dnf install -y "+packageName)
	case "pacman":
		command = exec.Command("sudo", "pacman", "-Suy", packageName)
	}

	// command.Stdin = os.Stdin
	// command.Stdout = os.Stdout
	// command.Stderr = os.Stderr
	// err := command.Run()

	// return err

	err := spinner.New().
		Type(spinner.MiniDot).
		Title(" Installing package " + packageName + "...").
		ActionWithErr(func(context.Context) error {
			// command.Stdin = os.Stdin
			// command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			err := command.Run()

			return err
		}).
		// Accessible(true).
		Run()

	if err != nil {
		fmt.Println(red(fmt.Sprintf("Failed to install %s: %v", packageName, err)))

		os.Exit(1)
	}

	fmt.Println(green(fmt.Sprintf("%s installed successfully.", packageName)))
}

func cacheSudoPassword() {
	command := exec.Command("sudo", "-v")
	err := command.Run()

	if err != nil {
		fmt.Println(red(fmt.Sprintf("%v", err)))
		os.Exit(1)
	}
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

func moveFolderIfExists(sourcePath, destinationPath string) {
	fileInfo, err := os.Lstat(sourcePath)

	if err != nil {
		if os.IsNotExist(err) {
			fmt.Sprintln("source path does not exist: %w", err)
			return
		}

		fmt.Sprintln("failed to stat source path: %w", err)
		return
	}

	if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		err = os.Rename(sourcePath, destinationPath)

		if err != nil {
			fmt.Sprintln("failed to move symlink: %w", err)
			return
		}

		return
	}

	if !fileInfo.IsDir() {
		fmt.Sprintln("source path is not a directory")
		return
	}

	// Ensure destination directory exists.
	destDir := filepath.Dir(destinationPath)
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		err := os.MkdirAll(destDir, 0755)
		if err != nil {
			fmt.Sprintln("failed to create destination directory: %w", err)
			return
		}
	}

	err = os.Rename(sourcePath, destinationPath)
	if err != nil {
		fmt.Sprintln("failed to move folder: %w", err)
		return
	}

	return
}

func stowFile(filename string, template string) {
	var command *exec.Cmd

	stowPath := path.Join(formValues.homeDir, filename)
	templatePath := path.Join(formValues.homeDir, ".dreitagebart")

	if fileExists(stowPath) {
		fileInfo, err := os.Stat(stowPath)

		if err != nil {
			fmt.Printf("Could not read fileinfo: %s", err)
			os.Exit(1)
		}

		if fileInfo.IsDir() {
			moveFolderIfExists(stowPath, formValues.backupPath)
		} else {
			err = os.Rename(stowPath, path.Join(formValues.backupPath, filename))

			if err != nil {
				fmt.Printf("failed to move file: %s", err)
				os.Exit(1)
			}
		}
	}

	command = exec.Command("stow", template, "--dir="+templatePath, "--dotfiles")

	err := spinner.New().Type(spinner.MiniDot).ActionWithErr(func(context.Context) error {
		_, err := command.CombinedOutput()

		return err
	}).Run()

	if err != nil {
		fmt.Println(red(fmt.Sprintf("Failed to stow %s", filename)))
		os.Exit(1)
	}
}
