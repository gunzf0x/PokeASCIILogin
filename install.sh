#!/bin/bash

# Defining some colors to terminal
RED='\033[0;31m'
CYAN='\033[0;36m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
LGREEN='\033[1;32m'
LGRAY='\033[0;37m'
NC='\033[0m' # No Color

PROJECT_NAME='PokeASCIILogin'

REPO="github.com/p4nchit0z/${PROJECT_NAME}"
ROOT_FOLDER="./${PROJECT_NAME}"
DATA_ASCII="${ROOT_FOLDER}/ascii_pokemon_files"
DATA_JSON="${ROOT_FOLDER}/data"
ABS_PATH_TO_BINARY="$(pwd)/${PROJECT_NAME}/PokeASCIILogin"

# Just a pretty character to display
createStar(){
  echo -en "${RED}[${CYAN}$1${RED}]${NC}" 
}


# Prints what this script will do
printWhatWillBeDone () {
  echo -e "${RED}This script will do the following:\n${NC}"
  createStar "1"
  echo -e " Using 'git' command download data from Github repository"
  createStar "2"
  echo -e " Check if data has been downloaded correctly comparing md5sum hashes"
  createStar "3"
  echo -e " Check current shell (bash and zsh automatically supported; other shells will need next step to be 'manual')"
  createStar "4"
  echo -e " Adds a line at your .rc file (.bashrc or .zshrc supported) so every time you start your shell, PokeASCII will be executed"
  echo -n "    If you no longer want PokeASCII login just remove the line that initiates"
  echo " ${ABS_PATH_TO_BINARY} in your .rc file and delete the repository folder"
  createStar "${YELLOW}!"
  echo -en " Do you want to continue? ${LGREEN}[Y]es/[N]o${NC}: "
}


# Simply prints a horizontal line
printSeparation () {
  local START=1
  local END=$(tput cols)
  let "END -= 3"

  echo -en "\n${LGRAY}<"
  i=$START
  while [[ $i -le $END ]]
  do
    echo -en "="
    ((i = i + 1))
  done
  echo -e ">\n${NC}"
}


# Checks input given by the user if he/she wants to keep executing the script
checkUserChoice () {
  if [[ $1 =~ [Yy][Ee]?[Ss]? ]]; then
    return

  elif [[ $1 =~ [Nn][Oo]? ]]; then
    echo -e "\nPlease read $0 script with your favorite text editor carefully. Then re-run if you want to"
    echo -e "Or try installing 'PokeASCIILogin' "
    exit 1

  else
    echo -e "\n${RED}Invalid input. Please follow the instructions. Select '[Y]es' or '[N]o'${NC}\nBye."
    exit 1
  fi
}

# Check if 'git' command is available since it does not commes installed by default in many Linux systems
checkGit() {
  createStar "*"
  if ! command -v git &> /dev/null; then 
  echo -e " ${RED} Warning! 'git' is apparently not installed. Fix this and re-run the script."
  exit 1

  else 
    echo -e " ${LGREEN}'git' command detected${NC}"
    echo -e "\n    Cloning ${REPO}...\n"
    local repositoryURL="test"
    # git clone $repositoryURL
  fi
}


# Using m5dsum hashes, check if files have been downloaded correctly
checkDownloadedASCIIFiles () {
  # Check ASCII Pokemon files
  createStar "*"
  echo -e "${BLUE} Checking if files have been downloaded correctly...${NC}"
  # Check md5sum file and check its number of lines
  local md5nameASCII="${ROOT_FOLDER}/_md5sum_ascii_pokemon.txt"
  local md5LinesASCII=$(wc -l < $md5nameASCII)
  # Start checking every ASCII.txt file
  for downloadedFile in $DATA_ASCII/*.txt; do 
    local tempMD5=$(md5sum $downloadedFile | awk '{print $1}')
    local counter=1
    # For every ASCII file, check if their name and hashes do match
    while read pokehash filename
    do
      let "counter += 1"
      if [[ $filename == $(basename $downloadedFile) && $tempMD5 == $pokehash ]]; then 
       break
      elif [[ $counter -gt $md5LinesASCII ]]; then 
        echo -e "${RED}    Warning! No md5sum hash match found for '$downloadedFile' file${NC}"
      fi
    done < $md5nameASCII
  done
  
  # Second, check JSON files
  local md5nameJSON="${ROOT_FOLDER}/_md5sum_json_files.txt"
  local md5LinesJSON=$(wc -l < $md5nameJSON)
  for jsonFile in $DATA_JSON/*.json; do 
    local tempMD5=$(md5sum $jsonFile | awk '{print $1}')
    local counter=1
    while read pokehash filename
    do
      let "counter += 1"
      if [[ $filename == $(basename $jsonFile) && $tempMD5 == $pokehash ]]; then 
       break
      elif [[ $counter -gt $md5LinesJSON ]]; then 
        echo -e "${RED}    Warning! No md5sum hash match found for '$jsonFile' file${NC}"
      fi
    done < $md5nameJSON
  done

}


# Check which shell is the default one and add executable to .rc file if recognized (Bash or ZSH)
checkShellandAppendExec() {
  createStar "*"
  echo -e " ${BLUE}Checking config files and appending executable to them...${NC}"
  # Check what the default shell is (not the current one we are using) since I assume you want to install 
  # this login in your default terminal which is the one you usually use...
  currentShell=$(basename $SHELL)
  if [[ $currentShell =~ [Zz][Ss][Hh] ]]; then 
    echo "    - Detected default ZSH shell"
    # local rc_file="$HOME/.zshrc"
    local rc_file="$HOME/testing.txt" # Delete after test
    local loginBasename=$(basename $rc_file)

  elif [[ $currentShell =~ [Bb][Aa][Ss][Hh] ]]; then 
    echo "    - Detected Bash shell"
    #local rc_file="$HOME/.bashrc"
    local rc_file="$HOME/testing.txt" # Delete after test
    local loginBasename=$(basename $rc_file)

  else 
    echo -e "${RED}    Could not identify your current shell (not Bash neither ZSH)${NC}"
    local rc_file="YourConfigFileHere"
    echo -e "\n    Try running the following command:\n"
    echo -e "${YELLOW}    echo '.${ABS_PATH_TO_BINARY}' >> $rc_file\n${NC}"
    echo "    where, for example, ${rc_file} in ZSH is $HOME/.zshrc and in Bash is $HOME/.bashrc"
    echo -e "    Restart your terminal and enjoy!"
    exit 0
  fi

  if [ ! -f $rc_file ]; then
    echo -e "    ${RED}Warning! $rc_file could not be found at $HOME. Try adding the following line manually to your .rc file:\n${NC}"
    echo -e "    ${CYAN}.${ABS_PATH_TO_BINARY}${NC}"
    exit 0
  fi 

  local loginBasename=$(basename $rc_file)
  # Create a backup for the file if something goes wrong
  cp $rc_file "${ROOT_FOLDER}/${loginBasename}_pokeascii_backup"

  # Write in .rc files the executable
  echo -e "    - Appending executable to ${rc_file}...\n"
  echo "# Execute PokeASCIILogin binary" >> ${rc_file}
  echo ".${ABS_PATH_TO_BINARY}" >> ${rc_file}

  createStar "*"
  echo -e "${BLUE} Finally, remember to run\n    ${CYAN}chmod +x ${ABS_PATH_TO_BINARY}"
  echo -e "    source $rc_file${NC}"
}



# Remove some files that are no longer needed
endScript () {
  createStar "*"
  echo -e "${YELLOW} Installation complete! Restart your terminal and you are ready to pokego!${NC}"
}


main (){
  # Clear the screen
  clear

  # Print what script will execute
  printWhatWillBeDone

 
  # Ask user if wants to proceed
  read userChoice

  # Check user input with simple RegEx
  checkUserChoice "${userChoice}"

 # Print a simple line 
  printSeparation

  # Check if the user has 'git' installed on his/her machine
  checkGit

  # Check that ASCII Pokemon files have been downloaded correctly
  checkDownloadedASCIIFiles

  # Check which shell user is using and append binary execution to init file
  checkShellandAppendExec

  # Done! 
  endScript
}

main
