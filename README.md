# PortableBuildTools
Create standalone instances of VS Build Tools.

> [!Warning]
> #### Under Development
> ###### There is no stable version.

# Requirements
- Go 1.22 or later.
- Rust 1.78 or later.
- GNU Make 4.4.1 or later.

## Building
- Run `make all`. Output is located in `build/`
- Run `portablebuildtools.exe` found under `build/`.
    - Use `portablebuildtools -h` to view a list of commands.

# Usage
- Run `portablebuildtools.exe` after building to download the default packages.
    - ARM/ARM64 targets can be downloaded by appending the `--arm-targets` flag.
    - Spectre libraries can be downloaded by appending the `--spectre-libs` flag.
    - LLVM Clang can be downloaded by appending the `--llvm-clang` flag.
    - UnitTest framework can be downloaded by appending the `--unittest` flag.
    - Cmake can be downloaded by appending the `--cmake` flag.
    - MFC/ATL libraries can be downloaded by appending the `--mfc-atl` flag.

# Registering (Experimental)
- Modify `_Instances/state.json` paths to match the download location. Modify additional files and values as needed.
- Run `vs_register.bat`
- Run `vswhere.exe` to check if your instance is registered.

## Testing
Rudimentary testing is located in `tests/*`
- After downloading the standalone SDK, run the tests via `build_win.bat`.
- View dumpbin header information to ensure the executables were built with the correct architecture.
- ARM and ARM64 targets have not been extensively tested on any substantial projects and may not function correctly.

## Information
- The default download options include (almost) everything you need to build and run the majority of C/C++ applications and libraries.
- Use [session.log](./session.log) to view all files downloaded with the `microsoft.vc.{msvc_ver}, microsoft.visualcpp, microsoft.vs, microsoft.visualstudio` prefixes.
- A modified version of ([MSIExtract](https://github.com/Super-Pizza/MSIExtract)) is used over options such as `msiexec` or `lessmsi` due to issues with how arguments and spaces in paths are handled with the alternatives.

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