#!/bin/bash
docker network create gokomodo-environment-development
echo "====== CREATE DB ======"
cd database
make dev
cd ..
echo "====== MIGRATE DB ======"
cd migration
make dev
cd ..
echo "====== RUNNING API ======"
cd api
make dev
echo "====== RUNNING TEST API ======"
make test
