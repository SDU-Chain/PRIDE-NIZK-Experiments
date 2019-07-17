#!/bin/bash
set -e
read -p "Number of commitments? " count
for i in {1..10}
do
echo ----$i----;
echo $count | ./run_car_emulator.sh
echo ----------;
done

