param (
    [string]$Mode = "install",
    [string]$BinaryPath
)

$ErrorActionPreference = "Stop"

Write-Host "[tfcred] Starting installation for local user..." -ForegroundColor Cyan

# ==============================================================================
# 1. Target Path Configurations
# ==============================================================================
$terraformPluginDir = "$env:APPDATA\terraform.d\plugins"
$targetFileName     = "terraform-credentials-tfcred.exe"
$fullBinaryPath     = Join-Path $terraformPluginDir $targetFileName

# Isolated application directory for administrative artifacts (Uninstaller, etc.)
$tfcredProgramDir   = "$env:LOCALAPPDATA\Programs\tfcred"
$configJsonPath     = Join-Path $env:USERPROFILE "terraform.tfrc.json"
$localContextPath   = Join-Path $env:LOCALAPPDATA "tfcred\contexts.json"

if (!(Test-Path $terraformPluginDir)) {
    New-Item -ItemType Directory -Path $terraformPluginDir -Force | Out-Null
}

if (!(Test-Path $tfcredProgramDir)) {
    New-Item -ItemType Directory -Path $tfcredProgramDir -Force | Out-Null
}

# ==============================================================================
# 2. Strict Lifecycle Rule Execution (Separating Install vs Upgrade)
# ==============================================================================
if ($Mode -eq "upgrade") {
    Write-Host "[tfcred] ℹ️ Processing tool upgrade. Preserving active user context database." -ForegroundColor Yellow
} else {
    if (Test-Path $localContextPath) {
        Remove-Item -Path $localContextPath -Force
        Write-Host "[tfcred] 🛡️ Fresh installation enforced. Cleared legacy context database." -ForegroundColor Green
    }
}

# ==============================================================================
# 3. Verify and Copy Release or Development Asset
# ==============================================================================
# Choose input source: Local dev argument OR side-by-side production build asset
if ($BinaryPath) {
    $srcBin = $BinaryPath
} else {
    $srcBin = Join-Path $PSScriptRoot "terraform-credentials-tfcred.exe"
}

$srcUninstall = Join-Path $PSScriptRoot "uninstall.ps1"

# Validate and deploy execution binary directly to the native Terraform path
if (-not (Test-Path $srcBin)) {
    Write-Host "[tfcred] ❌ Installation failed: Binary not found at: $srcBin" -ForegroundColor Red
    Exit 1
}

Copy-Item -Path $srcBin -Destination $fullBinaryPath -Force
Write-Host "[tfcred] ✅ Placed binary into native Terraform plugin path." -ForegroundColor Green

# Validate and deploy uninstaller wrapper to isolated local app tracking directory
if (Test-Path $srcUninstall) {
    $destUninstall = Join-Path $tfcredProgramDir "uninstall.ps1"
    Copy-Item -Path $srcUninstall -Destination $destUninstall -Force
    Write-Host "[tfcred] ✅ Placed uninstall script into isolated app storage path." -ForegroundColor Green
} else {
    Write-Host "[tfcred] ⚠️ Warning: 'uninstall.ps1' wrapper script missing from release package." -ForegroundColor Yellow
}

# ==============================================================================
# 4. Configure Native Windows App Paths Alias (Per-User tfcred command)
# ==============================================================================
$registryPath = "HKCU:\Software\Microsoft\Windows\CurrentVersion\App Paths\tfcred.exe"
if (!(Test-Path $registryPath)) {
    New-Item -Path $registryPath -Force | Out-Null
}
Set-Item -Path $registryPath -Value $fullBinaryPath
Set-ItemProperty -Path $registryPath -Name "Path" -Value $terraformPluginDir
Write-Host "[tfcred] ✅ Registered user execution alias: 'tfcred' -> '$targetFileName'" -ForegroundColor Green

# ==============================================================================
# 5. Programmatic Terraform CLI Configuration Upsert
# ==============================================================================
if (Test-Path $configJsonPath) {
    try {
        $configObject = Get-Content -Path $configJsonPath -Raw | ConvertFrom-Json
        Write-Host "[tfcred] ℹ️ Updating configuration at $configJsonPath." -ForegroundColor Yellow
    } catch {
        Write-Host "[tfcred] ❌ Installation failed: Configuration file '$configJsonPath' is corrupted." -ForegroundColor Red
        Exit 1
    }
} else {
    Write-Host "[tfcred] Creating a fresh configuration file at $configJsonPath..." -ForegroundColor Gray
    $configObject = [PSCustomObject]@{}
}

$configObject | Add-Member `
    -MemberType NoteProperty `
    -Name "credentials_helper" `
    -Value @{ "tfcred" = [PSCustomObject]@{} } `
    -Force

$configObject | ConvertTo-Json -Depth 5 | Set-Content -Path $configJsonPath -Force
Write-Host "[tfcred] ✅ Programmatically updated '$configJsonPath'." -ForegroundColor Green

# ==============================================================================
# 6. Environment Variable Linkage
# ==============================================================================
[Environment]::SetEnvironmentVariable("TF_CLI_CONFIG_FILE", $configJsonPath, "User")
$env:TF_CLI_CONFIG_FILE = $configJsonPath
Write-Host "[tfcred] ✅ Configured TF_CLI_CONFIG_FILE environment pointer." -ForegroundColor Green

Write-Host "[tfcred] 🎉 Installation completed successfully! Open a fresh terminal and run 'tfcred' to begin." -ForegroundColor Cyan
