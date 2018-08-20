#!/usr/bin/env bash

# exit if any command fails
set -e

set -x

# array of target os/arch
targets=( "darwin/amd64" "linux/amd64" "linux/arm" "windows/386" )
distPath="dist"

# download gpg keys to use for signing
runny aws s3 cp s3://$KEY_BUCKET/secret.key ./
runny gpg --import secret.key

cd cli/
for target in "${targets[@]}"
do
    # Build binary using gox
    gox -osarch="${target}" -output="${distPath}/${target}/awsops"

    # If OS is windows, append .exe to filename before signing
    if [ "${target}" == "windows/386" ]
    then
        fileToSign="${distPath}/${target}/awsops.exe"
    else
        fileToSign="${distPath}/${target}/awsops"
    fi

    # Sign file with GPG
    runny gpg --yes -a -o "${fileToSign}.sig" --detach-sig $fileToSign

done

# Push dist/ to S3 under folder for CI_BRANCH (ex: develop or 1.2.3)
CI_BRANCH=${CI_BRANCH:="unknown"}
aws s3 sync --acl public-read dist/ s3://$DOWNLOAD_BUCKET/$CI_BRANCH/
cd ../