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

New-Item -ItemType Directory -Path $OutputDir -Force | Out-Null

# version.yaml
@"
# yaml-language-server: `$schema=https://aka.ms/winget-manifest.version.schema.json
PackageIdentifier: $PackageIdentifier
PackageVersion: $Version
DefaultLocale: en-US
ManifestType: version
ManifestVersion: $ManifestVersion
"@ | Out-File "$OutputDir\$PackageIdentifier.yaml" -Encoding UTF8

# installer.yaml
@"
# yaml-language-server: `$schema=https://aka.ms/winget-manifest.installer.schema.json
PackageIdentifier: $PackageIdentifier
PackageVersion: $Version
InstallerLocale: en-US
InstallerType: zip
NestedInstallerType: portable
InstallModes:
  - silent
  - silentWithProgress
Commands:
  - tfcred
NestedInstallerFiles:
  - RelativeFilePath: tfcred.exe
  - RelativeFilePath: terraform-credentials-tfcred.exe
Installers:
  - Architecture: x64
    InstallerUrl: $InstallerUrl
    InstallerSha256: $InstallerSha256
    AppsAndFeaturesEntries:
      - DisplayName: tfcred
        Publisher: amiasea
        DisplayVersion: $Version
        ProductCode: tfcred-v$Version
ManifestType: installer
ManifestVersion: $ManifestVersion
"@ | Out-File "$OutputDir\$PackageIdentifier.installer.yaml" -Encoding UTF8

# locale.en-US.yaml
@"
# yaml-language-server: `$schema=https://aka.ms/winget-manifest.defaultLocale.schema.json
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
"@ | Out-File "$OutputDir\$PackageIdentifier.locale.en-US.yaml" -Encoding UTF8

Write-Host "✅ Winget manifests successfully generated for version $Version" -ForegroundColor Green
Write-Host "   Location: $OutputDir" -ForegroundColor Gray
Write-Host "   Remember: Users should run 'tfcred init' after installation." -ForegroundColor Cyan