#!/bin/bash
set -e

rm -rf ./gethaccounts/cars/*||true
for i in {1..1000}
do
	printf '\r'
	echo -n "($i/1000)"
	mkdir -p ./gethaccounts/cars/$i
	./bin/geth-timing --datadir ./gethaccounts/cars/$i --password ./password/password.txt account new > /dev/null 2>&1;
done

