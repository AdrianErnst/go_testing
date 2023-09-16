#!/bin/bash

(
    PROJECT_NAME=go_k8s_client
    docker build --build-arg PORT=9292 --no-cache --progress=plain --tag "${PROJECT_NAME}" .
)
