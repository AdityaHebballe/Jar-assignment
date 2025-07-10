#!/bin/bash
set -e

sudo docker build -t greeter-service:latest services/greeter-service
sudo docker tag greeter-service:latest adityahebballe/greeter-service:latest
sudo docker push adityahebballe/greeter-service:latest

sudo docker build -t echo-service:latest services/echo-service
sudo docker tag echo-service:latest adityahebballe/echo-service:latest
sudo docker push adityahebballe/echo-service:latest

sudo docker build -t admin-service:latest services/admin-service
sudo docker tag admin-service:latest adityahebballe/admin-service:latest
sudo docker push adityahebballe/admin-service:latest

echo "Docker images built and pushed to Docker Hub successfully!"
