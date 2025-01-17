import time
import requests

ADMIN_API_KEY = "edd1c9f034335f136f87ad84b625c8f1"
JWT_SECRET_KEY = open("secret.txt").read().strip()
JWT_KEY = "faas-app-key"

# Wait until apisix container is healthy
while True:
    try:
        response = requests.get("http://apisix:9180/apisix/admin/routes", headers={"X-API-KEY": ADMIN_API_KEY})
        if response.status_code == 200:
            break
    except requests.exceptions.RequestException:
        pass
    print("Waiting for apisix service to be ready...")
    time.sleep(1)

# Configure consumers
consumer_data = {
    "username": "jwt_consumer",
    "plugins": {
        "jwt-auth": {
            "key": JWT_KEY,
            "secret": JWT_SECRET_KEY
        }
    }
}
requests.put("http://apisix:9180/apisix/admin/consumers", headers={"X-API-KEY": ADMIN_API_KEY}, json=consumer_data)

# Define a route for auth-service
auth_service_route = {
    "uri": "/auth/*",
    "methods": ["GET", "POST", "PUT", "DELETE"],
    "upstream": {
        "type": "roundrobin",
        "nodes": {
            "auth-service:8081": 1
        }
    }
}
requests.put("http://apisix:9180/apisix/admin/routes/1", headers={"X-API-KEY": ADMIN_API_KEY}, json=auth_service_route)

# Define a route for registry-service
registry_service_route = {
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
}
requests.put("http://apisix:9180/apisix/admin/routes/2", headers={"X-API-KEY": ADMIN_API_KEY}, json=registry_service_route)

# Define a route for publisher-service
publisher_service_route = {
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
}
requests.put("http://apisix:9180/apisix/admin/routes/3", headers={"X-API-KEY": ADMIN_API_KEY}, json=publisher_service_route)

print("Routes configured successfully!")