#!/bin/bash

mode=${1:-"dev"}

docker compose -f docker/docker-compose."${mode}".yml down -v

echo "Removing certs..."

sudo rm -rf certs/*
sudo rm -rf secrets/*
sudo rm -rf docker/.env

docker rmi benlun1201/veg-store-backend-builder

echo "Done."