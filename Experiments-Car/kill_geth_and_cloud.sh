#!/bin/bash

# stop when error occurs
set -e

# kill server process gently
screen -S pride_exp_geth -X stuff $'\003' >/dev/null 2>&1 || true
screen -S pride_exp_cloud -X stuff $'\003' >/dev/null 2>&1 || true
# screen -S pride_exp_cloud -X exit >/dev/null 2>&1 || true

