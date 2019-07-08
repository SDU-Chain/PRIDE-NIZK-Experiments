#!/bin/bash
set -e
./init_geth.sh
echo deleteme | ./run_geth_and_cloud.sh
./deploy_contract.sh
