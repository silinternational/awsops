#!/usr/bin/env bash

# exit if any command fails
set -e

# set NPM cache
npm config set cache "${HOME}/cache/npm/"
export PATH="${HOME}/cache/npm/bin/:${PATH}"
export PREFIX="${HOME}/cache/npm/"
nvm install 9

# Upgrade Go to 1.10.3
GO_VERSION=1.10.3
source /dev/stdin <<< "$(curl -sSL https://raw.githubusercontent.com/codeship/scripts/master/languages/go.sh)"

# Get runny for exit code checking on subsequent commands
sudo wget -P /usr/local/bin/ https://raw.githubusercontent.com/silinternational/runny/0.2/runny
sudo chmod a+x /usr/local/bin/runny

# Install dependencies
sudo runny apt-get update -y
sudo runny apt-get install -y awscli
runny go get github.com/golang/dep/cmd/dep
runny go get github.com/mitchellh/gox
cd ~/src/github.com/silinternational/awsops
dep ensure