<#
.SYNOPSIS
    Generates Winget manifests for tfcred
#>

param(
    [Parameter(Mandatory=$true)]
    [string]$Version,

    [Parameter(Mandatory=$true)]
    [string]$InstallerUrl,

    [Parameter(Mandatory=$true)]
    [string]$InstallerSha256,

    [string]$OutputDir = "manifests/a/amiasea/tfcred/$Version"
)

$PackageIdentifier = "amiasea.tfcred"
$ManifestVersion = "1.12.0"

New-Item `
    -ItemType Directory `
    -Path $OutputDir `
    -Force | Out-Null


# version.yaml
@"
# yaml-language-server: `$schema=https://aka.ms/winget-manifest.version.1.12.0.schema.json
PackageIdentifier: $PackageIdentifier
PackageVersion: $Version
DefaultLocale: en-US
ManifestType: version
ManifestVersion: $ManifestVersion
"@ | Out-File `
    "$OutputDir\$PackageIdentifier.yaml" `
    -Encoding utf8


# installer.yaml
@"
# yaml-language-server: `$schema=https://aka.ms/winget-manifest.installer.1.12.0.schema.json
PackageIdentifier: $PackageIdentifier
PackageVersion: $Version
InstallerLocale: en-US

Installers:
  - Architecture: x64
    InstallerUrl: $InstallerUrl
    InstallerSha256: $InstallerSha256

    InstallerType: zip
    NestedInstallerType: exe

    NestedInstallerFiles:
      - RelativeFilePath: tfcred-bootstrap.exe

    Scope: user

    InstallModes:
      - silent
      - silentWithProgress

    AppsAndFeaturesEntries:
      - DisplayName: terraform-credentials-tfcred
        Publisher: amiasea
        DisplayVersion: $Version

ManifestType: installer
ManifestVersion: $ManifestVersion
"@ | Out-File `
    "$OutputDir\$PackageIdentifier.installer.yaml" `
    -Encoding utf8


# locale.en-US.yaml
@"
# yaml-language-server: `$schema=https://aka.ms/winget-manifest.defaultLocale.1.12.0.schema.json
PackageIdentifier: $PackageIdentifier
PackageVersion: $Version
PackageLocale: en-US

Publisher: amiasea
PublisherUrl: https://github.com/amiasea/terraform-credentials-tfcred

PackageName: tfcred
PackageUrl: https://github.com/amiasea/terraform-credentials-tfcred

License: MIT

ShortDescription: Terraform credential context manager

ManifestType: defaultLocale
ManifestVersion: $ManifestVersion
"@ | Out-File `
    "$OutputDir\$PackageIdentifier.locale.en-US.yaml" `
    -Encoding utf8


Write-Host "Winget manifests generated successfully." -ForegroundColor Green
Write-Host "Version: $Version" -ForegroundColor Gray
Write-Host "Installer: $InstallerUrl" -ForegroundColor Gray
Write-Host "SHA256: $InstallerSha256" -ForegroundColor Gray
Write-Host "Location: $OutputDir" -ForegroundColor Gray