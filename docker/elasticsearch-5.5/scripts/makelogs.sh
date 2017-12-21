#!/usr/bin/env bash
while ! curl --output /dev/null --silent --head --fail http://localhost:9200 -u elastic:changeme; do sleep 1 && echo "waiting for elastic to start..."; done;

echo " *** will create logs *** "
/usr/bin/makelogs --auth "elastic:changeme"
echo " *** finished create logs *** "
