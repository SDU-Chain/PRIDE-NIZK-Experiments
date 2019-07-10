#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import json
import sys

if __name__ == "__main__":
    if len(sys.argv) <= 1:
        sys.stderr.write(
            '''Usage: parse_json.py <key_1> <key_2> ... <key_n>

-- example bash command --
cat foo.json | parse_json.py \\\"bar\\\" \\\"baz\\\" 1
------------------------------

-- content of foo.json ---
{"bar":{"baz":["zero","one"]}}
--------------------------

--------- output ---------
one
--------------------------
''')
        sys.exit(1)
    else:

        json_str = ""
        for line in sys.stdin:
            json_str += line
        obj = json.loads(json_str)

        for i in range(1, len(sys.argv)):
            obj = obj[json.loads(sys.argv[i])]

        print(obj)

