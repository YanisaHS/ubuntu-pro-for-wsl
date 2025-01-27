<#
.Synopsis
    Build the WSL Pro Service debian package for local use.
#>

param(
    [Parameter(Mandatory = $False, HelpMessage = "The directory where the debian build artifacts will be stored in")]
    [string]$OutputDir
)

$projectRoot = "${PSScriptRoot}\..\.."

# By default, we store artifacts in the same location dpkg-buildpackage would store them in
if ( "${OutputDir}" -eq "" ) {
    $OutputDir = "${projectRoot}"
}

# Ensure Ubuntu-Preview is installed and registered
$appx = Get-AppxPackage -Name "CanonicalGroupLimited.UbuntuPreview"
if ( $appx -eq "" ) {
    Write-Error "Ubuntu Preview is not installed"
}

$env:WSL_UTF8 = 1

if ( "$(wsl --list --verbose | Select-String Ubuntu-Preview)" -eq "" ) {
    ubuntupreview.exe install --root --ui=none
    if ( "${LastExitCode}" -ne "0" ) {
        Write-Error "could not install Ubuntu-Preview"
        exit 1
    }

    Copy-Item -Path "${HOME}/.gitconfig" -Destination "\\wsl$\Ubuntu-Preview\etc\gitconfig"
}

# Write script to run
$scriptWindows = New-TemporaryFile

$scriptLinux = ( wsl.exe -d Ubuntu-Preview -- wslpath -ua `'${scriptWindows}`' )
if ( "${LastExitCode}" -ne "0" ) {
    Write-Error "could not get build script's linux path"
    exit 1
}

# Using WriteAllText to avoid CRLF
[IO.File]::WriteAllText($scriptWindows, @'
#!/bin/bash
set -eu

git config --get-all safe.directory                        \
    | grep "$(pwd)"                                        \
    || git config --global --add safe.directory "$(pwd)"

# Update internal dependencies in the repo because
# we need git to work, and .git is not rsync'd.
cd wsl-pro-service
go version || (
    sudo DEBIAN_FRONTEND=noninteractive apt update
    sudo DEBIAN_FRONTEND=noninteractive apt install -y golang-go
)
./debian/update-internal-dependencies
go mod tidy
go mod vendor
cd -

# Set up directory
build_dir="${HOME}/wsl-pro-service-build"

modules=("agentapi" "common" "contractsapi" "storeapi" "tools" "wsl-pro-service" "wslserviceapi")

for mod in ${modules[@]}; do
    rsync --recursive --quiet "${mod}" "${build_dir}"
done

# Build
bash -e "${build_dir}/tools/build/build-deb.sh"

# Export artifacts
cp -f ${build_dir}/wsl-pro-service_* "${OutputDir}"

'@)

# Set up output directory
New-Item -Force -ItemType "Directory" -Path "${OutputDir}" | Out-Null

$outputLinux = ( wsl.exe -d Ubuntu-Preview -- wslpath -ua `'${OutputDir}`' )
if ( "${LastExitCode}" -ne "0" ) {
    Write-Error "could not get output dir's linux path"
    exit 1
}

wsl.exe -d Ubuntu-Preview -u root -- bash -ec "chmod +x ${scriptLinux} 2>&1"
if ( "${LastExitCode}" -ne "0" ) {
    Write-Error "could not make build script executable"
    exit 1
}

wsl.exe -d Ubuntu-Preview -u root --cd "${projectRoot}" -- bash -ec "OutputDir=${outputLinux} ${scriptLinux} 2>&1"
if ( "${LastExitCode}" -ne "0" ) {
    Write-Error "could not build deb"
    exit 1
}

Write-Output "Artifacts stored in ${OutputDir}"

Remove-Item -Path "${scriptWindows}"

exit 0
