# syntax=docker/dockerfile:1

## UI
FROM node:19-alpine as ui

RUN npm i -g pnpm

WORKDIR /src/ui
COPY ui/package.json ui/pnpm-lock.yaml /src/ui/

RUN pnpm install --frozen-lockfile

WORKDIR /src
COPY ./ui/ ./ui/

WORKDIR /src/ui
RUN pnpm run build


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
COPY --from=ui /src/ui/dist ./public

EXPOSE 3000

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/server"]