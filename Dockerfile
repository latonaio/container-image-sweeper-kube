# syntax = docker/dockerfile:experimental
# Build Container
FROM golang:1.17.6 as builder

ENV GO111MODULE on
WORKDIR /go/src/github.com/latonaio

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o container-image-sweeper ./cmd

# Runtime Container
FROM alpine:3.15

RUN apk add --no-cache libc6-compat

ENV SERVICE=container-image-sweeper \
    POSITION=Runtime \
    AION_HOME="/var/lib/aion"

ENV APP_DIR="${AION_HOME}/${POSITION}/${SERVICE}"

WORKDIR ${AION_HOME}

COPY --from=builder /go/src/github.com/latonaio/container-image-sweeper .

CMD ["./container-image-sweeper"]
