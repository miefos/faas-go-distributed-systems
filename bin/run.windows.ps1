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