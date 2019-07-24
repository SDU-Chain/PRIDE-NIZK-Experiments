#!/usr/bin/env bash
read -p "Number of cars=? " count
for i in $( seq 1 $count )
do
	echo 100 | ./run_car_emulator.sh &
done

while true; do
	num=`pidof car | tr " " "\n" | wc -l`
	echo "Number of cars remaining: $num"
	if [ $num -eq 0 ]; then
		echo "No cars remaining. Exited."
		exit 0
	fi
	sleep 1
done
