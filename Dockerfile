FROM golang:1.23.1-alpine3.20

RUN go version

COPY ./ ./

RUN go mod download

RUN go build -o ./cmd/main/main.go ./.bin/main

CMD ["./.bin/main"]
