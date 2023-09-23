FROM alpine
ARG PORT

# Port on which the service will be exposed.
EXPOSE $PORT

WORKDIR /app
ADD build build
ENTRYPOINT build/go_k8s_client