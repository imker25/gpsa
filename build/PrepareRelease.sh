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

# ################################################################################################################
# functional code
# ################################################################################################################

if [ "$1" == "-help" ]; then
    print_usage
    exit 0
fi  

pushd "$BRANCH_ROOT"

actualBranch=$(git status -b -s)

echo "Branch: \" $actualBranch \""

popd
exit 0