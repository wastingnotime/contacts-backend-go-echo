FROM golang:1.24-alpine
WORKDIR /app

COPY . /app

RUN apk add --update gcc musl-dev

RUN go build -o app

VOLUME data

ENV DB_LOCATION=/data/contacts.db
ENV ENVIRONMENT=production

ENTRYPOINT ./app