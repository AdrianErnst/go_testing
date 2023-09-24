#!/bin/bash

source kind_build.sh
kubectl delete ns client
helm uninstall client
helm upgrade client ./../deployments/helm --install
sleep 20
kubectl logs -f svc/go-k8s-client -n client