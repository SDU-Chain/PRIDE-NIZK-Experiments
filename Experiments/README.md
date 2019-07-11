## Step 0: Prepare files
  Make sure the following commands are available.
  
  	bash python3 screen curl go gcc make truffle
  
  For example, on Ubuntu 18.04, you may install missing packages through apt-get.
  
  	apt install build-essential golang-go npm nodejs screen curl python3 bash
	npm -g install truffle

  Compile and put the executables in the bin folder as follows. 

	# Put the executables as follows:
	# ./bin/car
	# ./bin/cloudProvider
	# ./bin/geth-timing

  The following scripts might be handy.

	./compile_car_and_cloud.sh
	./compile_geth.sh

  If for some reasons the x-permission is lost, grant the x-permission on executables.

	chmod +x ./bin/*
	chmod +x ./*.sh
	chmod +x ./*.py

## Step 1: Prepare accounts
  The sealer account is already created. The address is placed at `./config/address`. The corresponding private key is placed at `./gethdata/keystore/` as well as `./gethaccount/sealer/keystore/`.
  
  However car accounts are not pre-created. Execute the following command to create 1000 car accounts. 

	echo 1000 | ./make_accounts.sh

## Step 2: Fund accounts
  The car accounts are not sealers, i.e. they can't mine (vote) to produce ether. The easiest way is to pre-fund them at the genesis block. It's advised to write a smart contract to fund new accounts in prodution, but we omit it because it is just an experiment.

	./make_genesis.py

## Step 3: [Server side] Initialize the blockchain and deploy the smart contract
  Set the account to the server account.

	echo 0 | ./set_geth_account.sh

  Note that the following script is written ONLY for truffle 5.0 and is not promised to be working in the future. If the script is not working, it's advised to write a program that send the RPC message `eth_sendTransaction` to deploy the contract as well as get the corresponding address, which should be write to `./config/contract`.

	./init_server_geth_and_deploy_contract.sh

## Step 4: [Server side] Start the sealer node and the cloud provider

  Find your ip address and write to `./config/bootnode_ip`.

	echo 192.168.1.1 > ./config/bootnode_ip

  Start the geth client as the sealer node, and output the timing log to `./output/foobar.txt`.

	echo foobar | ./run_server_geth_and_cloud.sh

  The script also start the cloud provider.

## Step 5: [Server side] Get the enode of the sealer node
  Note that if you just started ethereum client, please wait a few seconds for the client to prepare the RPC protocol. It might take ~10s.

	./get_geth_enode.sh

## Step 6: Start a car node

To be continued.
