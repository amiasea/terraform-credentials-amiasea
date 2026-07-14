$ErrorActionPreference = "Stop"

Write-Host "[tfcred] Starting uninstallation..." -ForegroundColor Cyan

# ==============================================================================
# 1. Path Configurations
# ==============================================================================
$terraformPluginDir = "$env:APPDATA\terraform.d\plugins"
$targetFileName     = "terraform-credentials-tfcred.exe"
$fullBinaryPath     = Join-Path $terraformPluginDir $targetFileName

# Track both local metadata storage and the program tracking directory
$tfcredLocalDataDir = Join-Path $env:LOCALAPPDATA "tfcred"
$tfcredProgramDir   = Join-Path $env:LOCALAPPDATA "Programs\tfcred"
$configJsonPath     = Join-Path $env:USERPROFILE "terraform.tfrc.json"

# ==============================================================================
# 2. Vault Cleanup (Purges keys explicitly tracked inside contexts.json before drop)
# ==============================================================================
$localContextPath = Join-Path $tfcredLocalDataDir "contexts.json"
if (Test-Path $localContextPath) {
    try {
        $contextStore = Get-Content -Path $localContextPath -Raw | ConvertFrom-Json
        if ($contextStore.Contexts) {
            Write-Host "[tfcred] Purging associated secrets from Windows Credential Manager..." -ForegroundColor Gray
            foreach ($contextKey in $contextStore.Contexts.PSObject.Properties.Name) {
                $entry = $contextStore.Contexts.$contextKey
                $domainFormatted = $entry.Domain -replace '\.', '_'
                $vaultTargetKey = "tfcred:domain:$($domainFormatted):$($entry.TokenType):$($entry.Org)"
                
                if (cmdkey /list | Where-Object { $_ -match [regex]::Escape($vaultTargetKey) }) {
                    cmdkey /delete:$vaultTargetKey | Out-Null
                }
            }
        }
    } catch {
        Write-Host "[tfcred] ⚠️ Failed to purge vault secrets programmatically: context registry unreadable." -ForegroundColor Yellow
    }
}

# ==============================================================================
# 3. Cleanly Remove tfcred Blocks From Terraform Config File
# ==============================================================================
if (Test-Path $configJsonPath) {
    try {
        $configObject = Get-Content -Path $configJsonPath -Raw | ConvertFrom-Json
        if ($configObject.PSObject.Properties.Name -contains "credentials_helper" -and $configObject.credentials_helper.PSObject.Properties.Name -contains "tfcred") {
            Write-Host "[tfcred] Programmatically stripping credentials_helper configuration block..." -ForegroundColor Gray
            
            # Remove the tfcred block specifically
            $configObject.credentials_helper.PSObject.Properties.Remove("tfcred")
            
            # If credentials_helper is now empty, strip the whole property out completely
            if (($configObject.credentials_helper.PSObject.Properties.Name).Count -eq 0) {
                $configObject.PSObject.Properties.Remove("credentials_helper")
            }
            
            # Save the scrubbed config file back out
            $configObject | ConvertTo-Json -Depth 5 | Set-Content -Path $configJsonPath -Force
            Write-Host "[tfcred] ✅ Programmatically cleaned '$configJsonPath'." -ForegroundColor Green
        }
    } catch {
        Write-Host "[tfcred] ⚠️ Warning: Failed to clean $configJsonPath safely. Leaving file intact." -ForegroundColor Yellow
    }
}

# ==============================================================================
# 4. Delete the Plugin Binary
# ==============================================================================
if (Test-Path $fullBinaryPath) {
    Remove-Item -Path $fullBinaryPath -Force
    Write-Host "[tfcred] ✅ Removed binary from native Terraform plugin path." -ForegroundColor Green
} else {
    Write-Host "[tfcred] ℹ️ Binary not found in plugin path. Skipping." -ForegroundColor Yellow
}

# ==============================================================================
# 5. Clear the User-Profile Locked Context Cache Folder Entirely
# ==============================================================================
if (Test-Path $tfcredLocalDataDir) {
    Remove-Item -Path $tfcredLocalDataDir -Recurse -Force
    Write-Host "[tfcred] ✅ Successfully deleted user context cache folder." -ForegroundColor Green
} else {
    Write-Host "[tfcred] ℹ️ Local context directory not found. Skipping." -ForegroundColor Yellow
}

# ==============================================================================
# 6. Remove the Native Windows App Paths Alias
# ==============================================================================
$registryPath = "HKCU:\Software\Microsoft\Windows\CurrentVersion\App Paths\tfcred.exe"
if (Test-Path $registryPath) {
    Remove-Item -Path $registryPath -Recurse -Force
    Write-Host "[tfcred] ✅ Removed global 'tfcred' execution alias mapping." -ForegroundColor Green
} else {
    Write-Host "[tfcred] ℹ️ App Paths registry key not found. Skipping." -ForegroundColor Yellow
}

# ==============================================================================
# 7. Clear the Environment Variable Linkage
# ==============================================================================
[Environment]::SetEnvironmentVariable("TF_CLI_CONFIG_FILE", $null, "User")
$env:TF_CLI_CONFIG_FILE = $null
Write-Host "[tfcred] ✅ Unlinked TF_CLI_CONFIG_FILE environment pointer." -ForegroundColor Green

# ==============================================================================
# 8. Complete Program Folder Purge
# ==============================================================================
# This works cleanly in a native WinGet uninstall flow because WinGet executes
# the script file out of a transient path outside of this target directory.
if (Test-Path $tfcredProgramDir) {
    Remove-Item -Path $tfcredProgramDir -Recurse -Force
    Write-Host "[tfcred] ✅ Successfully purged program installation tracking folder." -ForegroundColor Green
}

Write-Host "[tfcred] 🎉 Uninstallation completed successfully! Machine reverted to standard defaults." -ForegroundColor Cyan
