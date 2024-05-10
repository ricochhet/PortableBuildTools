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
- `--show-versions`
    - Action: Store True
    - Description: Show available MSVC and Windows SDK versions.

- `--accept-license`
    Action: Store True
    - Description: Automatically accept the license.

- `--msvc-version`
    - Argument: <version>
    - Description: Get a specific MSVC version.

- `--sdk-version`
    - Argument: <version>
    - Description: Get a specific Windows SDK version.

- `--preview`
    - Action: Store True
    - Description: Use the preview channel for Preview versions.

- `--dl-spectre`
    - Action: Store True
    - Description: Download spectre mitigated libraries.

- `--dl-arm-targets`
    - Action: Store True
    - Description: Download ARM/ARM64 targets.

## Testing
Rudimentary testing is located in `tests/*`
- After downloading the standalone SDK, run the tests via `build_win.bat`.
- View dumpbin header information to ensure the executables were built with the correct architecture.
- ARM and ARM64 targets have not been extensively tested on any substantial projects and may not function correctly.

## Information
- The default download options include (almost) everything you need to build and run the majority of C/C++ applications and libraries.
- Use [SdkStandaloneSessionReference.log](./SdkStandaloneSessionReference.log) to view all files downloaded with the `microsoft.vc.{msvc_ver}, microsoft.visualcpp, microsoft.vs, microsoft.visualstudio` prefixes.
    - Set `DOWNLOAD_ALL` to `True` if you want to download everything with the stated prefixes. ~20GB
- rust-msiexec ([MSIExtract](https://github.com/Super-Pizza/MSIExtract)) is used over options such as `msiexec` or `lessmsi` due to issues with how arguments and spaces in paths are handled with the alternatives.

### Alternatives
- `winlibs-mingw` [brechtsanders/winlibs_mingw](https://github.com/brechtsanders/winlibs_mingw/releases)
- `llvm-mingw` [mstorsjo/llvm-mingw](https://github.com/mstorsjo/llvm-mingw)
- `llvm-clang` [llvm-clang](https://releases.llvm.org/download.html)

### Tools
- `dumpbin /headers example.dll`
- `corflags example.dll`

# License
See [LICENSE](./LICENSE) file.

See [LICENSE-MSIExtract](./LICENSE-MSIExtract) file.