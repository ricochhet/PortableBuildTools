& "$PSScriptRoot/build.ps1"

Set-Location -Path "./build/"
python downloader.py
Set-Location -Path ".."