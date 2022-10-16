#!/bin/bash

export COMPOSE_IGNORE_ORPHANS=True

# postgresql
export POSTGRESQL_IMAGE=gokomodo-postgresql
export POSTGRESQL_IMAGE_TAG=production
export POSTGRESQL_CONTAINER=gokomodo-postgresql-production
export POSTGRESQL_HOST=gokomodo-postgresql.service
export POSTGRESQL_USER=gokomododev
export POSTGRESQL_PASSWORD=inisecret
export POSTGRESQL_DB=gokomodo
export POSTGRESQL_TIME_ZONE=Asia/Kuala_Lumpur
docker build -t "$POSTGRESQL_IMAGE:$POSTGRESQL_IMAGE_TAG" -f ./manifest-docker/Dockerfile.postgresql ./manifest-docker

# redis
export REDIS_IMAGE=gokomodo-redis
export REDIS_IMAGE_TAG=production
export REDIS_CONTAINER=gokomodo-redis-production
export REDIS_HOST=gokomodo-redis.service
docker build -t "$REDIS_IMAGE:$REDIS_IMAGE_TAG" -f ./manifest-docker/Dockerfile.redis ./manifest-docker

# pgadmin
export PGADMIN_IMAGE=gokomodo-pgadmin
export PGADMIN_IMAGE_TAG=production
export PGADMIN_CONTAINER=gokomodo-pgadmin-production
export PGADMIN_HOST=gokomodo-pgadmin.service
export PGADMIN_EMAIL=admin@gokomodo.com
export PGADMIN_PASSWORD=inisecret
docker build -t "$PGADMIN_IMAGE:$PGADMIN_IMAGE_TAG" -f ./manifest-docker/Dockerfile.pgadmin ./manifest-docker

docker-compose -f ./manifest/docker-compose.production.yaml up -d --build
