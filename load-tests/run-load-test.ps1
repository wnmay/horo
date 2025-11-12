# Load Test Runner
Write-Host "Horo Load Test Runner" -ForegroundColor Cyan

$envFile = Join-Path $PSScriptRoot ".env"
if (Test-Path $envFile) {
    Write-Host "Loading .env file..." -ForegroundColor Green
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^([^#][^=]+)=(.*)$') {
            $key = $matches[1].Trim()
            $value = $matches[2].Trim().Trim('"', "'")
            [Environment]::SetEnvironmentVariable($key, $value, "Process")
            Write-Host "  $key = SET" -ForegroundColor Gray
        }
    }
}

Write-Host ""
Write-Host "Config:" -ForegroundColor Cyan
$baseUrl = if ($env:BASE_URL) { $env:BASE_URL } else { "http://localhost:8080" }
Write-Host "  BASE_URL: $baseUrl" -ForegroundColor White
$apiKeyStatus = if ($env:FIREBASE_API_KEY) { "SET" } else { "NOT SET" }
Write-Host "  API_KEY: $apiKeyStatus" -ForegroundColor White
Write-Host ""

if (-not $env:FIREBASE_API_KEY) {
    Write-Host "ERROR: FIREBASE_API_KEY not set!" -ForegroundColor Red
    exit 1
}

Write-Host "Running K6 test..." -ForegroundColor Green
Set-Location $PSScriptRoot
k6 run full-flow-load-test.js
exit $LASTEXITCODE
