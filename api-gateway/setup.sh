#!/bin/bash

# Function to install curl
install_curl() {
  # Check for package manager and install curl accordingly
  if command -v apk > /dev/null 2>&1; then
    echo "Installing curl using apk (Alpine Linux)..."
    apk add --no-cache curl
  elif command -v apt-get > /dev/null 2>&1; then
    echo "Installing curl using apt-get (Debian/Ubuntu)..."
    apt-get update && apt-get install -y curl
  elif command -v yum > /dev/null 2>&1; then
    echo "Installing curl using yum (CentOS/RHEL)..."
    yum install -y curl
  elif command -v dnf > /dev/null 2>&1; then
    echo "Installing curl using dnf (Fedora)..."
    dnf install -y curl
  elif command -v zypper > /dev/null 2>&1; then
    echo "Installing curl using zypper (SUSE)..."
    zypper install -y curl
  else
    echo "Error: Could not detect package manager. Install curl manually."
    exit 1
  fi
}

# Check if curl is already installed
if command -v curl > /dev/null 2>&1; then
  echo "curl is already installed."
else
  echo "curl is not installed. Installing..."
  install_curl
fi

# Verify curl installation
if command -v curl > /dev/null 2>&1; then
  echo "curl is installed."
else
  echo "curl installation failed."
  exit 1
fi

ADMIN_API_KEY="edd1c9f034335f136f87ad84b625c8f1"
JWT_SECRET_KEY=$(<secret.txt)
JWT_KEY="faas-app-key"

# Wait a bit until docker containers start up
sleep 7

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
