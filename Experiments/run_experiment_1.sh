#!/bin/bash
set -e
for i in {1..10}
do
echo ----$i----;
echo $i | ./run_car_emulator.sh
echo ----------;
done

