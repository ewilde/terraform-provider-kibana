#!/usr/bin/env bash
docker run -d --name elasticsearch -p 9200:9200 elastic-local:6.0.0
echo "waiting for 15 seconds for elastic to finish starting up"
sleep 15
docker run --name kibana --link elasticsearch -p 5601:5601 docker.elastic.co/kibana/kibana-oss:6.0.0
