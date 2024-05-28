#!/bin/bash

# color codes
COLOR_RED=31
COLOR_GREEN=32
COLOR_BLUE=34
COLOR_MAGENTA=35
COLOR_CYAN=36

function install_homebrew() {
  # Check installation of homebrew
  echo "Check installation of $(colorize $COLOR_CYAN "homebrew")..."

  if command -v brew >/dev/null; then
    echo "homebrew already installed - nothing to do"
  else
    echo "Installing homebrew..."

    # install homebrew with official curl command - see https://brew.sh
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  fi

  echo ""
}

function install_component() {
  local bin="$1"
  local component="$1"

  if [ "$#" == 2 ]; then
    bin="$2"
  fi

  local pid=""
  local message="Check installation of $(colorize $COLOR_CYAN $component)..."

  # Check installation of this component
  echo "$message"

  if command -v $bin > /dev/null; then
    message="$(colorize $COLOR_CYAN $component) "
    message+="is already installed - nothing to do"
    echo "$message"
  else
    message="Installing "
    message+="$(colorize $COLOR_CYAN $component)..."
    echo "$message"
    # install component with homebrew
    # and temporally store the PID
    brew install $component 2>errors.log & pid=$!
    # pid=$({ brew install $component &> /dev/null; echo $!; })
    # message="pid is $pid"
    # echo "$message"
    show_spinner $pid
  fi
  echo ""
}

function show_spinner() {
  local spinner="\|/-"
  local i=0
  
  while ps -p $1 > /dev/null
  do
    printf "\b%c" "${spinner:i++%4:1}"
    sleep 0.1
  done
}

function stow_zshrc() {
  # Check if ~/.zshrc already exists
  local valid_option=false
  local backup_path="$(get_backup_path)"
 
  while [ $valid_option = false ]; do
    if [ -f "$HOME/.zshrc" ]; then
      echo "$(colorize $COLOR_CYAN "~/.zshrc") already exists. Should I make a backup?"  
      read -p "Otherwise the file gets overwritten. (y/n) " choice
       case $choice in
        y) 
          echo "$(colorize $COLOR_CYAN "~/.zshrc") moved to folder $(colorize $COLOR_CYAN "~/.dreitagebart/backups/$TIMESTAMP")"
          mkdir -p $backup_path
          cp -L ~/.zshrc $backup_path
          valid_option=true
          break;;
        n) 
          echo "No backup file will be created - file gets overwritten"
          valid_option=true
          break;;
        *) 
          echo ""
      esac
    else
      break
    fi
  done

  stow zsh
}

function stow_git() {
  # Check if ~/.gitconfig already exists
  local valid_option=false
  local backup_path="$(get_backup_path)"

  while [ $valid_option = false ]; do
    if [ -f "$HOME/.gitconfig" ]; then
      echo "$(colorize $COLOR_CYAN "~/.gitconfig") already exists. Should I make a backup?"  
      read -p "Otherwise the file gets overwritten. (y/n) " choice
       case $choice in
        y) 
          echo "$(colorize $COLOR_CYAN "~/.gitconfig") moved to folder $(colorize $COLOR_CYAN "~/.dreitagebart/backups/$TIMESTAMP")"
          mkdir -p $backup_path
          cp -L ~/.gitconfig $backup_path
          valid_option=true
          break;;
        n) 
          echo "No backup file will be created - file gets overwritten"
          valid_option=true
          break;;
        *) 
          echo ""
      esac
    else
      break
    fi
  done

  stow git
}

function stow_nvim() {
  # Check if ~/.config/nvim already exists
  local valid_option=false
  local backup_path="$(get_backup_path ".config/nvim")"

  while [ $valid_option = false ]; do
    if [ -d "$HOME/.config/nvim" ]; then
      echo "$(colorize $COLOR_CYAN "~/.config/nvim") already exists. Should I make a backup?"  
      read -p "Otherwise the folder gets overwritten. (y/n) " choice
       case $choice in
        y) 
          echo "$(colorize $COLOR_CYAN "~/.config/nvim") moved to folder $(colorize $COLOR_CYAN "~/.dreitagebart/backups/$TIMESTAMP/.config/nvim")"
          mkdir -p $backup_path
          mv ~/.config/nvim $backup_path
          # cp -L ~/.config/nvim $backup_path
          valid_option=true
          break;;
        n) 
          echo "No backup for tmux plugins. Folder will be deleted..."
          rm -rf ~/.tmux/plugins
          valid_option=true
          break;;
        *) 
          echo ""
      esac
    else
      break
    fi
  done

  stow nvim
}

function stow_tmux() {
  # Check if ~/.tmux.conf already exists
  local valid_option=false
  local backup_path="$(get_backup_path)"

  while [ $valid_option = false ]; do
    if [ -f "$HOME/.tmux.conf" ]; then
      echo "$(colorize $COLOR_CYAN "~/.tmux.conf") already exists. Should I make a backup?"  
      read -p "Otherwise the file gets overwritten. (y/n) " choice
       case $choice in
        y) 
          echo "$(colorize $COLOR_CYAN "~/.tmux.conf") moved to folder $(colorize $COLOR_CYAN "~/.dreitagebart/backups/$TIMESTAMP")"
          mkdir -p $backup_path
          cp -L ~/.tmux.conf $backup_path
          valid_option=true
          break;;
        n) 
          echo "No backup file will be created - file gets overwritten"
          valid_option=true
          break;;
        *) 
          echo ""
      esac
    else
      break
    fi
  done

  stow tmux

  # Check if folder ~/.tmux/plugins already exist
  valid_option=false
  backup_path="$(get_backup_path ".tmux/plugins")"

  while [ $valid_option = false ]; do
    if [ -d "$HOME/.tmux/plugins" ]; then
      echo "$(colorize $COLOR_CYAN "~/.tmux/plugins") already exists. Should I make a backup?"  
      read -p "Otherwise the folder gets overwritten. (y/n) " choice
       case $choice in
        y) 
          echo "$(colorize $COLOR_CYAN "~/.tmux/plugins") moved to folder $(colorize $COLOR_CYAN "~/.dreitagebart/backups/$TIMESTAMP/.tmux/plugins")"
          mkdir -p $backup_path
          mv ~/.tmux/plugins $backup_path
          valid_option=true
          break;;
        n) 
          echo "No backup for tmux plugins. Folder will be deleted..."
          rm -rf ~/.tmux/plugins
          valid_option=true
          break;;
        *) 
          echo ""
      esac
    else
      break    
    fi
  done

  git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm
}

function install_shell_components() {
  install_component "stow"
  install_component "zsh"
  install_component "neovim" "nvim"
  install_component "tmux"
  install_component "fzf"
  install_component "bat"
  install_component "eza"
  install_component "zoxide"
  install_component "thefuck"
}

function install_dev_components() {
  install_component "nvm"
  install_component "pnpm"
  install_component "yarn"
}  

function colorize() {
  local color="$1"
  local text="$2"

  echo -e "\e[${color}m${text}\e[0m"
}

function print_dreitagebart() {
  echo "                                                                               "
	echo "                                 ▒▓▓████████▓▓▒                                "
  echo "                              ▒▓█████████████████▓                             "
  echo "                            ▒▓█████████████████████▓                           "
	echo "                           ▒████████▓▓▓▓▓▓▓██████████▒                         "
	echo "                          ▒██████▓████████████▓███████                         "
	echo "                          ███████▓▓▒▒▒▒▒▒▒▒▒▒▓▓███████▒                        "
	echo "                         ▒████▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▓████                        "
	echo "                         ▓█▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓██                        "
	echo "                         ▒▒▒▒▒▓▓▓▓▓▓▓▒▒▒▒▒▒▓▓▓▓▓▓▓▒▒▒▒▓                        "
	echo "                         ▒██▓▓▒▒▓▓▓▓▓▓█████▓▓▓▓▓▒▒▓▓▓█▒▒                       "
	echo "                        ▒▒▒█▒▒▒▒▒▓▓▒▒▒█▓▒█▓▒▒▓▓▒▒▒▒▒▒▓▒▒                       "
	echo "                         ▓▒▓▒▒▒▒▒▒▒▒▒▓▓▒▒▒▓▒▒▒▒▒▒▒▒▒▓▒▓▒                       "
	echo "                         ▓▒▒▓▓▓▓▓▓▓▓▓▒▒▒▒▒▒▓▓▓▓▓▓▓▓▓▓▒▓▒                       "
	echo "                         ▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒█▒                       "
	echo "                         ▒▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒█                        "
	echo "                          █▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▓▓                        "
	echo "                          ▓▓▓▒▒▒▒▒▓▓▓▓▓███▓▓▓▓▒▒▒▒▒▒▓▓▒                        "
	echo "                           █▓▓▒▒▒▓▓▓▒▓▒▒▒▒▓▓▓▓▓▒▒▒▓▓▓▓                         "
	echo "                           ▒▓▓▓▓▓█▓▒▒▒▓▓▓▓▒▒▒▒█▓▓▓▓▓▓                          "
	echo "                             ▓█▓███▒▒▒▒▓▓▓▒▒▒▓██▓▓▓▒                           "
	echo "                              ▒▓███▓▓▓▓▓▓▓▓▓▓▓███▓                             "
	echo "                                ▒▓████████████▓▒                               "
	echo "                                   ▒▓▓████▓▓▒                                  "
  echo ""
  echo "        _              _  _                        _                   _   "
  echo "     __| | _ __   ___ (_)| |_   __ _   __ _   ___ | |__    __ _  _ __ | |_ "
  echo "    / _\` || '__| / _ \| || __| / _\` | / _\` | / _ \| '_ \  / _\` || '__|| __|"
  echo " _ | (_| || |   |  __/| || |_ | (_| || (_| ||  __/| |_) || (_| || |   | |_ "
  echo "(_) \__,_||_|    \___||_| \__| \__,_| \__, | \___||_.__/  \__,_||_|    \__|"
  echo "                                      |___/                                "
  echo ""
  echo "Press ENTER to continue..."
  read -r
}

function show_final_message() {
  echo ""
  echo "- - - - - -"
  echo ""
  echo "Installation has been completed."
  echo "Wait... before you get your hands dirty. There are some other steps you should perform."
  echo ""
  echo "- POST INSTALLATION STEPS ---"
  echo ""
  echo "1.) Start $(colorize $COLOR_CYAN "tmux") and hit $(colorize $COLOR_BLUE "Control + Space") and $(colorize $COLOR_BLUE "I") (a capital i). This installs the tmux plugins"
  echo "2.) Start $(colorize $COLOR_CYAN "neovim") by typing $(colorize $COLOR_BLUE "nvim") into terminal and type $(colorize $COLOR_BLUE ":Lazy") + $(colorize $COLOR_BLUE "I"). This will install all neovim plugins."
  echo "    After that you should run command $(colorize $COLOR_BLUE ":MasonInstallAll") to install all LSPs"
  echo "Press ENTER to finish setup..."
  read -r
}

function get_timestamp() {
  # Get current time in seconds since the Epoch (1970-01-01 00:00:00 UTC)
  local current_time=$(date +%s)

  # Extract individual components using format specifiers
  local year=$(date +%Y -d @$current_time)  # Year (YYYY format)
  local month=$(date +%m -d @$current_time)  # Month (MM format)
  local day=$(date +%d -d @$current_time)    # Day (DD format)
  local hour=$(date +%H -d @$current_time)  # Hours (24-hour format, HH)
  local minutes=$(date +%M -d @$current_time) # Minutes (MM format)
  local seconds=$(date +%S -d @$current_time) # Seconds (SS format)

  TIMESTAMP="${year}${month}${day}-${hour}${minutes}${seconds}"
}

function get_backup_path() {
  local path="$HOME/.dreitagebart/backups/$TIMESTAMP/$1"

  echo "$path"
}

print_dreitagebart
get_timestamp
install_homebrew
install_shell_components
stow_zshrc
stow_tmux
stow_git
stow_nvim
show_final_message