#!/bin/bash

# ################################################################################################################
# function definition
# ################################################################################################################
function print_usage()  {
    echo "Usage: $0 [name] [describtion] [prerelease] [api-token]"
    echo "This script will create a GitHub release and upload the binary as asset to the created release"
    echo "    name            The name of the new release (e. g. \"V1.0.2\")."
    echo "    describtion     The describtion of the new release (e. g. \"Release of V1.0.2\")"
    echo "    prerelease      Tell if this release should have a prerelease tag. [true | false]"
    echo "    api-token       GitHub API token used for authentification"
}

# ################################################################################################################
# prameter check
# ################################################################################################################
if [ "$1" = "" ]; then 
	echo " Error: No name given with the call"		
	print_usage
	exit 1
else 
	releaseName=$1
fi

if [ "$2" = "" ]; then 
	echo " Error: No describtion given with the call"		
	print_usage
	exit 1
else 
	releaseDescribtion=$2
fi

if [ "$3" = "" ]; then 
	echo " Error: No prerelease given with the call"		
	print_usage
	exit 1
else 
    if [ "$3" = "true" ] || [  "$3" = "false"  ]; then
	    preTag=$3
    else
        echo " Error: prerelease can only be true or false"
        print_usage
        exit 1
    fi
fi

if [ "$4" = "" ]; then 
	echo " Error: No api-token given with the call"		
	print_usage
	exit 1
else 
	apiToken=$4
fi

# ################################################################################################################
# variable asigenment
# ################################################################################################################
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
BRANCH_ROOT="$SCRIPT_DIR/.."
commitId=$(git rev-parse --verify HEAD)
requestTmpJSON="/dev/shm/GitHub-Release-Request.json"
responseTmpJSON="/dev/shm/GitHub-Release-Response.json"
uploadTmpJSON="/dev/shm/GitHub-Upload-Asset-Response.json"
fileToUpload="$BRANCH_ROOT/bin/gpsa"


# ################################################################################################################
# functional code
# ################################################################################################################
pushd "$BRANCH_ROOT"

echo "Call arguments"
echo "name: \"$releaseName\"; describtion: \"$releaseDescribtion\"; prerelease \"$preTag\"; api-token: \"${apiToken:0:3}...\";"
echo ""
echo "Commit \"$commitId\" will be released"

if [ -f "$fileToUpload" ]; then
	echo "\"$fileToUpload\" will be uploaded"
else 
	echo "The file to upload \"$fileToUpload\" can not be found"
	exit 1
fi

if [ -f $requestTmpJSON ]; then
    echo "Delete old tmp file $requestTmpJSON"
    rm $requestTmpJSON
fi

echo "Create new tmp file $requestTmpJSON"
# contend for $requestTmpJSON acording to https://developer.github.com/v3/repos/releases/#create-a-release 
echo "{" >> $requestTmpJSON
echo "  \"tag_name\":\"$releaseName\"," >> $requestTmpJSON
echo "  \"target_commitish\":\"$commitId\"," >> $requestTmpJSON
echo "  \"name\":\"$releaseName\"," >> $requestTmpJSON
echo "  \"body\":\"$releaseDescribtion\"," >> $requestTmpJSON
echo "  \"draft\":false," >> $requestTmpJSON
echo "  \"prerelease\":$preTag" >> $requestTmpJSON
echo "}" >> $requestTmpJSON

if [ -f $responseTmpJSON ]; then
    echo "Delete old tmp file $responseTmpJSON"
    rm $responseTmpJSON
fi

curl --data @"$requestTmpJSON" -H "Authorization: $apiToken" -X POST "https://api.github.com/repos/imker25/gpsa/releases?" > $responseTmpJSON
if [ $? -eq 0 ]; then
	echo "No error in curl"
else
	echo "curl reported a error code"
	exit 1
fi
releaseID=$(cat $responseTmpJSON | jq -r ".id")
if [[ "$releaseID" == "" || "$releaseID" == "null" ]]; then 
	echo "No release ID found in the response $responseTmpJSON"
	exit 1
fi

uploadURL=$(cat $responseTmpJSON | jq -r ".upload_url")
if [[ "$uploadURL" == ""  ||  "$uploadURL" == "null" ]]; then 
	echo "No upload_url found in the response $responseTmpJSON"
	exit 1
fi

if [ -f $uploadTmpJSON ]; then
    echo "Delete old tmp file $uploadTmpJSON"
    rm $uploadTmpJSON
fi

realUploadUrl="${uploadURL::-13}"
echo "Release with ID $releaseID was created"
echo "Upload \"$fileToUpload\" to $realUploadUrl"
curl --data @"$fileToUpload" -H "Content-Type: application/zip" -H "Authorization: $apiToken" -X POST "$realUploadUrl?name=gpsa&label=linux-executabel" > $uploadTmpJSON
if [ $? -eq 0 ]; then
	echo "No error in curl"
else
	echo "curl reported a error code"
	exit 1
fi

assetID=$(cat $uploadTmpJSON | jq -r ".id")
downloadURL=$(cat $uploadTmpJSON | jq -r ".browser_download_url")
if [[ "$assetID" == ""  ||  "$releaseID" == "null" ]]; then 
	echo "No asset ID found in the response $uploadTmpJSON"
	exit 1
fi
if [[ "$downloadURL" == "" ||  "$releaseID" == "null" ]]; then 
	echo "No download URL found in the response $uploadTmpJSON"
	exit 1
fi

echo "The asset with ID $assetID can be downloaded from $downloadURL"

popd

exit 0
