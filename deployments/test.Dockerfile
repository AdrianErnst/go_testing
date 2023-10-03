FROM golang@sha256:d722397d1dfe5b449d83922fe55cc93dcca3b8c71a0a6c218acb218f1df1c984

WORKDIR /app

COPY . .
RUN go mod download
RUN go mod verify
RUN go get golang.org/x/tools/cmd/cover

CMD [ "sh", "-c", "go test -v ./... -tags integration -coverpkg=./..." ]