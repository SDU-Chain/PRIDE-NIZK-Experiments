#!/usr/bin/env bash
read -p "Number of cars=? " count
for i in $( seq 1 $count )
do
	echo 100 | ./run_car_emulator.sh &
done
