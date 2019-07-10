#!/bin/bash
set -e
pushd ../Go
GOPATH=`pwd` go build -o ../Experiments/bin/bootnode geth-timing/cmd/bootnode
popd
