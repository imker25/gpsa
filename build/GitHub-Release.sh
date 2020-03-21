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

function upload_asset()  { 
	fileToUpload="$1"
	label="$2"
	fileName=$(basename $fileToUpload)
	upoadResponseJSON="$LOG_DIR/GitHub-Upload-Asset-Response.$fileName.json"

	if [ -f "$fileToUpload" ]; then
		echo "Upload \"$fileToUpload\" to $realUploadUrl"
	else 
		echo "The file to upload \"$fileToUpload\" can not be found"
		exit 1
	fi
	
	curl --data @"$fileToUpload" -H "Content-Type: application/zip" -H "Authorization: token $apiToken" -X POST "$realUploadUrl?name=$fileName&label=$label" > "$upoadResponseJSON"
	if [ $? -eq 0 ]; then
		echo "No error in curl"
	else
		echo "curl reported a error code"
		exit 1
	fi

	assetID=$(cat "$upoadResponseJSON" | jq -r ".id")
	downloadURL=$(cat "$upoadResponseJSON" | jq -r ".browser_download_url")
	if [[ "$assetID" == ""  ||  "$releaseID" == "null" ]]; then 
		echo "No asset ID found in the response \"$upoadResponseJSON\""
		exit 1
	fi
	if [[ "$downloadURL" == "" ||  "$releaseID" == "null" ]]; then 
		echo "No download URL found in the response \"$upoadResponseJSON\""
		exit 1
	fi

	echo "The asset with ID $assetID can be downloaded from $downloadURL"

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
LOG_DIR="$BRANCH_ROOT/logs"
commitId=$(git rev-parse --verify HEAD)
releaseRequestJSON="$LOG_DIR/GitHub-Release-Request.json"
releaseResponseJSON="$LOG_DIR/GitHub-Release-Response.json"

# fileToUpload="$BRANCH_ROOT/bin/gpsa"


# ################################################################################################################
# functional code
# ################################################################################################################
pushd "$BRANCH_ROOT"

echo "Call arguments"
echo "name: \"$releaseName\"; describtion: \"$releaseDescribtion\"; prerelease \"$preTag\"; api-token: \"${apiToken:0:3}...\";"
echo ""
echo "Commit \"$commitId\" will be released"

if [ -f "$releaseRequestJSON" ]; then
    echo "Delete old request file \"$releaseRequestJSON\""
    rm "$releaseRequestJSON"
fi

echo "Create new request file \"$releaseRequestJSON\""
# contend for $releaseRequestJSON acording to https://developer.github.com/v3/repos/releases/#create-a-release 
echo "{" >> "$releaseRequestJSON"
echo "  \"tag_name\":\"$releaseName\"," >> "$releaseRequestJSON"
echo "  \"target_commitish\":\"$commitId\"," >> "$releaseRequestJSON"
echo "  \"name\":\"$releaseName\"," >> "$releaseRequestJSON"
echo "  \"body\":\"$releaseDescribtion\"," >> "$releaseRequestJSON"
echo "  \"draft\":false," >> "$releaseRequestJSON"
echo "  \"prerelease\":$preTag" >> "$releaseRequestJSON"
echo "}" >>"$releaseRequestJSON"

if [ -f "$releaseResponseJSON" ]; then
    echo "Delete old tmp file \"$releaseResponseJSON\""
    rm "$releaseResponseJSON"
fi

curl --data @"$releaseRequestJSON" -H "Authorization: token $apiToken" -X POST "https://api.github.com/repos/imker25/gpsa/releases?" > "$releaseResponseJSON"
if [ $? -eq 0 ]; then
	echo "No error in curl"
else
	echo "curl reported a error code"
	exit 1
fi
releaseID=$(cat "$releaseResponseJSON" | jq -r ".id")
if [[ "$releaseID" == "" || "$releaseID" == "null" ]]; then 
	echo "No release ID found in the response \"$releaseResponseJSON\""
	exit 1
fi

uploadURL=$(cat "$releaseResponseJSON" | jq -r ".upload_url")
if [[ "$uploadURL" == ""  ||  "$uploadURL" == "null" ]]; then 
	echo "No upload_url found in the response \"$releaseResponseJSON\""
	exit 1
fi

if [ -f "$upoadResponseJSON" ]; then
    echo "Delete old tmp file \"$upoadResponseJSON\""
    rm "$upoadResponseJSON"
fi

realUploadUrl="${uploadURL::-13}"
echo "Release with ID $releaseID was created"

upload_asset "$BRANCH_ROOT/bin/gpsa" "linux-executable"
# upload_asset "$BRANCH_ROOT/bin/gpsa.exe" "windows-executable"



popd

exit 0
