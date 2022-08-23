
function updateGo() {
  local DEBUG=${1:-false}

  if [[ $DEBUG == "true" ]]; then 
    set -x
  fi
  
  local VERSION="1.19"
  local TARGET_FILE_NAME="./go.tar.gz"

  if [ ! -f $TARGET_FILE_NAME ]; then
    wget -O $TARGET_FILE_NAME https://go.dev/dl/go${VERSION}.linux-amd64.tar.gz
  fi 

  sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf $TARGET_FILE_NAME

  ## Avoid accidental appending if the binary isn't alreayd set
  if [[ ! $PATH == *'/usr/local/go/bin'* ]]; then
    printf "\nexport PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
  fi

  if [[ $DEBUG == "true" ]]; then 
    set +x
  fi
}

updateGo $1

unset updateGo

