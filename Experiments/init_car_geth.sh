#!/bin/bash
set -e
read -p "Car account id=?" carid
echo "./gethaccounts/cars/$carid" |./init_geth_at.sh
