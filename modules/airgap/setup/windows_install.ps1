<#
.SYNOPSIS
  A install script to Setup and Install standalone RKE2 in Windows for airgap to be used as Worker Nodes.
  This script enables features, sets up environment variables and adds default configuration that are needed to install RKE2 in Windows and join a cluster.`
.DESCRIPTION
  Run the script to setup and install all RKE2 related needs and to join a cluster.
.Parameter Token
    Token of Primary server.`

.EXAMPLE
  Usage:
    Invoke-WebRequest ((New-Object System.Net.WebClient).DownloadString('https://github.com/rancher/rke2/blob/master/windows/rke2-quickstart.ps1'))
    ./rke2-quickstart.ps1 $ServerIP <server-IP> $Token <server-token> $Mode <install-mode> $Version <rke2-version>
#>

[CmdletBinding()]
param (
    [Parameter(Mandatory=$true)]
    [String]
    $Token
)

function Write-InfoLog() {
    Write-Output "[INFO] $($args -join " ")"
}

function Write-WarnLog() {
    Write-Output "[WARN] $($args -join " ")"
}

function Update-Registry-File() {
    (Get-Content -Path c:/etc/rancher/rke2/config.yaml) |
    ForEach-Object {$_ -Replace '$TOKEN', $Token} |
        Set-Content -Path c:/etc/rancher/rke2/config.yaml
    Get-Content -Path c:/etc/rancher/rke2/config.yaml
}

function Update-Cert() {
    Import-Certificate -FilePath "C:\Users\Administrator\ca.pem" -CertStoreLocation cert:\LocalMachine\Root
}

function Setup-Environment-Variables() {
    Write-InfoLog "Setting up environment vars..."
    [System.Environment]::SetEnvironmentVariable(
        "Path",[System.Environment]::GetEnvironmentVariable(
            "Path", [System.EnvironmentVariableTarget]::Machine) + ";c:\var\lib\rancher\rke2\bin;c:\usr\local\bin",
    [System.EnvironmentVariableTarget]::Machine)
}

function Start-rke2() {
    Invoke-Expression -Command "C:\usr\local\bin\rke2.exe agent service --add"
    Get-Service -Name rke2
    Start-Service -Name rke2
}

Update-Registry-File
Update-Cert
Setup-Environment-Variables
Start-rke2