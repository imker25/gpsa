#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd "$SCRIPT_DIR"

JENKINS_URL="https://homer.tobi.backfrak.de/jenkins"
# JENKINS_CRUMB is needed if your Jenkins master has CRSF protection enabled as it should
JENKINS_CRUMB=`curl "$JENKINS_URL/crumbIssuer/api/xml?xpath=concat(//crumbRequestField,\":\",//crumb)"`
# echo "$JENKINS_CRUMB"
curl -X POST -H "$JENKINS_CRUMB" -F "jenkinsfile=<../Jenkinsfile" $JENKINS_URL/pipeline-model-converter/validate

popd
