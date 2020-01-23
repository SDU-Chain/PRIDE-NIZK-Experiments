#!/usr/bin/env bash
find ./gethaccounts | grep UTC | awk -F '--' '{print $3}'
