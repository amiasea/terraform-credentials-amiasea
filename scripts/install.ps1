param (
    [string]$Version = "dev"
)

$ErrorActionPreference = "Stop"

Write-Host "[tfcred] Starting installation..." -ForegroundColor Cyan
Write-Host "[tfcred] Version: $Version" -ForegroundColor Gray

# ==============================================================================
# Paths
# ==============================================================================

$installerDir = $PSScriptRoot

$terraformPluginDir = Join-Path `
    $env:APPDATA `
    "terraform.d\plugins"

$tfcredProgramDir = Join-Path `
    $env:LOCALAPPDATA `
    "Programs\amiasea\tfcred"

$tfcredConfigDir = Join-Path `
    $env:LOCALAPPDATA `
    "amiasea\tfcred"

$configHCLPath = Join-Path `
    $tfcredConfigDir `
    "terraform.tfrc"

$tfcredHelperLogPath = Join-Path `
    $tfcredConfigDir `
    "credentials-helper.log"

$cliSource = Join-Path `
    $installerDir `
    "tfcred.exe"

$bootstrapSource = Join-Path `
    $installerDir `
    "tfcred-bootstrap.exe"

$helperSource = Join-Path `
    $installerDir `
    "terraform-credentials-tfcred.exe"

$cliTarget = Join-Path `
    $tfcredProgramDir `
    "tfcred.exe"

$bootstrapTarget = Join-Path `
    $tfcredProgramDir `
    "tfcred-bootstrap.exe"

$helperTarget = Join-Path `
    $terraformPluginDir `
    "terraform-credentials-tfcred.exe"

# ==============================================================================
# Validate package contents
# ==============================================================================

foreach ($file in @(
    $cliSource,
    $bootstrapSource,
    $helperSource
)) {
    if (-not (Test-Path $file)) {
        throw "Missing package file: $file"
    }
}

# ==============================================================================
# Create directories
# ==============================================================================

foreach ($directory in @(
    $terraformPluginDir,
    $tfcredProgramDir,
    $tfcredConfigDir
)) {
    New-Item `
        -ItemType Directory `
        -Path $directory `
        -Force | Out-Null
}

# ==============================================================================
# Install binaries
# ==============================================================================

Copy-Item `
    -Path $cliSource `
    -Destination $cliTarget `
    -Force

Copy-Item `
    -Path $bootstrapSource `
    -Destination $bootstrapTarget `
    -Force

Copy-Item `
    -Path $helperSource `
    -Destination $helperTarget `
    -Force

Write-Host "[tfcred] Installed CLI:" -ForegroundColor Green
Write-Host "        $cliTarget"

Write-Host "[tfcred] Installed bootstrap:" -ForegroundColor Green
Write-Host "        $bootstrapTarget"

Write-Host "[tfcred] Installed credentials helper:" -ForegroundColor Green
Write-Host "        $helperTarget"

# ==============================================================================
# Ensure tfcred is on PATH
# ==============================================================================

$userPath = [Environment]::GetEnvironmentVariable(
    "Path",
    "User"
)

$pathEntries = @()

if ($userPath) {
    $pathEntries = $userPath.Split(";")
}

if ($pathEntries -notcontains $tfcredProgramDir) {

    $newPath = if ([string]::IsNullOrWhiteSpace($userPath)) {
        $tfcredProgramDir
    }
    else {
        "$userPath;$tfcredProgramDir"
    }

    [Environment]::SetEnvironmentVariable(
        "Path",
        $newPath,
        "User"
    )

    Write-Host "[tfcred] Added tfcred to user PATH." -ForegroundColor Green

    Add-Type @"
using System;
using System.Runtime.InteropServices;

public static class NativeMethods
{
    [DllImport("user32.dll", SetLastError = true, CharSet = CharSet.Auto)]
    public static extern IntPtr SendMessageTimeout(
        IntPtr hWnd,
        uint Msg,
        UIntPtr wParam,
        string lParam,
        uint fuFlags,
        uint uTimeout,
        out UIntPtr lpdwResult
    );
}
"@

    $HWND_BROADCAST = [IntPtr]0xffff
    $WM_SETTINGCHANGE = 0x001A
    $SMTO_ABORTIFHUNG = 0x0002

    [UIntPtr]$result = [UIntPtr]::Zero

    [void][NativeMethods]::SendMessageTimeout(
        $HWND_BROADCAST,
        $WM_SETTINGCHANGE,
        [UIntPtr]::Zero,
        "Environment",
        $SMTO_ABORTIFHUNG,
        5000,
        [ref]$result
    )

    Write-Host "[tfcred] Broadcasted environment change." -ForegroundColor Green
}
else {
    Write-Host "[tfcred] User PATH already configured." -ForegroundColor Gray
}

# ==============================================================================
# Configure Terraform CLI
# ==============================================================================

$configHCL = @'
credentials_helper "tfcred" {
  args = []
}
'@

$parentDir = Split-Path `
    -Path $configHCLPath `
    -Parent

if (-not (Test-Path $parentDir)) {
    New-Item `
        -ItemType Directory `
        -Path $parentDir `
        -Force | Out-Null
}

if (Test-Path $configHCLPath) {
    $fileContent = Get-Content `
        -Path $configHCLPath `
        -Raw

    $pattern = '(?ms)credentials_helper\s+"tfcred"\s*\{.*?\}'

    if ($fileContent -match $pattern) {
        $fileContent = [regex]::Replace(
            $fileContent,
            $pattern,
            $configHCL
        )
    }
    else {
        $fileContent = $fileContent.TrimEnd()
        $fileContent += "`r`n`r`n"
        $fileContent += $configHCL
        $fileContent += "`r`n"
    }
}
else {
    $fileContent = $configHCL + "`r`n"
}

# Set-Content `
#     -Path $configHCLPath `
#     -Value $fileContent `
#     -Encoding UTF8 `
#     -Force

[System.IO.File]::WriteAllText(
    $configHCLPath,
    $fileContent,
    [System.Text.UTF8Encoding]::new($false)
)

Write-Host "[tfcred] Configured Terraform credentials helper." -ForegroundColor Green
Write-Host "  $configHCLPath"

# ==============================================================================
# Configure Terraform CLI environment
# ==============================================================================

[Environment]::SetEnvironmentVariable(
    "TF_CLI_CONFIG_FILE",
    $configHCLPath,
    "User"
)

Write-Host "[tfcred] Set TF_CLI_CONFIG_FILE." -ForegroundColor Green

# ==============================================================================
# Configure credentials helper logging
# ==============================================================================

[Environment]::SetEnvironmentVariable(
    "TFCRED_CREDENTIALS_HELPER_LOG",
    $tfcredHelperLogPath,
    "User"
)

Write-Host "[tfcred] Set TFCRED_CREDENTIALS_HELPER_LOG." -ForegroundColor Green
Write-Host "  $tfcredHelperLogPath"

# ==============================================================================
# Install uninstall script
# ==============================================================================

$uninstallSource = Join-Path `
    $installerDir `
    "uninstall.ps1"

if (Test-Path $uninstallSource) {

    Copy-Item `
        -Path $uninstallSource `
        -Destination (Join-Path $tfcredProgramDir "uninstall.ps1") `
        -Force

    Write-Host "[tfcred] Installed uninstall script." -ForegroundColor Green
}

# ==============================================================================
# Register Windows Apps & Features entry
# ==============================================================================

$uninstallKey = `
    "HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall\tfcred"

New-Item `
    -Path $uninstallKey `
    -Force | Out-Null

Set-ItemProperty `
    -Path $uninstallKey `
    -Name "DisplayName" `
    -Value "terraform-credentials-tfcred"

Set-ItemProperty `
    -Path $uninstallKey `
    -Name "Publisher" `
    -Value "amiasea"

Set-ItemProperty `
    -Path $uninstallKey `
    -Name "DisplayVersion" `
    -Value $Version

Set-ItemProperty `
    -Path $uninstallKey `
    -Name "UninstallString" `
    -Value "`"$bootstrapTarget`" uninstall"

Write-Host "[tfcred] Registered Apps & Features entry." -ForegroundColor Green

Write-Host "[tfcred] Installation completed successfully." -ForegroundColor Cyan