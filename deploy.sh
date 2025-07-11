#!/bin/bash
set -e

# 3 microservices
kubectl apply -f kubernetes/greeter-service/
kubectl apply -f kubernetes/echo-service/
kubectl apply -f kubernetes/admin-service/


kubectl apply -f kubernetes/access-control.yaml
kubectl apply -f kubernetes/app-secrets.yaml
kubectl apply -f kubernetes/dummy-volume.yaml
kubectl apply -f kubernetes/ingress.yaml
kubectl apply -f kubernetes/grafana-ingress.yaml

echo "Kubernetes manifests applied successfully!"
