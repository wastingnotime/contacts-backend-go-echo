FROM golang:alpine
WORKDIR /app
ADD . /app
RUN cd /app && go build -o app

ENV ENVIRONMENT=production

ENTRYPOINT ./app