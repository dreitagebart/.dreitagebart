package main

import (
	"embed"

	"github.com/fatih/color"
)

//go:embed git nvim/.config p10k tmux zsh
var templateFiles embed.FS

var (
	red   = color.New(color.FgRed).SprintfFunc()
	green = color.New(color.FgGreen).SprintfFunc()
	// cyan   = color.New(color.FgCyan).SprintFunc()
	// blue   = color.New(color.FgBlue).SprintFunc()
	// yellow = color.New(color.FgYellow).SprintFunc()
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
