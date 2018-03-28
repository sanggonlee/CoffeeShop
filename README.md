# Coffee Shop

This is a toy RESTful API server written in Go

## Getting Started

First, you need to install goose.
```
go get -u github.com/pressly/goose/cmd/goose
```

Initialize the Postgres DB
```
make db_init
make db_migrate
```

### Running with a Docker image
Pull the docker image
```
docker pull sanggonlee/coffeeshop
```

Run the image in a container
```
docker run -d -p 8080:8080 --add-host=database:<YOUR-IP-ADDRESS> coffeeshop
```

Then the server will listen on port 8080. Good to go!

### Running natively
Put this directory under $GOPATH/src/github.com.
Go to the project root and run
```
go get -d .
```
to install dependencies.

Then
```
make run
```
to start the server.

## API

### Insert a drink
```
POST http://localhost:8080/drink
{
	"name": "strawberry mocha",
	"price": "5.40",
	"start": "2018-02-06T00:00:00+00:00",
	"end": "2018-02-24T00:00:00+00:00",
	"ingredients": [
		"milk",
		"cocoa",
		"espresso",
		"cream",
		"strawberry"
	]
}
```

### Delete a drink
```
DELETE http://localhost:8080/drink?id=deffdea5-85ac-45c0-a859-261855b19544
```

### Query drinks
The following query parameters are supported:
**name** - Name of the drink to search for.  
**date** - Filter the drinks by the date they're available on. Only takes the RFC3339 format.  
**ingredients** - List of ingredients separated by commas, matches the drinks that contain all of the given ingredients.  
**offset** - Used for pagination. Starting offset to fetch the results.  
**limit** - Used for pagination. The size of the batch query request.  

The response has the following structure:
```
{
  "offset_previous": The offset given by the request (defaults to 0),
  "offset_current": The current offset; offset_previous + number of hits returned
  "hits": The matched results
}
```

Example:
Query drinks with offset 0 and limit 5, available on January 4th of 2018 and has ingredient milk
```
GET http://localhost:8080/drinks?ingredients=milk&name=black%20mocha&offset=1&date=2018-01-04T00%3A00%3A00%2B00%3A00
```

Response:
```
{
  "offset_previous": 0
  "offset_current": 5,
  "hits": [
    {
      "id": "f3221e11-bc39-46ca-b514-abbd4349c4ba",
      "name": "black mocha",
      "price": "5.10",
      "start": "2018-01-02T00:00:00Z",
      "end": "2018-02-22T00:00:00Z",
      "ingredients": [
        "cocoa",
        "cream",
        "espresso",
        "milk"
      ]
    },
    {
      "id": "60d04dc9-ef0f-44bc-9b45-9570465a1e33",
      "name": "strawberry mocha",
      "price": "5.40",
      "start": "2018-01-01T00:00:00Z",
      "end": "2018-01-31T00:00:00Z",
      "ingredients": [
        "cocoa",
        "cream",
        "espresso",
        "milk",
        "strawberry"
      ]
    }
  ]
}
```