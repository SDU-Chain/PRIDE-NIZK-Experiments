#!/bin/bash
read -p "count=?" count
contract="0x5E619911b8358861992365EDB67f51ffBC531618"
account="b3270be37a758e67a67fc6f2b62247cc58e0e61f"
killall geth-timing >/dev/null 2>&1
killall cloudProvider >/dev/null 2>&1
nohup ./bin/geth-timing --timing.output=./output/$count.txt --datadir ./gethdata --networkid 1114 --targetgaslimit 75200240 --rpc --rpcaddr 0.0.0.0 --rpcport 8545 --rpccorsdomain "*" --rpcapi "db,eth,net,web3,miner,personal" --nodiscover --mine --minerthreads 1 --unlock $account --password ./password/password.txt>/dev/null 2>&1 &
nohup ./bin/cloudProvider -port 12345 >/dev/null 2>&1 &
sleep 1
for i in {1..10}
do
./bin/car -count=$count -contract $contract -cloud localhost:12345;
done
sleep 1
killall cloudProvider
killall geth-timing


