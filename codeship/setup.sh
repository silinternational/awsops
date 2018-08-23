#!/usr/bin/env bash

# exit if any command fails
set -e

# Upgrade Go to 1.10.3
GO_VERSION=1.10.3
source /dev/stdin <<< "$(curl -sSL https://raw.githubusercontent.com/codeship/scripts/master/languages/go.sh)"

# Install dependencies
sudo apt-get update -y
sudo apt-get install -y awscli
go get github.com/golang/dep/cmd/dep
go get github.com/mitchellh/gox
cd ~/src/github.com/silinternational/awsops
dep ensure