#!/bin/bash
set -e
wget `cat ./config/httpserver_url`/config/bootnode -O - | tee ./config/bootnode
wget `cat ./config/httpserver_url`/config/contract -O - | tee ./config/contract
rm -rf ./gethaccounts/cars || true
wget `cat ./config/httpserver_url`/gethaccounts/cars.tar.gz -O - | tar -xz -C ./gethaccounts
wget `cat ./config/httpserver_url`/genesis.json -O - > ./genesis.json
