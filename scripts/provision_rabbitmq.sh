#!/bin/bash

DOCKER=docker
DOCKER_PATH=$(command -v $DOCKER)

if [ "$CI_COMMIT_REF_NAME" = "" ]; then
    echo "Starting RabbitMQ..."

    if [ -z "$DOCKER_PATH" ]; then
        echo "Docker not found. Checking for podman."

        DOCKER=podman
        PODMAN_PATH=$(command -v podman)
        if [ -z "$PODMAN_PATH" ]; then
            echo "Podman not found. No suitable container engine."
            exit 127
        fi
    fi

    $DOCKER run --name my_rabbitmq -p 5672:5672 -p 15672:15672 -e RABBITMQ_DEFAULT_USER=user -e RABBITMQ_DEFAULT_PASS=1234567 -d rabbitmq:3-management



    printf "\nWaiting for RabbitMQ to be fully available."
    while ! curl -s -u user:1234567 http://localhost:15672/api/overview > /dev/null; do
        echo "Waiting for RabbitMQ..."
        sleep 5
    done

    echo "\n RabbitMQ ready!"
else
    echo "No RabbitMQ"
fi

echo "Provisioning psql..."

sh ./scripts/provision_psql.sh