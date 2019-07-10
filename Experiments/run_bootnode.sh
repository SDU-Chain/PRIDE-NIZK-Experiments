#!/bin/bash

# stop when error occurs
set -e

# kill server process
screen -S pride_bootnode -X quit >/dev/null 2>&1 || true
killall bootnode >/dev/null 2>&1 || true

# run server at background
screen -Sdm pride_bootnode ./bin/bootnode -nodekey ./bootnode.key -verbosity 9 -addr :30310
