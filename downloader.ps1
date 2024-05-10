param(
    [Parameter(Mandatory=$false)]
    [string[]]$PythonScriptArguments
)

if (-not $PythonScriptArguments) {
    $PythonScriptArguments = @()
}

& "$PSScriptRoot/build.ps1"

Set-Location -Path "./build/"
python downloader.py $PythonScriptArguments
Set-Location -Path ".."