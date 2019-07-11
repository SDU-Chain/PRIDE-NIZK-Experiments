#!/bin/bash
set -e
read -p "What's the geth data folder?" custom_dir
./kill_geth_and_cloud.sh
pushd $custom_dir
rm geth.ipc || true
rm -rf geth/ || true
popd
./bin/geth-timing --datadir "$custom_dir" init ./genesis.json

