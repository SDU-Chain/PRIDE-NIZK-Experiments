#!/bin/bash

contract=`cat ./config/contract`
bootnode=`cat ./config/bootnode`

# stop when error occurs
set -e

read -p "Car account id=?" carid

car_account=`cat ./gethaccounts/cars/$carid/keystore/* | head -n 1 | ./parse_json.py \"address\"`
echo "Car account: $car_account"

# kill server process gently
./kill_geth_and_cloud.sh

# run server at background
screen -Sdm pride_exp_geth ./bin/geth-timing --timing.output=./output/$carid.txt --datadir ./gethaccounts/cars/$carid --syncmode 'full' --port 30310 --bootnodes "$bootnode"  --networkid 1114 --targetgaslimit 75200240 --rpc --rpcaddr 0.0.0.0 --rpcport 8545 --rpccorsdomain "*" --rpcapi "admin,db,eth,net,web3,miner,personal" --mine --minerthreads 1 --unlock $car_account --password ./password/password.txt


