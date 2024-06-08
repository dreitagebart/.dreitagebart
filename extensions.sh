# color codes
COLOR_RED=31
COLOR_GREEN=32
COLOR_BLUE=34
COLOR_MAGENTA=35
COLOR_CYAN=36

function colorize() {
  local color="$1"
  local text="$2"

  echo -e "\e[${color}m${text}\e[0m"
}

function main() {
  local valid_option=false
  local flag_upload=false
  local flag_download=false

  for arg in "$@"; do
    if [[ "$arg" == "--upload" ]]
    then
      flag_upload=true
    elif [[ "$arg" == "-u" ]]; then
      flag_upload=true
    fi

    if [[ "$arg" == "--download" ]]
    then
      flag_download=true
    elif [[ "$arg" == "-d" ]]; then
      flag_download=true
    fi
  done

  if [[ $flag_upload == true && $flag_download == true ]]
  then
    echo "You can only $(colorize $COLOR_CYAN "download") OR $(colorize $COLOR_CYAN "upload") your VSCodium extensions"
    exit 1
  fi

  if [[ $flag_upload == false && 
        $flag_download == false ]]
  then
    while [ $valid_option = false ]
    do
      echo "What do you want to do with your $(colorize $COLOR_CYAN "VSCodium") extensions?"  
      
      read -p "Upload (u) or download (d)?" choice
      
      case $choice in
        u) 
          echo "upload"
          flag_upload=true
          valid_option=true
          break;;
        d) 
          echo "download"
          flag_download=true
          valid_option=true
          break;;
        *) 
          echo ""
      esac
    done
  fi

  if [[ $flag_upload == true ]]
  then
    upload
  fi

  if [[ $flag_download == true ]]
  then
    download
  fi
}

function download() {
  mkdir -p ~/.dreitagebart/vscodium/.config/VSCodium/User
  cp ~/.config/VSCodium/User/settings.json ~/.dreitagebart/vscodium/.config/VSCodium/User/settings.json
  cp ~/.config/VSCodium/User/keybindings.json ~/.dreitagebart/vscodium/.config/VSCodium/User/keybindings.json
  cp -r ~/.config/VSCodium/User/snippets ~/.dreitagebart/vscodium/.config/VSCodium/User/snippets
  codium --list-extensions > ~/.dreitagebart/vscodium/extensions.dat

  echo "doing download"
}

function upload() {
  local file_path="$HOME/.dreitagebart/vscodium/extensions.dat"

  if [[ -f $file_path ]]
  then
    while IFS= read -r line
    do
      echo "Processing extension: $(colorize $COLOR_CYAN $line)"
      echo "- - - - - - - - - - - - - - - - - - - - - - - - - - -"
      codium --install-extension $line
      echo ""
    done < "$file_path"

    echo "All extensions installed"
  else 
    echo "Extension file at $(colorize $COLOR_CYAN "~/.dreitagebart/vscodium/extensions.dat") does not exist"
    exit 1
  fi
}

main $@

