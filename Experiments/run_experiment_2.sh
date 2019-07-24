#!/usr/bin/env bash
read -p "Number of cars=? " count
last_pid=0
for i in $( seq 1 $count )
do
	while true; do
		echo 100 | ./run_car_emulator.sh &
		this_pid=%?
		if [$this_pid -ne $last_pid]; then
			last_pid=this_pid
			break
		fi
	done

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
