#!/usr/bin/env bash
read -p "Number of cars=? " count
for i in $( seq 1 $count )
do
	echo 100 | ./run_car_emulator.sh 1>/dev/null 2>/dev/null &
	echo -n "Starting car ($i/$count)"
	printf "\r"
done

while true; do
	num=`pidof car | tr " " "\n" | wc -l`
	echo -n "Number of cars remaining: $num"
	printf "\r"
	if [ $num -eq 0 ]; then
		echo "No cars remaining. Exited."
		exit 0
	fi
	sleep 1
done
