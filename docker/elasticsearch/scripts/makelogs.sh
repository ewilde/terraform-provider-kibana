#!/usr/bin/env bash
echo waiting for elastic to start
while ! curl --output /dev/null --silent --head --fail http://localhost:9200 -u elastic:changeme; do sleep 1 && echo -n .; done;
echo waiting elastic started
/usr/bin/makelogs --auth "elastic:changeme"
echo finished creating logs
