Step 0: Compile and put the executables in the bin folder as follows.

	./bin/car
	./bin/cloudProvider
	./bin/geth-timing

Step 1: Initialize a private ethereum network

    ./bin/geth-timing --datadir ./gethdata init ./private.json

Step 2: Run the node

    ./bin/geth-timing --datadir ./gethdata --timing.output=./output/deleteme.txt --networkid 1114 --rpc --rpcaddr 0.0.0.0 --rpcport 8545 --rpccorsdomain "*" --rpcapi "db,eth,net,web3,miner,personal" --nodiscover --mine --minerthreads 1 --unlock b3270be37a758e67a67fc6f2b62247cc58e0e61f --password ./password/password.txt

Step 3: Deploy the contract. Write down the contract address.

Step 4: Close the node

Step 5: Modify `./config/contract` with the valid contract address, e.g.:

    0x5E619911b8358861992365EDB67f51ffBC531618

Step 6: Execute `run.sh` to run the experiments
