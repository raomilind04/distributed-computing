#!/bin/bash

CHUNKS="$1"

echo "CHUNKS=$CHUNKS"
export CHUNKS 
python3 chunking-script.py

if [ ! -d "data" ]; then
    echo "Python script failed."
    exit 1
fi

if [ -d "pi_service/data" ] && [ "$(ls -A pi_service/data)" ]; then
    echo "'pi_service/data' directory exists."
    rm -rf pi_service/data/*
fi

echo "Move data directory"
mv data pi_service/

if [ ! -d "pi_service/data" ]; then
    echo "Failed to move"
    exit 1
fi

DOCKER_COMPOSE_PATH="../dockerCompose/$CHUNKS/docker-compose.yml"

if [ -f "$DOCKER_COMPOSE_PATH" ]; then
    echo "Moving $DOCKER_COMPOSE_PATH"
    if [ -f "docker-compose.yml" ]; then
        echo "Replacing"
        rm -f docker-compose.yml
    fi
    cp "$DOCKER_COMPOSE_PATH" .
else
    echo "No docker-compose.yml file found at $DOCKER_COMPOSE_PATH"
fi

echo "Start Docker"
docker-compose up --build

if [ $? -ne 0 ]; then
    echo "Docker Compose failed"
    exit 1
fi

echo "Completed."

