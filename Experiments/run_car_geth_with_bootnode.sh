#!/bin/bash

contract=`cat ./config/contract`
bootnode=`cat ./config/bootnode`
account=`cat ./config/account`

# stop when error occurs
set -e

read -p "Output filename=?" filename

# kill server process gently
./kill_geth_and_cloud.sh

# run server at background
screen -Sdm pride_exp_geth ./bin/geth-timing --timing.output=./output/$filename.txt --datadir ./gethdata --syncmode 'full' --port 30310 --bootnodes "$bootnode" --networkid 1114 --rpc --rpcaddr 0.0.0.0 --rpcport 8545 --rpccorsdomain "*" --rpcapi "admin,db,eth,net,web3,miner,personal" --unlock $account --password ./password/password.txt

