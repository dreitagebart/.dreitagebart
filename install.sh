#!/bin/bash

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
  linux)
    case "$ARCH" in
      x86_64) BINARY="dotfiles-installer-linux-amd64";;
      *) echo "Unsupported architecture"; exit 1;;
    esac
    ;;
  # darwin)
  #   case "$ARCH" in
  #     x86_64) BINARY="dotfiles-installer-darwin-amd64";;
  #     arm64) BINARY="dotfiles-installer-darwin-arm64";;
  #     *) echo "Unsupported architecture"; exit 1;;
  #   esac
  #   ;;
  *) echo "Unsupported operating system"; exit 1;;
esac

curl -L "https://github.com/yourusername/yourrepo/releases/latest/download/$BINARY" -o dotfiles-installer
chmod +x dotfiles-installer
./dotfiles-installer
