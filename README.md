# contacts-backend-go-echo

**contacts-backend-go-echo** is part of "contacts" project that is an initiative where we try to explore frontend and backend implementations in order to better understand it cutting-edge features. This repository presents a golang rest API sample.

## stack
* golang 1.22
* echo
* sqlite
* gorm

## features
* migrations
* high concurrent
* small footprint

## get started (linux only)

### option 1 - just build and use as docker image
build a local docker image
```
docker build --tag contacts.backend.go.echo .
```

execute the local docker image
```
docker run -p 8010:8010 contacts.backend.go.echo
```
### option 2 - execute from source code 

- first, install golang 1.22+, if you don't have it on your computer:  [how to install golang](https://go.dev/doc/install)
- go to root of solution and execute the commands below

set environment for development
```
cp .env_example .env
```

update deps
```
 go get -u -v
 go mod tidy
```

and then run the application
```
go run .
```

## testing
create a new contact
```
curl --request POST \
  --url http://localhost:8010/contacts \
  --header 'Content-Type: application/json' \
  --data '{
	"firstName": "Albert",
	"lastName": "Einstein",
	"phoneNumber": "2222-1111"
  }'
```

retrieve existing contacts
```
curl --request GET \
  --url http://localhost:8010/contacts
```
more examples and details about requests on (link) *to be defined