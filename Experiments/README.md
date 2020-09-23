## Step 0: [Both sides] Prepare files
  Make sure the following commands are available.

```bash
bash python3 screen curl wget go gcc make truffle tar
```

  Note that this experiment is done in 2019. Go 1.14 is not supported. You may choose the go toolchain that comes with Ubuntu 18.04.

  For example, on Ubuntu 18.04, you may install missing packages through apt-get.

```bash
apt install build-essential golang-go npm nodejs screen curl wget python3 bash tar
npm -g install truffle
```

  Compile and put the executables in the bin folder as follows. 

```bash
# Put the executables as follows:
# ./bin/car
# ./bin/cloudProvider
# ./bin/geth-timing
```

  The following scripts might be handy.

```bash
./compile_car_and_cloud.sh
./compile_geth.sh
```

  If for some reasons the x-permission is lost, grant the x-permission on executables.

```bash
chmod +x ./bin/*
chmod +x ./*.sh
chmod +x ./*.py
```

## Step 1: [Server side] Prepare accounts
  The sealer account is already created. The address is placed at `./config/address`. The corresponding private key is placed at `./gethaccount/sealer/keystore/`.

  However car accounts are not pre-created. Execute the following command to create 1000 car accounts. 

```bash
echo 1000 | ./make_accounts.sh
```

## Step 2: [Server side] Fund accounts
  The car accounts are not sealers, i.e. they can't mine (vote) to produce ether. The easiest way is to pre-fund them at the genesis block. It's advised to write a smart contract to fund new accounts for prodution use, but we omit it because it is just an experiment.

```bash
./make_genesis.py
```

## Step 3: [Server side] Initialize the blockchain
  Set the account to the server account.

```bash
echo 0 | ./set_geth_account.sh
```

## Step 4: [Server side] Start the sealer node and the cloud provider

  Find your ip address and write to `./config/bootnode_ip`.

```bash
echo 192.168.1.1 > ./config/bootnode_ip
```

  Start the geth client as the sealer node, and output the timing log to `./output/foobar.txt`.

```bash
echo foobar | ./run_server_geth_and_cloud.sh
```

  The script also start the cloud provider. It wait a few seconds and then get the enode address.

  Start the HTTP server powered by Python 3. The HTTP server provide information about the sealer node for the cars.

```bash
./run_http_server.sh
```

## Step 5: [Client side] Get information from the HTTP server

  Edit these two files:
  - `./config/cloudprovider_ip_port`
  - `./config/httpserver_url`

  Download configuration from the server.

  ```bash
./download_config_from_http_server.sh
  ```


## Step 6: [Client side] Start the node

  Select which car account to use. The first one, for example.
```bash
echo 1 | ./set_geth_account.sh
```

  Start the node.

```bash
./run_car_geth_with_bootnode.sh
```

  Wait ~30 seconds. Check if the node is connected to the sealer node.

```bash
./get_geth_nodelist.sh
```

## Step 7: [Server side] Deploy the smart contract

  Note that the following script is written ONLY for truffle 5.0 and is not promised to be working in the future. If the script is not working, it's advised to write a program that send the RPC message `eth_sendTransaction` to deploy the contract as well as get the corresponding address, which should be write to `./config/contract`.

```bash
./deploy_contract.sh
```


## Step 8: [Client side] Start a car

The car generates 1000 commitments and upload the proof.


```bash
echo 1000 | ./run_car_emulator.sh
```

## Step 9: [Client side] kernel.pid_max

Linux kernel has a limitation on maximum number of process in order to prevent fork bombs.

For experiment 2, the client might start a huge number of cars, in which case it will reach the limitation.

Make sure the Linux kernel is running on 64-bit mode. The value must be a power of 2, e.g. 32768. The maximum value is 4194304 for a 64-bit kernel.

To avoid some strange behaviors, the following commands and experiment 2 should be run by root, instead of `sudo`. Use `sudo -s` command to switch to root account.

```bash
echo 4194304 > /proc/sys/kernel/pid_max
ulimit -u 4194304
ulimit -n 1027204
```

To make these option persists after reboot, edit `/etc/sysctl.conf` by adding the following line:

```ini
kernel.pid_max = 4194304
```

and `/etc/security/limits.conf` file by adding the following line:

```
* soft nproc 4194304
* hard nproc 4194304
* soft nofile 1027204
* hard nofile 1027204
```

Be sure to check all files at `/etc/sysctl.d/` and `/etc/security/limits.d/` in case your configuration be overridden.

## Step 10: [Both sides] Stop the nodes

The following command stop the ethereum client **safely**. It sends the `Ctrl+C` signal to the process, instead of just killing them.

```bash
./kill_geth_and_cloud.sh
```

