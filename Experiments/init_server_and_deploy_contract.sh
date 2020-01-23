#!/bin/bash
set -e

#echo 0 | ./set_geth_account.sh

./init_geth.sh
echo deleteme | ./run_server_geth_and_cloud.sh
./deploy_contract.sh
echo "sleep 5s..."
sleep 5
./kill_geth_and_cloud.sh
echo -n "Clear the output folder?[y/N] "
read ANS
if [ $ANS="Yes" -o $ANS="yes" -o $ANS="y" -o $ANS="Y" ] ; then
	./clear_output.sh
fi
