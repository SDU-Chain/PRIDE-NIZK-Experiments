#!/bin/bash

contract=`cat ./config/contract`
account=`cat ./config/account`

# stop when error occurs
set -e

# the number of commitments commited by a car 
read -p "Number of commitments per car=?" count

./run_geth_and_cloud.sh

# wait for the geth client to unlock the account
# according to the CPU and disk I/O performance of your platform, the waiting time varies
# e.g: Raspberry Pi 3 on aarch64 needs ~8s
#      Intel Core i5-8500 on amd64 needs ~1s
echo "sleep 8s..."
sleep 8

# run cars one by one
echo $count | ./run_car.sh

sleep 1

# kill server process
./kill_geth_and_cloud.sh

sleep 1

