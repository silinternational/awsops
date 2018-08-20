#!/usr/bin/env bash

# Upgrade Go to 1.10.3
GO_VERSION=1.10.3
source /dev/stdin <<< "$(curl -sSL https://raw.githubusercontent.com/codeship/scripts/master/languages/go.sh)"

# Get runny for exit code checking on subsequent commands
sudo wget -P /usr/local/bin/ https://raw.githubusercontent.com/silinternational/runny/0.2/runny
sudo chmod a+x /usr/local/bin/runny

# Update packages and install dependencies
sudo runny apt-get update -y
sudo runny apt-get install -y awscli git
runny go get -u github.com/golang/dep/cmd/dep
runny go get github.com/mitchellh/gox

dep ensure