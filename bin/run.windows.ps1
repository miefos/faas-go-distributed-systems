# Color definitions
$GREEN = "`e[32m"
$YELLOW = "`e[33m"
$CYAN = "`e[36m"
$RED = "`e[31m"
$RESET = "`e[0m"

### 1. Generate Secret
if (-not (Test-Path -Path "secret.txt")) {
    Write-Host "$CYAN[1] Generating Secret$RESET"

    # Generate a random secret and save it to secret.txt
    $secretBytes = New-Object byte[] 32
    (New-Object Security.Cryptography.RNGCryptoServiceProvider).GetBytes($secretBytes)
    $secret = [Convert]::ToBase64String($secretBytes)
    Set-Content -Path "secret.txt" -Value $secret
    Set-Content -Path ".\auth-service\secret.txt" -Value $secret
    Set-Content -Path ".\api-gateway\secret.txt" -Value $secret

    Write-Host "$GREEN Secret generated and copied to auth-service and api-gateway.$RESET"
    Get-Content -Path "secret.txt"
    Get-Content -Path ".\auth-service\secret.txt"
    Get-Content -Path ".\api-gateway\secret.txt"
    Write-Host ""
} else {
    Write-Host "$YELLOW secret.txt already exists. Skipping secret generation.$RESET"
}

### 2. Build images
Write-Host "$CYAN[2] Building images$RESET"
docker build -t apisix-configurator:latest .\api-gateway
docker build -t auth-service:latest .\auth-service
docker build -t registry-service:latest .\registry-service
docker build -t spawner-service:latest .\spawner-service
docker build -t publisher-service:latest .\publisher-service
Write-Host "$GREEN[2] Done: Images built$RESET"
Write-Host ""

### 3. Init swarm
Write-Host "$CYAN[3] Initializing Swarm$RESET"
$swarmState = docker info --format '{{.Swarm.LocalNodeState}}'
if ($swarmState -ne "active") {
    Write-Host "$YELLOW Swarm is not initialized. Initializing Swarm...$RESET"
    docker swarm init
} else {
    Write-Host "$GREEN Swarm is already initialized.$RESET"
}
Write-Host "$GREEN[3] Done: Swarm initialized$RESET"
Write-Host ""

### 4. Docker compose up
Write-Host "$CYAN[4] Starting services$RESET"
docker compose up -d
Write-Host "$GREEN[4] Done: Services started$RESET"
Write-Host ""

### 5. Setup API gateway
Write-Host "$CYAN[5] Setting up API Gateway$RESET"

# Get container IDs that match the "apisix" filter
$container_ids = docker ps --filter "name=apisix" --format "{{.ID}}"

if (-not $container_ids) {
    Write-Host "$RED[ERROR] No containers matching 'apisix' found.$RESET"
    exit 1
}

# Loop over each container and run setup.sh
foreach ($container_id in $container_ids) {
    Write-Host "$CYAN  Setting up API Gateway in container $container_id$RESET"
    docker exec --user root $container_id powershell -Command "./setup.ps1"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "$RED[ERROR] Failed to run setup.ps1 in container $container_id.$RESET"
        exit 1
    }
    Write-Host "$GREEN  Done: API Gateway setup in container $container_id$RESET"
}
Write-Host "$GREEN[5] Done: API Gateway setup$RESET"
Write-Host ""

# 6. Finish
Write-Host "$GREEN Setup complete.$RESET"