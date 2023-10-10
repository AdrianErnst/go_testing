# First Go Repo
This repo builds the base for future go projects

It provides the following features:
* Openapi generated gin web server with k8s client and open api html endpoints
* Kubernetes Client with in-cluster and local client configs + integration tests for both setups
* Test split by build flags
* Environment using startup flags
* Pre-Commit Hooks
  * Linting
  * Gitleaks secret leak checks
* Prettier config
* Gitignore config
* Yaml based kind cluster and registry using ctlptl included in tilt up/down
* Dockerfile with image builder with high security, test Dockerfile, local Dockerfile
* Tilt with live update and testing for unit and integration tests in and out of cluster

