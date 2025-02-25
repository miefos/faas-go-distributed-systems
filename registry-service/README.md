# Registry Service
This service is responsible for managing the functions registered by the users.

## Environment Variables
The following environment variables are required to run the service:
- `NATS_URL`: The URL of the NATS server.
- `SERVER_ADDRESS`: The port on which the server will listen.
- `BUCKET_NAME`: The name of the bucket in which the functions will be stored.

## REST API
The REST API to the registry service is described below.

| Endpoint             | Method | Description                                   |
|----------------------|--------|-----------------------------------------------|
| registry/register    | POST   | Register a new function for the user.         |
| registry/retrieve    | GET    | Retrieve a specific function by ID            |
| registry/delete      | DELETE | Delete a specific function by ID              |
| registry/list        | GET    | List all the functions registered by the user |

Here's an example of the body of the request to register a new function:
```json
{
    "name":"MyFunction",
    "description":"Test function",
    "payload":"myrepo/myfunction:v1"
}
```

Here's an example of the body of the request to retrieve a function:
```json
{
    "name":"MyFunction"
}
```

Here's an example of the body of the request to delete a function:
```json
{
    "name":"MyFunction"
}
```

Here's an example of the body of the request to list all the functions:
```json
{}
```