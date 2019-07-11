#!/bin/bash
set -e

read -p "How many car accounts do you want to generate? " count

rm -rf ./gethaccounts/cars/*||true

START=1
END=$count

echo "Generating $count car accounts..."

i=$START
while [[ $i -le $END ]]
do
	echo -n "($i/$END)"
	mkdir -p ./gethaccounts/cars/$i
	./bin/geth-timing --datadir ./gethaccounts/cars/$i --password ./password/password.txt account new > /dev/null 2>&1
	printf '\r'
	((i = i + 1))
done

echo "Creating tar archive..."
mkdir -p ./http/gethaccounts
cd ./gethaccounts
tar -zcf ../http/gethaccounts/cars.tar.gz cars
