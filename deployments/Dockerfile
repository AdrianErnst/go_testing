ARG PROJECT_NAME="k8s_go_client"
ARG PORT=9292

############################
# STEP 1 build executable binary
############################
FROM golang@sha256:d722397d1dfe5b449d83922fe55cc93dcca3b8c71a0a6c218acb218f1df1c984 AS builder
ARG PORT

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser.
ENV USER=appuser
ENV UID=10001 

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .

# Fetch dependencies.
# Using go mod with go 1.11
RUN go mod download
RUN go mod verify
# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o  /go/bin/k8s_go_client ./cmd/web/

############################
# STEP 2 build a small image
############################
FROM scratch
ARG PORT
ENV GIN_MODE=release

# Import the user and group files from the builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable.
COPY --from=builder /go/bin/k8s_go_client /go/bin/k8s_go_client

# Use an unprivileged user.
USER appuser:appuser

# Port on which the service will be exposed.
EXPOSE $PORT

# Run the binary.
ENTRYPOINT ["/go/bin/k8s_go_client"]