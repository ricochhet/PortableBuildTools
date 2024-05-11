# SdkStandalone
Create standalone instances of VS Build Tools.

> [!Warning]
> #### Under Development
> ###### Do not expect anything to work as expected. Feel to contribute.

# Requirements
- Go 1.22 or later.
- Rust 1.78 or later.

## Usage
- Run `build.bat` under `project/`.
- Run `downloader.exe` found under `project/build/`.

## Testing
Rudimentary testing is located in `tests/*`
- After downloading the standalone SDK, run the tests via `build_win.bat`.
- View dumpbin header information to ensure the executables were built with the correct architecture.
- ARM and ARM64 targets have not been extensively tested on any substantial projects and may not function correctly.

## Information
- The default download options include (almost) everything you need to build and run the majority of C/C++ applications and libraries.
- Use [SdkStandaloneSessionReference.log](./SdkStandaloneSessionReference.log) to view all files downloaded with the `microsoft.vc.{msvc_ver}, microsoft.visualcpp, microsoft.vs, microsoft.visualstudio` prefixes.
    - Set `DOWNLOAD_ALL` to `True` if you want to download everything with the stated prefixes. ~20GB
    - Uses `python/downloader.py` script.
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