#!/bin/bash

export COMPOSE_IGNORE_ORPHANS=True

# db migration
export DB_MIGRATION_IMAGE=gokomodo-db-migration
export DB_MIGRATION_IMAGE_TAG=development
export DB_MIGRATION_CONTAINER=gokomodo-db-migration-development
export DB_MIGRATION_HOST=gokomodo-db-migration.service
export BACKEND_STAGE=development

docker build -t "$DB_MIGRATION_IMAGE:$DB_MIGRATION_IMAGE_TAG" .
docker-compose -f ./manifest/docker-compose.development.yaml up -d --build
