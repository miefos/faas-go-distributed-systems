# FaaS Project

## Description
### Structure
### Scaling
### Security

## Launching options
### Preparation
The whole project can be automatically built and deployed using the `./bin/run.sh` script on UNIX machines or `.\bin\run.windows.ps1` on Windows machines. The script will build all the needed docker images and then start the docker compose process, alternatively, you can build and deploy the services manually.

First build the docker images:
```bash
docker build -t apisix-configurator:latest ./api-gateway
docker build -t auth-service:latest ./auth-service
docker build -t registry-service:latest ./registry-service
docker build -t spawner-service:latest ./spawner-service
docker build -t publisher-service:latest ./publisher-service
```

### Docker compose deployment
To deploy the services using docker compose, simply run the following command:
```bash
docker compose up
```
 To manually scale services in docker compose, use the `--scale` option:
```bash
docker compose up --scale <service-name>=<number>
```
### Swarm deployment
To use swarm, first initialize the swarm, then deploy the stack:
```bash
docker swarm init
docker stack deploy -c docker-compose.yml faas
```

To manually scale services in swarm, use the `docker service scale` command:
```bash
docker service scale faas_auth-service=<number> faas_registry-service=<number> faas_execution-service=<number>
```

To remove the swarm stack, use the `docker stack rm` command:
```bash
docker stack rm faas
```

## Modules
- The auth service is available INTERNALLY at http://auth-service:8081, EXTERNALLY at http://localhost/auth.
- The registry service is available INTERNALLY at http://registry-service:8082, EXTERNALLY at http://localhost/registry.
- The publisher service is available INTERNALLY at http://publisher-service:8083, EXTERNALLY at http://localhost/publisher.

### API Gateway / APISIX
The entrypoint to all the services, this is the only service that is exposed to the outside world.

### Auth service
User auth & registration

Details can be found here: [Auth service](auth-service/README.md)

### Registry service
Registering & unregistering functions

Details can be found here: [Registry service](registry-service/README.md)

### Publisher service
It exposes a REST API to spawn the worker that will execute the desired function.

Details can be found here: [Publisher service](publisher-service/README.md)

### Spawner service
It executes functions by spawning workers as containers from an image reference and a string argument.

Details can be found here: [Spawner service](spawner-service/README.md)

## Execution Example
### User Registration
```bash
curl -X POST http://localhost/auth/register -d '{"username":"user","password":"password"}'
```
This will return a 200 OK status code if the registration was successful and the user uuid.

### User Login
```bash
curl -X POST http://localhost/auth/login -d '{"username":"user","password":"password"}'
```
This will return an authorization token that can be used to authenticate the user in the other services.

### Function Image Creation
You can either create a new image locally or use an existing image from Docker Hub.
#### Create Dockerfile for the function
```Dockerfile
FROM python:3.9-slim
WORKDIR /usr/src/app
COPY simple_function.py .
RUN chmod +x simple_function.py
ENTRYPOINT ["python", "simple_function.py"]
```

#### Create the function file
```python
import sys

def simple_function(input_string):
    return input_string.upper()

if __name__ == "__main__":
    input_string = sys.argv[1]
    result = simple_function(input_string)
    print(result)
```

#### Build the image
```bash
docker build -t simple_function:latest .
```

#### Push the image to Docker Hub (optional)
```bash
docker tag simple_function:latest <docker-hub-username>/simple_function:latest
```

### Function Registration
```bash
curl -X POST http://localhost/registry/register -d '{
    "name":"HumanFriendlyName",
    "description":"ServiceDescription",
    "payload":"docker-image-reference",
    }' -H "Authorization: Bearer <token>"
```

### Function Execution
```bash
curl -X POST http://localhost/publisher/publish -d '{
    "name":"HumanFriendlyName",
    "argument":"string-argument"
    }' -H "Authorization Bearer <token>"
```

To test other features such as function CRUD operations, please refer to the README.md files in the respective services.

## Replication
Since we did not use kubernetes, there was no classic standard way to dinamically replicate the services. However, we can use the `docker-compose` or `docker swarm` scaling options to replicate the services.

As of now, the services are statically programmed to be launched with healthchecks and automatic reloads in case of failure. This is done by the `docker-compose.yml`.

The static scaling has been created with this numbers:
- 3 auth-service
- 3 registry-service
- 3 publisher-service
- 10 spawner-service
- 1 apisix-configurator (that dies after the configuration is done)
- 2 nats

And one of all accessories services for APISIX.

Moreover, each service is programmed to be able to switch to another NATS server in case of failure.
The NATS servers are configured in cluster mode with a replication factor of 2.