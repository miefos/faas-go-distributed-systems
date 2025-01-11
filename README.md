# FaaS

# Modules
## Auth service
User auth & registration

## Registry service
Registering & unregistering functions

## Execution service
Executing functions

# NATS
## Monitoring at localhost:8222
![alt text](images/nats-monitoring.png)

# Limitations
Since function activations are synchronous, it has certain limitations
* If message queue is long, then the HTTP request may timeout.
* If a function has a long execution, then the HTTP request may timeout and would affect other users as well. 