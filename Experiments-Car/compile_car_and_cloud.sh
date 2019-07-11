#!/bin/bash
set -e
pushd ../Go
GOPATH=`pwd` go build -o ../Experiments/bin/car PRIDE-Exp/Exe/Car
GOPATH=`pwd` go build -o ../Experiments/bin/cloudProvider PRIDE-Exp/Exe/CloudProvider
popd
