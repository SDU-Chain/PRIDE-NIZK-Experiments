#!/bin/bash
set -e

./kill_http_server.sh

cd ./http
screen -Sdm pride_exp_http python3 -m http.server 8000
