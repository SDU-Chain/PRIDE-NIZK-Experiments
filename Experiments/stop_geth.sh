#!/bin/bash

# stop when error occurs
set -e

# kill server process
screen -S pride_exp_geth -X quit >/dev/null 2>&1 || true
screen -S pride_exp_cloud -X quit >/dev/null 2>&1 || true
killall geth-timing >/dev/null 2>&1 || true
killall cloudProvider >/dev/null 2>&1 || true

