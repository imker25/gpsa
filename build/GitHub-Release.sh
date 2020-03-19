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
commitId=$(git rev-parse --verify HEAD)
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

echo "Create new tmp file $tmpJSON"
# contend for $tmpJSON acording to https://developer.github.com/v3/repos/releases/#create-a-release 
echo "{" >> $tmpJSON
echo "  \"tag_name\":\"$releaseName\"," >> $tmpJSON
echo "  \"target_commitish\":\"$commitId\"," >> $tmpJSON
echo "  \"name\":\"$releaseName\"," >> $tmpJSON
echo "  \"body\":\"$releaseDescribtion\"," >> $tmpJSON
echo "  \"draft\":true," >> $tmpJSON
echo "  \"prerelease\":$preTag" >> $tmpJSON
echo "}" >> $tmpJSON

curl -X POST --data @/dev/shm/GitHub-Release.json "https://api.github.com/repos/imker25/gpsa/releases?access_token=$apiToken"

popd

