#!/bin/bash

status=$1

echo $status

docker_compose=docker-compose.yaml

if [ "$status" = "up" ]; then
    sudo docker-compose -f $docker_compose --env-file=.env --project-name=init-database up -d --remove-orphans
fi

if [ "$status" = "down" ]; then
    sudo docker-compose -f $docker_compose --env-file=.env --project-name=init-database down
fi

