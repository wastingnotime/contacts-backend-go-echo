FROM golang:1.24-alpine as build-stage

WORKDIR /app

COPY . /app

# to supports cgo as sqlite/driver dependency
RUN apk add --update gcc musl-dev

RUN go mod download

RUN go build -o app



FROM alpine:latest as deploy-stage

ENV DB_LOCATION=/data/contacts.db
ENV ENVIRONMENT=production

# act as doc only
EXPOSE 8010

# principle of least privilege
RUN addgroup -S nonroot && adduser -S appuser -G nonroot

# just to allow appuser to have access to the data volume
WORKDIR /data
RUN chown appuser .

VOLUME /data

WORKDIR /app
RUN chown appuser .

COPY --from=build-stage /app/app .

USER appuser

ENTRYPOINT ./app