#!/bin/bash

ADMIN_API_KEY="edd1c9f034335f136f87ad84b625c8f1"
JWT_SECRET_KEY=$(<secret.txt)
JWT_KEY="faas-app-key"

# Wait until apisix container is healthy
until curl -sf http://apisix:9180/apisix/admin/routes -H "X-API-KEY: $ADMIN_API_KEY" > /dev/null; do
    echo "Waiting for apisix service to be ready..."
    sleep 1
done

curl -X PUT http://apisix:9180/apisix/admin/consumers \
-H "X-API-KEY: $ADMIN_API_KEY" \
-d '{
    "username": "jwt_consumer",
    "plugins": {
        "jwt-auth": {
            "key": "'$JWT_KEY'",
            "secret": "'$JWT_SECRET_KEY'"
        }
    }
}'

# Define a route for auth-service
curl -X PUT http://apisix:9180/apisix/admin/routes/1 \
  -H "X-API-KEY: $ADMIN_API_KEY" \
  -d '{
        "uri": "/auth/*",
        "methods": ["GET", "POST", "PUT", "DELETE"],
        "upstream": {
          "type": "roundrobin",
          "nodes": {
            "auth-service:8081": 1
          }
        }
      }'

# Define a route for registry-service
curl -X PUT http://apisix:9180/apisix/admin/routes/2 \
  -H "X-API-KEY: $ADMIN_API_KEY" \
  -d '{
        "uri": "/registry/*",
        "methods": ["GET", "POST", "PUT", "DELETE"],
        "upstream": {
          "type": "roundrobin",
          "nodes": {
            "registry-service:8082": 1
          }
        },
        "plugins": {
          "jwt-auth": {}
        }
      }'

# Define a route for publisher-service
curl -X PUT http://apisix:9180/apisix/admin/routes/3 \
  -H "X-API-KEY: $ADMIN_API_KEY" \
  -d '{
        "uri": "/publisher/*",
        "methods": ["GET", "POST", "PUT", "DELETE"],
        "upstream": {
          "type": "roundrobin",
          "nodes": {
            "publisher-service:8083": 1
          }
        },
        "plugins": {
          "jwt-auth": {}
        }
      }'

echo "Routes configured successfully!"
