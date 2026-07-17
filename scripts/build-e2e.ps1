param (
    [string]$OutputPath = "dist-dev\e2e"
)

$ErrorActionPreference = "Stop"

$Version = "0.0.0-dev"

Write-Host "[tfcred] Building E2E installer package..." -ForegroundColor Cyan
Write-Host "[tfcred] Version: $Version" -ForegroundColor Gray

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$outputDir = Join-Path $repoRoot $OutputPath

# ==============================================================================
# Clean output directory
# ==============================================================================

if (Test-Path $outputDir) {
    Remove-Item `
        -Path $outputDir `
        -Recurse `
        -Force
}

New-Item `
    -ItemType Directory `
    -Path $outputDir `
    -Force | Out-Null


# ==============================================================================
# Build user-facing CLI
# ==============================================================================

$cliOutput = Join-Path `
    $outputDir `
    "tfcred.exe"

Write-Host "[tfcred] Building CLI..." -ForegroundColor Yellow

go build `
    -ldflags "-X main.version=$Version" `
    -o $cliOutput `
    ".\cmd\tfcred"

if ($LASTEXITCODE -ne 0) {
    throw "Failed to build tfcred.exe"
}


# ==============================================================================
# Build Terraform credentials helper
# ==============================================================================

$helperOutput = Join-Path `
    $outputDir `
    "terraform-credentials-tfcred.exe"

Write-Host "[tfcred] Building Terraform credentials helper..." -ForegroundColor Yellow

go build `
    -ldflags "-X main.version=$Version" `
    -o $helperOutput `
    ".\cmd\terraform-credentials-tfcred"

if ($LASTEXITCODE -ne 0) {
    throw "Failed to build terraform-credentials-tfcred.exe"
}


# ==============================================================================
# Build installer bootstrap
# ==============================================================================

$bootstrapOutput = Join-Path `
    $outputDir `
    "tfcred-bootstrap.exe"

Write-Host "[tfcred] Building bootstrap installer..." -ForegroundColor Yellow

go build `
    -ldflags "-X main.version=$Version" `
    -o $bootstrapOutput `
    ".\cmd\installer"

if ($LASTEXITCODE -ne 0) {
    throw "Failed to build tfcred-bootstrap.exe"
}


# ==============================================================================
# Copy installer scripts
# ==============================================================================

Write-Host "[tfcred] Copying installer scripts..." -ForegroundColor Yellow

Copy-Item `
    -Path (Join-Path $repoRoot "scripts\install.ps1") `
    -Destination (Join-Path $outputDir "install.ps1") `
    -Force

Copy-Item `
    -Path (Join-Path $repoRoot "scripts\uninstall.ps1") `
    -Destination (Join-Path $outputDir "uninstall.ps1") `
    -Force


# ==============================================================================
# Validate package layout
# ==============================================================================

$requiredFiles = @(
    "tfcred.exe",
    "terraform-credentials-tfcred.exe",
    "tfcred-bootstrap.exe",
    "install.ps1",
    "uninstall.ps1"
)

foreach ($file in $requiredFiles) {
    $path = Join-Path $outputDir $file

    if (-not (Test-Path $path)) {
        throw "Missing E2E package file: $path"
    }
}


# ==============================================================================
# Summary
# ==============================================================================

Write-Host ""
Write-Host "[tfcred] E2E installer package ready:" -ForegroundColor Green

Get-ChildItem `
    -Path $outputDir |
    ForEach-Object {
        Write-Host "  $($_.Name)"
    }

Write-Host ""
Write-Host "[tfcred] Package path: $outputDir" -ForegroundColor Cyan