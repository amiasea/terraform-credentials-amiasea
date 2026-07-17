param ()

$ErrorActionPreference = "Stop"

Write-Host "[tfcred] Starting uninstallation..." -ForegroundColor Cyan

# ==============================================================================
# Paths
# ==============================================================================

$terraformPluginDir = Join-Path `
    $env:APPDATA `
    "terraform.d\plugins"

$helperBinary = Join-Path `
    $terraformPluginDir `
    "terraform-credentials-tfcred.exe"

$tfcredProgramDir = Join-Path `
    $env:LOCALAPPDATA `
    "Programs\amiasea\tfcred"

$cliBinary = Join-Path `
    $tfcredProgramDir `
    "tfcred.exe"

$bootstrapBinary = Join-Path `
    $tfcredProgramDir `
    "tfcred-bootstrap.exe"

$uninstallScript = Join-Path `
    $tfcredProgramDir `
    "uninstall.ps1"

$tfcredConfigDir = Join-Path `
    $env:LOCALAPPDATA `
    "amiasea\tfcred"

$configHCLPath = Join-Path `
    $tfcredConfigDir `
    "terraform.tfrc"

$uninstallKey = `
    "HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall\tfcred"

# ==============================================================================
# Remove Terraform credentials helper binary
# ==============================================================================

if (Test-Path $helperBinary) {

    Remove-Item `
        -Path $helperBinary `
        -Force

    Write-Host "[tfcred] Removed Terraform credentials helper binary." -ForegroundColor Green
}
else {

    Write-Host "[tfcred] Credentials helper binary not found. Skipping." -ForegroundColor Yellow
}

# ==============================================================================
# Remove Terraform credentials helper configuration
#
# tfcred owns the credentials_helper section.
# Terraform only permits one credentials helper.
# ==============================================================================

if (Test-Path $configHCLPath) {

    try {

        $fileContent = Get-Content `
            -Path $configHCLPath `
            -Raw

        $pattern = '(?ms)credentials_helper\s+"tfcred"\s*\{.*?\}\s*'

        if ($fileContent -match $pattern) {

            $fileContent = [regex]::Replace(
                $fileContent,
                $pattern,
                ""
            )

            $fileContent = $fileContent.TrimEnd() + "`r`n"

            Set-Content `
                -Path $configHCLPath `
                -Value $fileContent `
                -Encoding UTF8 `
                -Force

            Write-Host "[tfcred] Removed Terraform credentials helper configuration." -ForegroundColor Green
        }

    }
    catch {

        Write-Host "[tfcred] Warning: Unable to update Terraform configuration." -ForegroundColor Yellow
        Write-Host $_.Exception.Message -ForegroundColor Yellow
    }
}

# ==============================================================================
# Remove tfcred from user PATH
# ==============================================================================

$userPath = [Environment]::GetEnvironmentVariable(
    "Path",
    "User"
)

if ($userPath) {

    $newPath = ($userPath.Split(";") |
        Where-Object {
            $_ -and ($_ -ne $tfcredProgramDir)
        }) -join ";"

    [Environment]::SetEnvironmentVariable(
        "Path",
        $newPath,
        "User"
    )

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

    Write-Host "[tfcred] Removed tfcred from user PATH." -ForegroundColor Green
}

# ==============================================================================
# Remove application files
# ==============================================================================

foreach ($file in @(
    $cliBinary,
    $bootstrapBinary,
    $uninstallScript
)) {

    if (Test-Path $file) {

        Remove-Item `
            -Path $file `
            -Force

        Write-Host "[tfcred] Removed $file" -ForegroundColor Green
    }
}

# ==============================================================================
# Remove Apps & Features registration
# ==============================================================================

if (Test-Path $uninstallKey) {

    Remove-Item `
        -Path $uninstallKey `
        -Recurse `
        -Force

    Write-Host "[tfcred] Removed Apps & Features registration." -ForegroundColor Green
}

# ==============================================================================
# Remove installation directory
# ==============================================================================

if (Test-Path $tfcredProgramDir) {

    try {

        Remove-Item `
            -Path $tfcredProgramDir `
            -Recurse `
            -Force

        Write-Host "[tfcred] Removed installation directory." -ForegroundColor Green
    }
    catch {

        Write-Host "[tfcred] Warning: Unable to remove installation directory immediately." -ForegroundColor Yellow
        Write-Host $_.Exception.Message -ForegroundColor Yellow
    }
}

# ==============================================================================
# Intentionally preserved:
#
# - %LOCALAPPDATA%\amiasea\tfcred
# - terraform.tfrc
# - tfcred contexts
# - Windows Credential Manager tokens
# - user state
# - TF_CLI_CONFIG_FILE environment variable
#
# Full data cleanup is user-controlled:
#
#     tfcred purge
#
# ==============================================================================

Write-Host "[tfcred] Uninstallation completed successfully." -ForegroundColor Cyan