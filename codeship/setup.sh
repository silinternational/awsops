#!/usr/bin/env bash

# exit if any command fails
set -e

# Upgrade Go
GO_VERSION=1.16
source /dev/stdin <<< "$(curl -sSL https://raw.githubusercontent.com/codeship/scripts/master/languages/go.sh)"

# Install dependencies
sudo apt-get update -y
sudo apt-get install -y awscli
cd ~/src/github.com/silinternational/awsops
go get ./...
