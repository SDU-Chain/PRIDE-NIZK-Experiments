#!/bin/bash
set -e
./kill_geth_and_cloud.sh
pushd ./gethdata
rm geth.ipc || true
rm -rf geth/ || true
popd
./bin/geth-timing --datadir ./gethdata init ./genesis.json


