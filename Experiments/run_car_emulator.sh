#!/bin/bash

contract=`cat ./config/contract`
account=`cat ./config/account`
geth_rpc_url=`cat ./config/geth_rpc_url || echo "http://localhost:8545"`

# stop when error occurs
set -e

# the number of commitments commited by a car 
read -p "Number of commitments per car=?" count

# run cars one by one
for i in {1..1}
do
./bin/car -count=$count -contract $contract -cloud `cat ./config/cloudprovider_ip_port` -ethereum $geth_rpc_url ;
done

