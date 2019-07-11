#!/bin/bash

# stop when error occurs
set -e

# kill server process gently
screen -S pride_exp_http -X stuff $'\003' >/dev/null 2>&1 || true

