#!/bin/bash

# Check installation of homebrew
echo "Check installation of homebrew..."
if command -v brew >/dev/null; then
  echo "homebrew already installed - nothing to do"
else
  echo "Installing homebrew..."

  # install homebrew with official curl command - see https://brew.sh
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
fi
echo ""

# Check installation of GNU stow
echo "Check installation of GNU stow..."
if command -v stow >/dev/null; then
  echo "GNU stow is already installed - nothing to do"
else
  echo "Installing GNU stow..."

  # install GNU stow with homebrew
  brew install stow
fi
echo ""

# Check install of zsh
echo "Check installation of zsh..."
if command -v zsh >/dev/null; then
  echo "zsh is already installed - nothing to do"
else
  echo "Installing zsh..."

  # install zsh with homebrew
  brew install zsh
fi
echo ""

# Check if ~/.zshrc already exists
valid_option=false
while [ $valid_option = false ]; do
  if [ -f "$HOME/.zshrc" ]; then
    echo "~/.zshrc already exists. Should I make a backup?"  
    read -p "Otherwise the file gets overwritten. (y/n) " choice

    case $choice in
      y) 
        echo "~/.zshrc moved to folder ~/.dreitagebart/backup/.zshrc"
        mkdir -p ~/.dreitagebart/backup
        mv ~/.zshrc ~/.dreitagebart/backup/.zshrc
        valid_option=true
        break;;
      n) 
        echo "No backup file will be created - file gets overwritten"
        valid_option=true
        break;;
      *) 
        echo ""
    esac
  fi
done
