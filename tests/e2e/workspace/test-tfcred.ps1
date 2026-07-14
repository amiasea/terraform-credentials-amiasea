$tfcred = "$env:APPDATA\terraform.d\plugins\terraform-credentials-tfcred.exe"

$domain = "app.terraform.io"

Write-Host "Calling tfcred:"
Write-Host "  Binary: $tfcred"
Write-Host "  Domain: $domain"
Write-Host ""

$output = & $tfcred get $domain 2>&1

if ($LASTEXITCODE -ne 0) {
    Write-Error "tfcred failed with exit code $LASTEXITCODE"
    Write-Host $output
    exit $LASTEXITCODE
}

Write-Host "Response:"
$output