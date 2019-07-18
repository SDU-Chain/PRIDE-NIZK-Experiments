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

echo "Write the order of number of commitments to ./output/order.txt"
printf '%s\n' "${arr[@]}" |tee ./output/order.txt

for i in "${arr[@]}"
do
	echo ----$i----;
	echo $i | ./run_car_emulator.sh ;
	echo ----------;
done

