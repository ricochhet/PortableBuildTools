mkdir build
cd go
call build.bat
cd ..\rust
call build.bat
cd ..

copy go\PortableBuildTools.exe build\PortableBuildTools.exe
copy rust\target\release\MSIExtract.exe build\MSIExtract.exe