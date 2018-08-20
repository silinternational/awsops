#!/usr/bin/env bash

wget -P /usr/local/bin/ https://raw.githubusercontent.com/silinternational/runny/0.2/runny
rc=$?;
if [[ $rc != 0]]; then
    exit $rc;
fi

chmod a+x /usr/local/bin/runny
rc=$?;
if [[ $rc != 0]]; then
    exit $rc;
fi

runny apt-get update -y
runny apt-get install -y awscli git
runny go get -u github.com/golang/dep/cmd/dep
runny go get github.com/mitchellh/gox

dep ensure