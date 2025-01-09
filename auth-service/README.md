# FaaS (Function as a Service)
A simple Go-based project using Gorilla Mux for routing, JWT for authentication, and bcrypt for password hashing.

* NATS KV as database.

* To start the project, run
`docker-compose up`

* If needed, download dependencies
Download required dependencies:

```
go mod download github.com/golang-jwt/jwt/v4
go mod download github.com/gorilla/mux
go mod download github.com/nats-io/nats.go
```

# Test the API
## Registration with username validation
`curl -X POST http://127.0.0.1:8080/register -H "Content-Type: application/json" -d '{"username":"your_username", "password":"your_password"}'`

![alt text](images/console-registration.png)