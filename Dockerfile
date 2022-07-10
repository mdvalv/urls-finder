FROM golang:1.18 AS build-env
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.mod /src/
RUN go mod download
COPY . .
RUN  go build -a -ldflags="-s -w" -gcflags="all=-trimpath=/src" -asmflags="all=-trimpath=/src"

FROM alpine:latest

RUN apk add --no-cache ca-certificates \
    && rm -rf /var/cache/*

RUN mkdir -p /app \
    && adduser -D finder \
    && chown -R finder:finder /app

USER finder
WORKDIR /app

COPY --from=build-env /src/urls-finder .

ENTRYPOINT [ "./urls-finder" ]
