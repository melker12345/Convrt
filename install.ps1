# Requires -RunAsAdministrator

$ErrorActionPreference = "Stop"

# Define the installation directory
$installDir = "$env:LOCALAPPDATA\Programs\convrt"

# Create the installation directory if it doesn't exist
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
    Write-Host "Created installation directory: $installDir"
}

# Copy the executable to the installation directory
Copy-Item "convrt.exe" -Destination $installDir -Force
Write-Host "Copied convrt.exe to installation directory"

# Add to PATH if not already present
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable(
        "Path",
        "$userPath;$installDir",
        "User"
    )
    Write-Host "Added convrt to PATH"
    Write-Host "Please restart your terminal for the PATH changes to take effect"
} else {
    Write-Host "convrt is already in PATH"
}

Write-Host "`nInstallation complete! You can now use 'convrt' from anywhere." -ForegroundColor Green
Write-Host "Try running 'convrt --help' to get started." -ForegroundColor Green
