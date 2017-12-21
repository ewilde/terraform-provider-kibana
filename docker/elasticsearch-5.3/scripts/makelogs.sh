#!/usr/bin/env bash
while ! curl --output /dev/null --silent --head --fail http://localhost:9200; do sleep 1 && echo -n .; done;

/usr/bin/makelogs
