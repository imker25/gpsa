#!/bin/bash

# ################################################################################################################
# function definition
# ################################################################################################################
function print_usage()  {
    echo "Usage: $0 [filepath]"
    echo "This script will calculate the \"git hight\" for the given file"
    echo ""
    echo "The exit code of the script will be the \"git hight\" in case it 0 or grater. And will be -1 in error cases"
}

# ################################################################################################################
# prameter check
# ################################################################################################################
if [ "$1" = "" ]; then 
	echo " Error: No filepath given with the call"		
	print_usage
	exit -1
else 
	filepath=$1
fi

# ################################################################################################################
# variable asigenment
# ################################################################################################################
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
BRANCH_ROOT="$SCRIPT_DIR/.."

# ################################################################################################################
# functional code
# ################################################################################################################
pushd "$BRANCH_ROOT"

if [ -f "$filepath" ]; then 
    echo "Calculate \"git hight\" for \"$filepath\""
else 
    echo "File \"$filepath\" not found"
    exit -1
    popd
fi

lastFileChange=$(git log --pretty=format:"%H" -n 1  --follow "$filepath")
if [ $? -ne 0 ]; then
    echo "Error while searching for the last change in \"$filepath\""
    exit -1
fi

head=$(git log -n 1 --pretty=format:"%H")
if [ $? -ne 0 ]; then
    echo "Error while searching for the last change in repo"
    exit -1
fi

echo "Count commit between \"$lastFileChange\" and \"$head\""
commitCount=$(git rev-list --count $lastFileChange..$head)
if [ $? -ne 0 ]; then
    echo "Error while counting commits"
    exit -1
fi

echo "git hight: $commitCount"
exit $commitCount

popd
