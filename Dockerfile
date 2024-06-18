# syntax=docker/dockerfile:1

## Build
FROM golang:1-alpine AS build

WORKDIR /src

COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o server cmd/server/main.go


## Deploy
FROM alpine

RUN apk add --no-cache tini ca-certificates

WORKDIR /
COPY --from=build /src/server server

EXPOSE 3000

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/server"]