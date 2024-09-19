FROM golang:1.23.1-alpine3.20 AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc  gettext musl-dev

RUN go version

COPY . .

RUN go mod download

RUN go build -o .bin/main cmd/main/main.go

FROM alpine:3.20


COPY --from=builder /usr/local/src/.bin/main .

COPY config/config.yaml config/config.yaml
EXPOSE 5000

CMD ["./main"]
