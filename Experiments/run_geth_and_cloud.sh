#!/bin/bash

contract=`cat ./config/contract`
account=`cat ./config/account`

# stop when error occurs
set -e

# the number of commitments commited by a car 
read -p "Number of commitments per car=?" count

# kill server process
screen -S pride_exp_geth -X quit >/dev/null 2>&1 || true
screen -S pride_exp_cloud -X quit >/dev/null 2>&1 || true
killall geth-timing >/dev/null 2>&1 || true
killall cloudProvider >/dev/null 2>&1 || true

# run server at background
screen -Sdm pride_exp_cloud ./bin/cloudProvider -port 12345
screen -Sdm pride_exp_geth ./bin/geth-timing --timing.output=./output/$count.txt --datadir ./gethdata --networkid 1114 --targetgaslimit 75200240 --rpc --rpcaddr 0.0.0.0 --rpcport 8545 --rpccorsdomain "*" --rpcapi "db,eth,net,web3,miner,personal" --nodiscover --mine --minerthreads 1 --unlock $account --password ./password/password.txt

