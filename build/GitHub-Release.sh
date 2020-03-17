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
    if [ "$3" = "true" || "$3" = "false"  ]; then
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
commitId=$(git describe --always --long)
tmpJSON="/dev/shm/GitHub-Release.json"


# ################################################################################################################
# functional code
# ################################################################################################################
pushd "$BRANCH_ROOT"

echo "Call arguments"
echo "name: \"$releaseName\"; describtion: \"$releaseDescribtion\"; prerelease \"$preTag\"; api-token: \"${apiToken:0:3}...\";"
echo ""
echo "Commit \"$commitId\" will be released"

if [ -f $tmpJSON ]; then
    echo "Delete old tmp file $tmpJSON"
    rm $tmpJSON
fi

# example contend for $tmpJSON acording to https://developer.github.com/v3/repos/releases/#create-a-release 
# {
#   "tag_name": "v1.0.0",
#   "target_commitish": "master",
#   "name": "v1.0.0",
#   "body": "Description of the release",
#   "draft": false,
#   "prerelease": false
# }

popd

