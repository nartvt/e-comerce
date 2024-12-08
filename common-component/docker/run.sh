#!/bin/bash

status=$1

echo $status

docker_compose=docker-compose.yml

if [ "$status" = "start" ]; then
    sudo docker-compose -f $docker_compose  --project-name=init-database up -d --remove-orphans
fi

if [ "$status" = "stop" ]; then
    sudo docker-compose -f $docker_compose down
fi

