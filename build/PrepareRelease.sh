#!/bin/bash

# ################################################################################################################
# function definition
# ################################################################################################################
function print_usage()  {
    echo "Usage: $0 options"
    echo "-help    Print this help"
    echo "This script will do all needed steps to do a new gpsa release localy."
    echo "And ask in the end if the result should be published on GutHub"
}

# ################################################################################################################
# variable asigenment
# ################################################################################################################
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
BRANCH_ROOT="$SCRIPT_DIR/.."
LOG_DIR="$BRANCH_ROOT/logs"
RELEASE_BASE_BRANCH="feature/BetterVersioning"

# ################################################################################################################
# functional code
# ################################################################################################################

if [ "$1" == "-help" ]; then
    print_usage
    exit 0
fi  

pushd "$BRANCH_ROOT"

actualBranch=$(git status -b -s)
if [ $? != 0 ]; then
    echo "Error: Can not get branch information"
    popd
    exit -1
fi 

statusLines=$(echo "$actualBranch" | wc -l)
# if [ "$statusLines" != "1" ]; then
#     echo "Error: There are change files that are not checked in."
#     popd
#     exit -1
# fi 

# if [[ $actualBranch == *"[ahead"* ]]; then 
#     echo "Error: Local repository is ahead of remote"
#     popd
#     exit -1
# fi 

echo "Branch Status: \"$actualBranch\" "
expectedStatus="## $RELEASE_BASE_BRANCH...origin/$RELEASE_BASE_BRANCH [ahead 6]"
if [ "$expectedStatus" != "$actualBranch" ]; then
    echo "Error: Not running on $RELEASE_BASE_BRANCH branch"
    popd
    exit -1
fi 

popd
exit 0