# Requires -RunAsAdministrator

$ErrorActionPreference = "Stop"

# Define the installation directory
$installDir = "$env:LOCALAPPDATA\Programs\convrt"

# Remove from PATH
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -like "*$installDir*") {
    $newPath = ($userPath -split ';' | Where-Object { $_ -ne $installDir }) -join ';'
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "Removed convrt from PATH"
}

# Remove installation directory
if (Test-Path $installDir) {
    Remove-Item -Path $installDir -Recurse -Force
    Write-Host "Removed installation directory"
}

Write-Host "`nUninstallation complete!" -ForegroundColor Green
