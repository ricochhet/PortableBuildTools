mkdir build
cd go
call build.bat
cd ..\rust
call build.bat
cd ..

copy go\downloader.exe build\downloader.exe
copy rust\target\release\rust-msiexec.exe build\rust-msiexec.exe