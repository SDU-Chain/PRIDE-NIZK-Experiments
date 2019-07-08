#!/bin/bash

contract=`cat ./config/contract`
account=`cat ./config/account`

# stop when error occurs
set -e

# the number of commitments commited by a car 
read -p "Number of commitments per car=?" count

# run cars one by one
for i in {1..10}
do
./bin/car -count=$count -contract $contract -cloud localhost:12345;
done

