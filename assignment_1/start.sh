#!/bin/bash

CHUNKS="$1"  

echo "Chunk"
CHUNKS=$CHUNKS python3 chunking-script.py

if [ ! -d "data" ]; then
    echo "Python script failed."
    exit 1
fi

if [ -d "pi_service/data" ] && [ "$(ls -A pi_service/data)" ]; then
    echo "'pi_service/data' directory exists."
    rm -rf pi_service/data/*
fi

echo "Move"
mv data pi_service/

if [ ! -d "pi_service/data" ]; then
    echo "Failed to move"
    exit 1
fi

echo "Start Docker"
docker-compose up --build

if [ $? -ne 0 ]; then
    echo "Docker Compose failed"
    exit 1
fi

echo "Completed."


