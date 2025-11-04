# # Load environment variables from .env file
# $envFile = Join-Path $PSScriptRoot ".env"

# if (Test-Path $envFile) {
#     Write-Host "Loading environment variables from .env file..." -ForegroundColor Green
    
#     Get-Content $envFile | ForEach-Object {
#         if ($_ -match '^\s*([^#][^=]+)=(.*)$') {
#             $name = $matches[1].Trim()
#             $value = $matches[2].Trim()
            
#             # Remove quotes if present
#             $value = $value -replace '^"(.*)"$', '$1'
#             $value = $value -replace "^'(.*)'$", '$1'
            
#             [Environment]::SetEnvironmentVariable($name, $value, "Process")
#             Write-Host "  Loaded: $name" -ForegroundColor Cyan
#         }
#     }
    
#     Write-Host ""
# } else {
#     Write-Host "Warning: .env file not found at $envFile" -ForegroundColor Yellow
#     Write-Host "Create a .env file with AUTH_TOKEN=your_token" -ForegroundColor Yellow
#     Write-Host ""
# }

# # Run k6 with the test file passed as argument
# $testFile = $args[0]

# if (-not $testFile) {
#     Write-Host "Usage: .\run-test.ps1 <test-file>" -ForegroundColor Red
#     Write-Host "Example: .\run-test.ps1 full-flow-test.js" -ForegroundColor Yellow
#     Write-Host "  .\run-test.ps1 full-flow-test.js         (k6)" -ForegroundColor Yellow
#     Write-Host "  .\run-test.ps1 chat-rabbit-latency-test.mjs  (Node.js)" -ForegroundColor Yellow
#     exit 1
# }

# Write-Host "Running k6 test: $testFile" -ForegroundColor Green
# Write-Host "==========================================`n" -ForegroundColor Green

# switch ($ext) {
#     ".mjs" {
#         # Node.js (RabbitMQ latency test)
#         Write-Host "Detected .mjs -> running with Node.js" -ForegroundColor Cyan
#         node $testFile
#     }
#     ".js" {
#         # k6 (load test)
#         Write-Host "Detected .js -> running with k6" -ForegroundColor Cyan
#         k6 run $testFile
#     }
#     default {
#         Write-Host "Unknown extension '$ext'. Use .js for k6 or .mjs for Node.js." -ForegroundColor Red
#         exit 1
#     }
# }

# k6 run $testFile

# Load environment variables from .env file
$envFile = Join-Path $PSScriptRoot ".env"

if (Test-Path $envFile) {
    Write-Host "Loading environment variables from .env file..." -ForegroundColor Green
    
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^\s*([^#][^=]+)=(.*)$') {
            $name = $matches[1].Trim()
            $value = $matches[2].Trim()
            
            # Remove quotes if present
            $value = $value -replace '^"(.*)"$', '$1'
            $value = $value -replace "^'(.*)'$", '$1'
            
            [Environment]::SetEnvironmentVariable($name, $value, "Process")
            Write-Host "  Loaded: $name" -ForegroundColor Cyan
        }
    }
    
    Write-Host ""
} else {
    Write-Host "Warning: .env file not found at $envFile" -ForegroundColor Yellow
    Write-Host "Create a .env file with AUTH_TOKEN=your_token" -ForegroundColor Yellow
    Write-Host ""
}

# Run test file passed as argument
$testFile = $args[0]

if (-not $testFile) {
    Write-Host "Usage: .\run-test.ps1 <test-file>" -ForegroundColor Red
    Write-Host "Examples:" -ForegroundColor Yellow
    Write-Host "  .\run-test.ps1 full-flow-test.js" -ForegroundColor Yellow
    Write-Host "  .\run-test.ps1 chat-rabbit-latency-test.mjs" -ForegroundColor Yellow
    exit 1
}

$ext = [System.IO.Path]::GetExtension($testFile).ToLower()

Write-Host "Running test: $testFile" -ForegroundColor Green
Write-Host "=========================================`n" -ForegroundColor Green

switch ($ext) {
    ".mjs" {
        Write-Host "Detected .mjs -> running with Node.js" -ForegroundColor Cyan
        node $testFile
    }
    ".js" {
        Write-Host "Detected .js -> running with k6" -ForegroundColor Cyan
        k6 run $testFile
    }
    default {
        Write-Host "Unknown extension '$ext'. Use .js for k6 or .mjs for Node.js." -ForegroundColor Red
        exit 1
    }
}
