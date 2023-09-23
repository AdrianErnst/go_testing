#!/bin/bash

source kind_build.sh
helm uninstall client
helm upgrade client ./helm --install
sleep 20
kubectl logs -f svc/go-k8s-client