#!/bin/bash
set -e
pushd ../Go/src/geth-timing/
make geth
popd
mv ../Go/src/geth-timing/build/bin/geth ./bin/geth-timing
