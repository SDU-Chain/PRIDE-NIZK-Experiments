#!/bin/bash
set -e

read -p "Car account id (0 represents the server)=?" carid
if ! [[ $carid == ?(-)+([0-9]) ]] ; then
	echo "error: Not a number" >&2; exit 1
fi

account_dir="./gethaccounts/cars/$carid"

if [[ "s$carid" == "s0" ]] ; then
	account_dir="./gethaccounts/sealer"
fi

account=`cat "$account_dir/keystore/"* | head -n 1 | ./parse_json.py \"address\"`
echo $account > ./config/account
echo "Account: $account"

rm -rf "./gethdata/*" || true
cp -r "$account_dir" "./gethdata"

./init_geth.sh
