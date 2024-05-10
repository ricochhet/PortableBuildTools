# SdkStandalone
Create standalone instances of VS Build Tools.

# Requirements
- Rust 1.78 or later.
- Python 3.11 or later.
- PowerShell 7 or later.

## Usage
- Run the `downloader.ps1` script. This will run `build.ps1` and `downloader.py`.
- Output is located under `build/sdk_standalone`

## Options
- Set `DOWNLOAD_ARM_TARGETS` in `downloader.py` to `True` if you want to develop for ARM devices.
    - ARM compatibility is untested, and may be broken. Feel free to create a pull request fixing it if it is broken.
- Set `DOWNLOAD_SPECTRE_LIBS` in `downloader.py` to `True` if you want to use the Spectre mitigated libraries.
    - These libraries are not included in `set_vars` scripts.

## Additional
- The default download options include (almost) everything you need to build and run the majority of C/C++ applications and libraries.
- Use [SdkStandaloneSessionReference.log](./SdkStandaloneSessionReference.log) to view all files downloaded with the `microsoft.vc.{msvc_ver}, microsoft.visualcpp, microsoft.vs, microsoft.visualstudio` prefixes.
    - Set `DOWNLOAD_ALL` to `True` if you want to download everything. ~20GB
    - This is useful for instances where you wish to download only a small subset of what comes with VS.
- rust-msiexec ([MSIExtract](https://github.com/Super-Pizza/MSIExtract)) is used over options such as `msiexec` or `lessmsi` due to issues with how arguments and spaces in paths are handled with the alternatives.

### Tools
- `dumpbin /headers example.dll`
- `corflags example.dll`

# License
See [LICENSE](./LICENSE) file.

See [LICENSE-MSIExtract](./LICENSE-MSIExtract) file.