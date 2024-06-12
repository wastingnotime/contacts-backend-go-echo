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

## build and use
build a local docker image
```
docker build --tag contacts.backend.go.echo .
```

execute the local docker image
```
docker run -p 8010:80 contacts.backend.go.echo
```

create a new contact
```
curl --request POST \
  --url http://localhost/contacts \
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
  --url http://localhost/contacts
```
more examples and details about requests on (link)


## development 

install golang 1.22+ [(how to install golang)](https://go.dev/doc/install)

### update
```
 go get -u -v
 go mod tidy
```

### execute

run
```
go run .
```

### testing
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
