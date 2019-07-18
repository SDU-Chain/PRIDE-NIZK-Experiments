#!/bin/bash
set -e


arr[0]=1000;
arr[1]=2000;
arr[2]=5000;
arr[3]=10000;
arr[4]=20000;
arr[5]=50000;
#rand=$[$RANDOM % ${#arr[@]}]
arr=( $(shuf -e "${arr[@]}") );

echo "You have 10 seconds to record the order of number of commitments!"
printf '%s\n' "${arr[@]}"
sleep 10

for i in "${arr[@]}"
do
	echo ----$i----;
	echo $count | ./run_car_emulator.sh ;
	echo ----------;
done

