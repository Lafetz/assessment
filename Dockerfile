FROM golang:1.22.3-alpine3.18 AS builder

ENV APP_HOME=/go/src/web

WORKDIR "${APP_HOME}"

COPY ./go.mod ./go.sum ./

RUN go mod download
RUN go mod verify

COPY ./internal ./internal
COPY ./cmd ./cmd

RUN go build -o ./bin/web ./cmd

FROM alpine:latest

ENV APP_HOME=/go/src/web
RUN mkdir -p "${APP_HOME}"

WORKDIR "${APP_HOME}"

COPY --from=builder "${APP_HOME}"/bin/web "${APP_HOME}"

ENV PORT=8080

EXPOSE ${PORT}

CMD ["./web"]
