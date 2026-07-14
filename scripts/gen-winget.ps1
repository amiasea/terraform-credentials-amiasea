param(
    [Parameter(Mandatory=$true)]
    [string]$Version
)

$ErrorActionPreference = "Stop"

$owner = "amiasea"
$repo = "terraform-credentials-tfcred"
$packageId = "amiasea.tfcred"

# CORRECTED: Standardized lowercase matching to mirror your exact .goreleaser.yml outputs
$archiveName = "terraform-credentials-tfcred_windows_amd64.zip"
$distPath = ".\dist\$archiveName"

if (-not (Test-Path $distPath)) {
    throw "GoReleaser archive package not found at $distPath. Please execute your local goreleaser build cycle first."
}

# Fully qualified absolute URL routing WinGet natively to the case-sensitive GitHub release asset
$url = "https://github.com/$owner/$repo/releases/download/v$Version/terraform-credentials-tfcred_windows_amd64.zip"
$hash = (Get-FileHash $distPath -Algorithm SHA256).Hash.ToLower()

$manifestDir = ".\.winget\manifests\a\$($owner)\tfcred\$($Version)"

if (Test-Path $manifestDir) {
    Remove-Item -Recurse -Force $manifestDir
}

New-Item -ItemType Directory -Force -Path $manifestDir | Out-Null

$versionManifest = @"
PackageIdentifier: $($packageId)
PackageVersion: $($Version)
ManifestType: version
ManifestVersion: 1.12.0
"@

$installerManifest = @"
PackageIdentifier: $($packageId)
PackageVersion: $($Version)
InstallerLocale: en-US
Architecture: x64
InstallerType: zip
NestedInstallerType: custom
Scope: user
InstallModes:
  - silent
  - silentWithProgress
InstallerSwitches:
  Silent: -NoProfile -ExecutionPolicy Bypass -File .\install.ps1 -Mode install
  SilentWithProgress: -NoProfile -ExecutionPolicy Bypass -File .\install.ps1 -Mode install
  Custom: -NoProfile -ExecutionPolicy Bypass -File .\install.ps1 -Mode install
  Upgrade: -NoProfile -ExecutionPolicy Bypass -File .\install.ps1 -Mode upgrade
UninstallCommand: powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%USERPROFILE%\AppData\Local\Programs\tfcred\uninstall.ps1"
Commands:
  - tfcred
  - terraform-credentials-tfcred
NestedInstallerFiles:
  - RelativeFilePath: install.ps1
Installers:
  - Architecture: x64
    InstallerUrl: $($url)
    InstallerSha256: $($hash)
AppsAndFeaturesEntries:
  - DisplayName: tfcred
    Publisher: $($owner)
    DisplayVersion: $($Version)
    ProductCode: tfcred-v$($Version)
ManifestType: installer
ManifestVersion: 1.6.0
"@

$localeManifest = @"
PackageIdentifier: $($packageId)
PackageVersion: $($Version)
PackageLocale: en-US
Publisher: $($owner)
PackageName: tfcred
License: MIT
ShortDescription: Terraform credential context manager
ManifestType: defaultLocale
ManifestVersion: 1.6.0
"@

$versionManifest | Out-File "$manifestDir\$($packageId).yaml" -Encoding UTF8 -NoNewline
$installerManifest | Out-File "$manifestDir\$($packageId).installer.yaml" -Encoding UTF8 -NoNewline
$localeManifest | Out-File "$manifestDir\$($packageId).locale.en-US.yaml" -Encoding UTF8 -NoNewline

Write-Host "[winget] Local multi-manifest configuration blocks created successfully at: $manifestDir" -ForegroundColor Green
