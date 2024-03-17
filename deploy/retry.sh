#!/bin/bash
echo "Stopping containers, removing images and retrying..."
docker-compose down && docker image rm deploy-filmoteka && docker-compose up -d