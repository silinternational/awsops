#!/usr/bin/env bash

# Exit if any command below fails
set -e
set -x

#if [[ "${CI_BRANCH}" != "develop" && "${CI_BRANCH}" != "master" ]]; then
#    # Only deploy for develop and master branches
#    exit 0
#fi

STAGE="dev"
if [[ "${CI_BRANCH}" == "master" ]]; then
    STAGE="prod"
fi

# Install serverless framework
npm install -g serverless

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

FUNCTIONS=( "ecs-right-size-cluster" )
for func in "${FUNCTIONS[@]}"
do
    echo "Deploying function ${func}..."
    cd $SCRIPTDIR/$func
    serverless deploy -v --stage $STAGE
done
