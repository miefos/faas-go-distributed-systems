# Registry Service
This service is responsible for managing the functions registered by the users.

## Environment Variables
The following environment variables are required to run the service:
- `NATS_URL`: The URL of the NATS server.
- `SERVER_ADDRESS`: The port on which the server will listen.
- `BUCKET_NAME`: The name of the bucket in which the functions will be stored.

## REST API
The REST API to the registry service is described below. Note that all calls should contain a `UserId` header with the user's ID.

| Endpoint       | Method | Description                                   |
|----------------|--------|-----------------------------------------------|
| /register      | POST   | Register a new function for the user.         |
| /retrieve/:id  | GET    | Retrieve a specific function by ID            |
| /update/:id    | PUT    | Update a specific function by ID              |
| /delete/:id    | DELETE | Delete a specific function by ID              |
| /list          | GET    | List all the functions registered by the user |

Here's an example of the body of the request to register a new function:
```json
{
    "id":"func1",
    "name":"MyFunction",
    "description":"Test function",
    "payload":"print(42)"
}
```