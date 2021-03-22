#!/bin/bash

# Script to test the gpsa comandline interface
# The script expects to find the gpsa executable ./../bin/gpsa to exist
# You might want to run the gradle build command before

# Uses https://github.com/lehmannro/assert.sh as unit test framework.

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
gpsa="$SCRIPT_DIR/../bin/gpsa"
csv_out="$SCRIPT_DIR/../logs/test.csv"
json_out="$SCRIPT_DIR/../logs/test.json"
testdata="$SCRIPT_DIR/../testdata"
valid_gpx="$testdata/valid-gpx/*.gpx"


pushd "$SCRIPT_DIR"

echo "# ###################################################################"
echo "$(date) - Prepare for testing"
echo "# ###################################################################"
if [ -f "$SCRIPT_DIR/assert.sh" ]; then
    echo "Remove old $SCRIPT_DIR/assert.sh"
    rm -rf "$SCRIPT_DIR/assert.sh"
fi
wget -O "$SCRIPT_DIR/assert.sh" https://raw.githubusercontent.com/lehmannro/assert.sh/v1.1/assert.sh
 
if [ -f "$SCRIPT_DIR/assert.sh" ]; then
    chmod 700 "$SCRIPT_DIR/assert.sh"
    source "$SCRIPT_DIR/assert.sh"
else
    popd
    echo "Error while getting https://github.com/lehmannro/assert.sh"
    exit -1
fi


if [ -f "$gpsa" ]; then
    echo "Use testee $gpsa"
else
    popd
    echo "Testee $gpsa does not exist"
    exit -1
fi


echo "# ###################################################################"
echo "$(date) - Start testing"
echo "# ###################################################################"


# Test the version output
assert_raises "$gpsa -version | grep Version: &> /dev/null" 0

# Test output to file
if [ -f "$csv_out" ]; then
    rm -rf "$csv_out"
fi
assert "$gpsa -out-file=$csv_out $valid_gpx" ""
if [ ! -f "$csv_out" ]; then
    assert "echo $csv_out not found" ""
fi

if [ -f "$csv_out" ]; then
    rm -rf "$csv_out"
fi
assert_raises "$gpsa -verbose -out-file=$csv_out $valid_gpx | grep \"16 of 16 files processed successfully\" &> /dev/null" 0
if [ ! -f "$csv_out" ]; then
    assert "echo $csv_out not found" ""
fi

if [ -f "$json_out" ]; then
    rm -rf "$json_out"
fi
assert "$gpsa -out-file=$json_out $valid_gpx" ""
if [ ! -f "$json_out" ]; then
    assert "echo $json_out not found" ""
fi

# Test stdout
assert_raises "$gpsa $valid_gpx | grep \"not valid\" &> /dev/null" 0
assert_raises "$gpsa -verbose $valid_gpx | grep \"16 of 16 files processed successfully\" &> /dev/null" 0
assert_raises "$gpsa -verbose -skip-error-exit $valid_gpx $testdata/invalid-tcx/02.tcx | grep \"16 of 17 files processed successfully\" &> /dev/null" 0
assert_raises "$gpsa -verbose -skip-error-exit $valid_gpx $testdata/invalid-tcx/02.tcx | grep \"At least one error occurred\" &> /dev/null" 1

# Test stdin
if [ -f "$csv_out" ]; then
    rm -rf "$csv_out"
fi
assert_raises "find $testdata -name \"*.gpx\" | $gpsa -verbose -skip-error-exit -out-file=$csv_out " 255
if [ ! -f "$csv_out" ]; then
    assert "echo $csv_out not found" ""
fi
assert_raises "find $testdata/valid-gpx -name \"*.gpx\" | $gpsa " 0

popd
assert_end gpsa_IntegrationTests


