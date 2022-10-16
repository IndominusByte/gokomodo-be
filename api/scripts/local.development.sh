#!/bin/bash

export COMPOSE_IGNORE_ORPHANS=True

export BACKEND_IMAGE=gokomodo-go
export BACKEND_IMAGE_TAG=development
export BACKEND_CONTAINER=gokomodo-go-development
export BACKEND_HOST=gokomodo-go.service
export BACKEND_STAGE=development

docker build -t "$BACKEND_IMAGE:$BACKEND_IMAGE_TAG" .
docker-compose -f ./manifest/docker-compose.development.yaml up -d --build
