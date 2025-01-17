#!/bin/bash

# Color definitions
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
RED='\033[0;31m'
RESET='\033[0m'

### 1. Generate Secret
if [ ! -f secret.txt ]; then
    echo -e "${CYAN}[1] Generating Secret${RESET}"

    touch secret.txt
    chmod 600 secret.txt
    openssl rand -base64 32 | tr -d '\n\r\t ' > secret.txt

    cp -pv secret.txt ./auth-service/secret.txt
    cp -pv secret.txt ./api-gateway/secret.txt

    echo -e "${GREEN}Secret generated and copied to auth-service and api-gateway.${RESET}"
    cat secret.txt; echo
    cat ./auth-service/secret.txt; echo
    cat ./api-gateway/secret.txt; echo
    echo
else
    echo -e "${YELLOW}secret.txt already exists. Skipping secret generation.${RESET}"
fi

### 2. Build images
echo -e "${CYAN}[2] Building images${RESET}"
docker build -t auth-service:latest ./auth-service
docker build -t registry-service:latest ./registry-service
docker build -t spawner-service:latest ./spawner-service
docker build -t publisher-service:latest ./publisher-service
echo -e "${GREEN}[2] Done: Images built${RESET}"
echo

### 3. Init swarm
echo -e "${CYAN}[3] Initializing Swarm${RESET}"
if [[ "$(docker info --format '{{.Swarm.LocalNodeState}}')" != "active" ]]; then
  echo -e "${YELLOW}Swarm is not initialized. Initializing Swarm...${RESET}"
  docker swarm init
else
  echo -e "${GREEN}Swarm is already initialized.${RESET}"
fi
echo -e "${GREEN}[3] Done: Swarm initialized${RESET}"
echo

### 4. Docker compose up
echo -e "${CYAN}[4] Starting services${RESET}"
docker compose up -d --scale auth-service=2 --scale registry-service=2 --scale spawner-service=2 --scale publisher-service=2
echo -e "${GREEN}[4] Done: Services started${RESET}"
echo

### 5. Setup API gateway
echo -e "${CYAN}[5] Setting up API Gateway${RESET}"
### Get container IDs that match the "apisix" filter
container_ids=$(docker ps -q --filter "name=apisix")

if [ -z "$container_ids" ]; then
  echo -e "${RED}[ERROR] No containers matching 'apisix' found.${RESET}"
  exit 1
fi

# Loop over each container and run setup.sh
for container_id in $container_ids; do
  echo -e "${CYAN}  Setting up API Gateway in container $container_id${RESET}"
  docker exec --user root "$container_id" /bin/sh -c "./setup.sh"
  if [ $? -ne 0 ]; then
    echo -e "${RED}[ERROR] Failed to run setup.sh in container $container_id.${RESET}"
    exit 1
  fi
  echo -e "${GREEN}  Done: API Gateway setup in container $container_id${RESET}"
done
echo -e "${GREEN}[5] Done: API Gateway setup${RESET}"
echo

# 6. Finish
echo -e "${GREEN}Setup complete.${RESET}"