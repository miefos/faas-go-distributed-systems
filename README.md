# FaaS Project

## Description

## Launching options
To launch the project with docker compose, simply run the following command:
```bash
docker compose up
```

To scale services, simply add the `--scale` flag to the command:
```bash
docker compose up --scale auth-service=<number> --scale registry-service=<number> --scale execution-service=<number>
```

## Modules

## APISIX
I have no freaking clue of what this does or how it works, we'll explore it.

### API Gateway
The entrypoint to all the services, this is the only service that is exposed to the outside world.

- [ ] Define the API functionalities
- [ ] Implement rerouting to Auth service
- [ ] Implement rerouting to Registry service
- [ ] Implement connection to NATS Messaging service

### Auth service
User auth & registration
- [x] User registration
- [ ] Add user id to user model
- [ ] Login

Details can be found here: [Auth service](auth-service/README.md)

### Registry service
Registering & unregistering functions
- [x] Register function
- [x] Unregister function
- [x] Get function by id
- [x] Get all functions for user
- [x] Update function by id

Details can be found here: [Registry service](registry-service/README.md)


### Execution service
Executing functions