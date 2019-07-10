## Step 0: Prepare files
  Make sure the following commands are available.
  
  	bash python3 screen curl go gcc make

  Compile and put the executables in the bin folder as follows.

	./bin/car
	./bin/cloudProvider
	./bin/geth-timing

  The following scripts might be handy.

	./compile_car_and_cloud.sh
	./compile_geth.sh

  If for some reasons the x-permission is lost, grant the x-permission on executables.

	chmod +x ./bin/*
	chmod +x ./*.sh
	chmod +x ./*.py

## Step 1: Prepare accounts
  The sealer account is already created with address b3270be37a758e67a67fc6f2b62247cc58e0e61f. The corresponding private key is placed at ./gethdata/keystore as well as ./gethaccount/sealer/keystore.
  Execute the following command to create 1000 car accounts. 

	echo 1000 | ./make_accounts.sh

## Step 2: Fund accounts
  The car accounts are not sealers, i.e. they can't mine (vote) to produce ether. The easiest way is to pre-fund them at the genesis block. It's advised to write a smart contract to fund new accounts in prodution, but we omit it because it is just an experiment.

	./make_genesis.py

## Step 3: Initialize the blockchain

	./init_geth.sh

## Step 4: Deploy the smart contract
  Note that the following script is written ONLY for truffle 5.0 and is not promised to be working in the future. If the script is not working, it's advised to write a program that send the RPC message `eth_sendTransaction` to deploy the contract as well as get the corresponding address, which should be write to `./config/contract`.

	./deploy_contract.sh

## Step 5: Start the sealer node and the cloud provider
  Start the geth client as the sealer node, and output the timing log to ./output/foobar.txt.

	echo foobar | ./run_geth_and_cloud.sh

  The script also start the cloud provider.

## Step 6: Get the enode of the sealer node
  Note, if you just started ethereum client, please wait a few seconds for the client to prepare the RPC protocol. It might take ~10s.

	./get_geth_enode.sh

## Step 7: Start a car node

To be continued.
