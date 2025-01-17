# Publisher service
This service will listen for REST calls and will act accordingly to spawn the worker that will execute the desired function.

## REST Interface
It exposes the following endpoint:
- `POST /publish`

It expects the following JSON payload:
```json
{
    "uui": "random-uuid",
    "name": "NameOfTheFunction",
    "argument": "string-argument"
}
```

## NATS Interface
After recieving the REST call, it will first get from the registry service the image reference name, then it will perform a request to the NATS queue where the spawner will be listening.

When the spawner receives the message, it will spawn a worker container with the image reference and the string argument.

When done, the spawner will return the result to the publisher service, which will then return the result to the caller.