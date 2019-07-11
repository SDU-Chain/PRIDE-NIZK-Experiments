#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import json
import os

if __name__ == '__main__':
    genesis = None
    # read template
    with open('./genesis-template.json', 'r') as f:
        genesis = json.loads(f.read())
        f.close()

    # find all car's address
    car_addresses = list()
    for root, dirs, files in os.walk("./gethaccounts/cars/"):
        for file in files:
            try:
                path = os.path.join(root, file)
                with open(path, 'r') as f:
                    t = json.loads(f.read())
                    car_addresses.append(t['address'])
            except:
                pass

    # fund all the cars
    for car_addr in car_addresses:
        genesis['alloc'][car_addr] = {"balance": "0x100000000000000000000000000000000000000000000000000000000000000"}

    # output
    with open('./genesis.json', 'w') as f:
        f.write(json.dumps(genesis))
        f.flush()
        f.close()

    print("Generated genesis.json file with", len(car_addresses), "car accounts pre-funded.")

