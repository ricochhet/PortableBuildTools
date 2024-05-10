@echo off
echo Building arm64...
call "../build/sdk_standalone/set_vars_arm64.bat"
cl.exe /EHsc /Fe:main_arm64.exe main.cpp
dumpbin /headers main_arm64.exe