# My personal dotfiles

This directory contains the dotfiles i use for my daily work.

# Wait... what the hell are dotfiles?

Dotfiles are configuration files used in Unix-based operating systems (including macOS and Linux) that typically start with a leading dot (.). These files store settings and preferences for various programs, the shell environment, and even your system as a whole.

Warning:
If you want to give the whole thing a chance, keep looking in the [“How can i use it?”](#how-can-i-use-it) section. Try to understand the code, change it to your needs and do not blindly accept everything unless you know what you are doing. Use at your own risk!

# How can I use it?

1. Clone this repository with following command \
   `git clone https://github.com/dreitagebart/.dreitagebart ~/.dreitagebart` \
   This repository assumes that it is stored in your user directory under ~/.dreitagebart
2. If you have not already done so, make the file `install.sh` executable with the command `chmod +x install.sh`
3. Run the installation process with `./install.sh`

# What does install.sh do?

This script will ask you a few questions and guide you through the installation process. The necessary components such as `fzf` or `eza` etc. are installed via [homebrew](https://brew.sh). Before the dotfiles are copied to your home directory, you will be ased for a backup. You will find the backup files under `~/.dreitagebart/backups/` in a subfolder with the name of the current timestamp. Most of the required dotfiles are actually symlinks pointing to `~/.dreitagebart`. For example the file `~/.zshrc` is a symlink to `~/.dreitagebart/zsh/.zshrc`. This has the advantage that all changes are kept synchronized.

# What's inside?

- `/zsh` \
  This folder contains the configurations and plugins for zsh shell.
- `/nvim` \
  Neovim settings powered by [NVChad](https://nvchad.com)
- `/tmux` \
  Tmux configurations, key bindings as well as tmux plugins.
- `/p10k` \
  My configuration for the zsh theme powerlevel10k
- `/vscodium` \
  I use [VSCodium](https://vscodium.com), not VSCode from microsoft - This thing is work in progress.

My personal workspace is powered by following awesome tools:

## Shell / Terminal

- [zsh](https://www.zsh.org/)
- [tmux](https://github.com/tmux/tmux)
- [neovim](https://neovim.io/)
- [nvchad](https://nvchad.com/)
- [powerlevel10k](https://github.com/romkatv/powerlevel10k)
- [zinit](https://github.com/zdharma-continuum/zinit)
- [fzf](https://github.com/junegunn/fzf)
- [eza](https://github.com/eza-community/eza)
- [bat](https://github.com/sharkdp/bat)
- [thefuck](https://github.com/nvbn/thefuck)
- [ripgrep](https://github.com/BurntSushi/ripgrep)
- [zoxide](https://github.com/ajeetdsouza/zoxide)

## Others

- [GNU stow](https://www.gnu.org/software/stow/)
- [homebrew](https://brew.sh/)

# Supported platforms

I have tested it with Fedora, Ubuntu and Arch (btw). Since Homebrew supports Linux and MacOS, maybe there will be support for Mac OS as well. I don't know, I don't have a Mac. But you can let me know if something doesn't work as expected by creating an issue.

# A few concluding sentences

Understand that this is probably the first bash script I have ever made. I'm not an expert on this, but you get the idea that I'm constantly learning.
I know that settings are very opinionated and branded. But hey, that's open source. You can fork the repository, customize it to your taste or simply leave it as it is.
I hope that it will help you or simply be useful in everyday life.

# Thanks to...

- Jake Wiesler, who got the ball rolling for me on the topic of .dotfiles - [github](https://github.com/jakewies/.dotfiles) / [youtube](https://www.youtube.com/watch?v=70YMTHAZyy4)
- Elliot Minns (dreams of code / autonomy on youtube) for giving me this inspiration in one some of his youtube videos - [github](https://github.com/elliottminns/dotfiles) / [youtube](https://www.youtube.com/watch?v=ud7YxC33Z3w)
- Mathias Bynense who pointed me into the right direction for custom dotfiles - [github](https://github.com/mathiasbynens/dotfiles)
