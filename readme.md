# Info

Implementation of a go program out to solve the problem discussed in `task.md`.

* `./cmd/main.go` contains the entry point to the application
* `./api/...` contains an HTTP Handler and implements the RESTful related logic
* `./data/...` contains an interface definition of what the DB service is expected to provide. It also contains an in memory mock and a filesystem implementation of said interface
* `./helpers/...` contains any common utils for testing
* `./models/...` contains some code related to  modelling and validating the input JSON file format of the program

`go run cmd/main.go` will run the server, by default it will listen on http://localhost:8080

## cURL Examples

### Create an Appointment

#### Request:
```
curl \
  --location \
  --request POST 'localhost:8080/api' \
  --header 'Content-Type: application/json' \
  --data-raw '{
    "trainer_id": 1,
    "user_id": 1,
    "starts_at": "2020-01-24T09:00:00-08:00",
    "ends_at": "2020-01-24T09:30:00-08:00"
  }'
```

#### Response:
```
{
    "id": 1,
    "trainer_id": 1,
    "starts_at": "2020-01-24T09:00:00-08:00",
    "ends_at": "2020-01-24T09:30:00-08:00"
}
```

### Trainer's Scheduled Appointments

#### Request:
```
curl \
  --location \
  --request GET 'localhost:8080/api?trainer_id=1&filter=scheduled'
```

#### Response:
```
[
    {
        "id": 1,
        "trainer_id": 1,
        "starts_at": "2020-01-24T09:00:00-08:00",
        "ends_at": "2020-01-24T09:30:00-08:00"
    }
]
```

### Trainer's Available Appointments

Yet to be implemented!