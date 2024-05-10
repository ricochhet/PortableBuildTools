@echo off
echo Building arm...
call "../build/sdk_standalone/set_vars_arm32.bat"
cl.exe /EHsc /Fe:main_arm32.exe main.cpp
dumpbin /headers main_arm32.exe