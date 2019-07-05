#!/bin/bash
set -e
for i in {0..1000..40}
do
echo ----$i----;
echo $i | ./run.sh
echo ----------;
done

