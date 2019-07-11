#!/bin/bash
set -e 
pushd ../SmartContract/truffle/
truffle migrate --reset | tee >( awk 'BEGIN{flg=0}{if($0=="2_deploy_contracts.js")flg=1;if(flg&&match($0,"^   > contract address\\:.*")!=0){print;exit;}}' | awk '{print $4}' | tee ../../Experiments/config/contract)
popd
