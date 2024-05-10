@echo off
echo Building x86...
call "../build/sdk_standalone/set_vars32.bat"
cl.exe /EHsc /Fe:main_x86.exe main.cpp
dumpbin /headers main_x86.exe