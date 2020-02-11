#!/bin/bash
set -e

runTests() {
    echo "#####################"
    echo test client"$1"
    pushd clients/"$1" && bash create_ch.sh && bash create_ca.sh && bash call_ca.sh
    popd
}

runTests 0
