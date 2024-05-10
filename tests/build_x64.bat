@echo off
echo Building x64...
call "../build/sdk_standalone/set_vars64.bat"
cl.exe /EHsc /Fe:main_x64.exe main.cpp
dumpbin /headers main_x64.exe