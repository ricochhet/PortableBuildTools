/*
 * PortableBuildTools
 * Copyright (C) 2024 PortableBuildTools contributors
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package flag

import (
	"strings"
)

//nolint:lll // constructing batch scripts.
func NewMsvcX64Environ(msvcVersion, sdkVersion, targetA, targetB, host string, flags *Flags) string {
	base := []string{
		`@echo off`,
		"",
		`SET "ROOT=%~dp0"`,
		`SET "VSINSTALLDIR=%ROOT%\"`,
		`SET "VCINSTALLDIR=%VSINSTALLDIR%VC\"`,
		`SET "VS140COMNTOOLS=%VSINSTALLDIR%Common7\Tools\"`,
		`SET "UCRTVersion=` + sdkVersion + `"`,
		`SET "WindowsSdkDir=%VSINSTALLDIR%Windows Kits\10\"`,
		`SET "UniversalCRTSdkDir=%WindowsSdkDir%"`,
		`SET "WindowsSDKVersion=` + sdkVersion + `\"`,
		`SET "WindowsSDKLibVersion=` + sdkVersion + `\"`,
		`SET "WindowsSDK_ExecutablePath_x64=%VSINSTALLDIR%Windows Kits\10\BIN\%WindowsSDKVersion%` + targetA + `\"`,
		"",
		`SET "LIB="`,
		`SET "INCLUDE="`,
		`SET "LIBPATH="`,
		"",
		`@if exist "%VSINSTALLDIR%Common7\Tools" set "PATH=%VSINSTALLDIR%Common7\Tools;%PATH%"`,
		`@if exist "%VSINSTALLDIR%Common7\IDE" set "PATH=%VSINSTALLDIR%Common7\IDE;%PATH%"`,
		`@if exist "%VCINSTALLDIR%VCPackages" set "PATH=%VCINSTALLDIR%VCPackages;%PATH%"`,
		`@if exist "%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\BIN\Host` + host + `\` + targetA + `" set "PATH=%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\BIN\Host` + host + `\` + targetA + `;%PATH%"`,
		"",
		`@if not "%UCRTVersion%" == "" @set "INCLUDE=%UniversalCRTSdkDir%include\%UCRTVersion%\ucrt;%INCLUDE%"`,
		`@if not "%UCRTVersion%" == "" @set "LIB=%UniversalCRTSdkDir%lib\%UCRTVersion%\ucrt\` + targetA + `;%LIB%"`,
		"",
		`@if not "%WindowsSdkDir%" == "" @set "PATH=%WindowsSdkDir%BIN\` + sdkVersion + `\` + targetA + `;%WindowsSdkDir%BIN\` + sdkVersion + `\` + targetB + `;%PATH%"`,
		`@if not "%WindowsSdkDir%" == "" @set "INCLUDE=%WindowsSdkDir%include\%WindowsSDKVersion%shared;%WindowsSdkDir%include\%WindowsSDKVersion%um;%WindowsSdkDir%include\%WindowsSDKVersion%winrt;%INCLUDE%"`,
		`@if not "%WindowsSdkDir%" == "" @set "LIB=%WindowsSdkDir%lib\%WindowsSDKLibVersion%um\` + targetA + `;%LIB%"`,
		`@if not "%WindowsSdkDir%" == "" @set "LIBPATH=%WindowsLibPath%;%ExtensionSDKDir%\Microsoft.VCLibs\14.0\References\CommonConfiguration\neutral;%LIBPATH%"`,
		"",
		`@if not "%WindowsSDK_ExecutablePath_x64%" == "" @set "PATH=%WindowsSDK_ExecutablePath_x64%;%PATH%"`,
		"",
		`@if exist "%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\LIB\` + targetA + `\store" set "LIB=%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\LIB\` + targetA + `\store;%LIB%"`,
		"",
		`@if exist "%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\ATLMFC\INCLUDE" set "INCLUDE=%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\ATLMFC\INCLUDE;%INCLUDE%"`,
		`@if exist "%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\INCLUDE" set "INCLUDE=%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\INCLUDE;%INCLUDE%"`,
		"",
		`@if exist "%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\ATLMFC\LIB\` + targetA + `" set "LIB=%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\ATLMFC\LIB\` + targetA + `;%LIB%"`,
		`@if exist "%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\LIB\` + targetA + `" set "LIB=%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\LIB\` + targetA + `;%LIB%"`,
		"",
		`@if exist "%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\ATLMFC\LIB\` + targetA + `" set "LIBPATH=%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\ATLMFC\LIB\` + targetA + `;%LIBPATH%"`,
		`@if exist "%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\LIB\` + targetA + `" set "LIBPATH=%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\LIB\` + targetA + `;%LIBPATH%"`,
		"",
		`@if exist "%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\LIB\` + targetA + `\store" set "LIBPATH=%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\LIB\` + targetA + `\store;%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\LIB\` + targetA + `\store\references;%LIBPATH%"`,
	}

	if flags.SpectreLibs {
		base = append(base, "",
			`@if exist "%VCINSTALLDIR%\Tools\MSVC\`+msvcVersion+`\ATLMFC\LIB\SPECTRE\`+targetA+`" set "LIB=%VCINSTALLDIR%\Tools\MSVC\`+msvcVersion+`\ATLMFC\LIB\SPECTRE\`+targetA+`;%LIB%"`,
			`@if exist "%VCINSTALLDIR%\Tools\MSVC\`+msvcVersion+`\LIB\SPECTRE\`+targetA+`" set "LIB=%VCINSTALLDIR%\Tools\MSVC\`+msvcVersion+`\LIB\SPECTRE\`+targetA+`;%LIB%"`,
		)

		base = append(base, "",
			`@if exist "%VCINSTALLDIR%\Tools\MSVC\`+msvcVersion+`\ATLMFC\LIB\SPECTRE\`+targetA+`" set "LIBPATH=%VCINSTALLDIR%\Tools\MSVC\`+msvcVersion+`\ATLMFC\LIB\SPECTRE\`+targetA+`;%LIBPATH%"`,
			`@if exist "%VCINSTALLDIR%\Tools\MSVC\`+msvcVersion+`\LIB\SPECTRE\`+targetA+`" set "LIBPATH=%VCINSTALLDIR%\Tools\MSVC\`+msvcVersion+`\LIB\SPECTRE\`+targetA+`;%LIBPATH%"`,
		)
	}

	if flags.LlvmClang {
		base = append(base, "",
			`@if exist "%VCINSTALLDIR%\Tools\LLVM\`+targetA+`" set "LLVM_PATH=%VCINSTALLDIR%\Tools\LLVM\`+targetA+`"`,
			`@if exist "%LLVM_PATH%\BIN" set "PATH=%LLVM_PATH%\BIN;%PATH%"`,
			`@if exist "%LLVM_PATH%\LIB\CLANG\17\LIB" set "LIB=%LLVM_PATH%\LIB\CLANG\17\LIB;%LIB%"`,
			`@if exist "%LLVM_PATH%\LIB\CLANG\17\LIB" set "LIBPATH=%LLVM_PATH%\LIB\CLANG\17\LIB;%LIBPATH%"`,
			`@if exist "%LLVM_PATH%\LIB\CLANG\17\INCLUDE" set "INCLUDE=%LLVM_PATH%\LIB\CLANG\17\INCLUDE;%INCLUDE%"`,
		)
	}

	return strings.Join(base, "\n")
}

//nolint:lll // constructing batch scripts.
func NewMsvcX86Environ(msvcVersion, sdkVersion, targetA, targetB, host string, flags *Flags) string {
	base := []string{
		`@echo off`,
		"",
		`SET "ROOT=%~dp0"`,
		`SET "VSINSTALLDIR=%ROOT%\"`,
		`SET "VCINSTALLDIR=%VSINSTALLDIR%VC\"`,
		`SET "VS140COMNTOOLS=%VSINSTALLDIR%Common7\Tools\"`,
		`SET "UCRTVersion=` + sdkVersion + `"`,
		`SET "WindowsSdkDir=%VSINSTALLDIR%Windows Kits\10\"`,
		`SET "UniversalCRTSdkDir=%WindowsSdkDir%"`,
		`SET "WindowsSDKVersion=` + sdkVersion + `\"`,
		`SET "WindowsSDKLibVersion=` + sdkVersion + `\"`,
		`SET "WindowsSDK_ExecutablePath_x64=%VSINSTALLDIR%Windows Kits\10\BIN\%WindowsSDKVersion%` + targetA + `\"`,
		"",
		`SET "LIB="`,
		`SET "INCLUDE="`,
		`SET "LIBPATH="`,
		"",
		`@if exist "%VSINSTALLDIR%Common7\Tools" set "PATH=%VSINSTALLDIR%Common7\Tools;%PATH%"`,
		`@if exist "%VSINSTALLDIR%Common7\IDE" set "PATH=%VSINSTALLDIR%Common7\IDE;%PATH%"`,
		`@if exist "%VCINSTALLDIR%VCPackages" set "PATH=%VCINSTALLDIR%VCPackages;%PATH%"`,
		`@if exist "%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\BIN\Host` + host + `\` + targetB + `" set "PATH=%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\BIN\Host` + host + `\` + targetB + `;%PATH%"`,
		`@if exist "%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\BIN\Host` + host + `\` + targetA + `" set "PATH=%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\BIN\Host` + host + `\` + targetA + `;%PATH%"`,
		"",
		`@if not "%UCRTVersion%" == "" @set "INCLUDE=%UniversalCRTSdkDir%include\%UCRTVersion%\ucrt;%INCLUDE%"`,
		`@if not "%UCRTVersion%" == "" @set "LIB=%UniversalCRTSdkDir%lib\%UCRTVersion%\ucrt\` + targetA + `;%LIB%"`,
		"",
		`@if not "%WindowsSdkDir%" == "" @set "PATH=%WindowsSdkDir%BIN\` + sdkVersion + `\` + targetA + `;%WindowsSdkDir%BIN\` + sdkVersion + `\` + targetA + `;%PATH%"`,
		`@if not "%WindowsSdkDir%" == "" @set "INCLUDE=%WindowsSdkDir%include\%WindowsSDKVersion%shared;%WindowsSdkDir%include\%WindowsSDKVersion%um;%WindowsSdkDir%include\%WindowsSDKVersion%winrt;%INCLUDE%"`,
		`@if not "%WindowsSdkDir%" == "" @set "LIB=%WindowsSdkDir%lib\%WindowsSDKLibVersion%um\` + targetA + `;%LIB%"`,
		`@if not "%WindowsSdkDir%" == "" @set "LIBPATH=%WindowsLibPath%;%ExtensionSDKDir%\Microsoft.VCLibs\14.0\References\CommonConfiguration\neutral;%LIBPATH%"`,
		"",
		`@if not "%WindowsSDK_ExecutablePath_x64%" == "" @set "PATH=%WindowsSDK_ExecutablePath_x64%;%PATH%"`,
		"",
		`@if exist "%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\ATLMFC\INCLUDE" set "INCLUDE=%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\ATLMFC\INCLUDE;%INCLUDE%"`,
		`@if exist "%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\INCLUDE" set "INCLUDE=%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\INCLUDE;%INCLUDE%"`,
		"",
		`@if exist "%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\ATLMFC\LIB\` + targetA + `" set "LIB=%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\ATLMFC\LIB\` + targetA + `;%LIB%"`,
		`@if exist "%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\LIB\` + targetA + `" set "LIB=%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\LIB\` + targetA + `;%LIB%"`,
		"",
		`@if exist "%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\LIB\` + targetA + `\store" set "LIB=%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\LIB\` + targetA + `\store;%LIB%"`,
		"",
		`@if exist "%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\ATLMFC\LIB\` + targetA + `" set "LIBPATH=%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\ATLMFC\LIB\` + targetA + `;%LIBPATH%"`,
		`@if exist "%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\LIB\` + targetA + `" set "LIBPATH=%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\LIB\` + targetA + `;%LIBPATH%"`,
		"",
		`@if exist "%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\LIB\` + targetA + `\store" set "LIBPATH=%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\LIB\` + targetA + `\store;%VCINSTALLDIR%\Tools\MSVC\` + msvcVersion + `\LIB\` + targetA + `\store\references;%LIBPATH%"`,
	}

	if flags.SpectreLibs {
		base = append(base, "",
			`  @if exist "%VCINSTALLDIR%\Tools\MSVC\`+msvcVersion+`\ATLMFC\LIB\SPECTRE\`+targetA+`" set "LIB=%VCINSTALLDIR%\Tools\MSVC\`+msvcVersion+`\ATLMFC\LIB\SPECTRE\`+targetA+`;%LIB%"`,
			`  @if exist "%VCINSTALLDIR%\Tools\MSVC\`+msvcVersion+`\LIB\SPECTRE\`+targetA+`" set "LIB=%VCINSTALLDIR%\Tools\MSVC\`+msvcVersion+`\LIB\SPECTRE\`+targetA+`;%LIB%"`,
		)

		base = append(base, "",
			`@if exist "%VCINSTALLDIR%\Tools\MSVC\`+msvcVersion+`\ATLMFC\LIB\SPECTRE\`+targetA+`" set "LIBPATH=%VCINSTALLDIR%\Tools\MSVC\`+msvcVersion+`\ATLMFC\LIB\SPECTRE\`+targetA+`;%LIBPATH%"`,
			`@if exist "%VCINSTALLDIR%\Tools\MSVC\`+msvcVersion+`\LIB\SPECTRE\`+targetA+`" set "LIBPATH=%VCINSTALLDIR%\Tools\MSVC\`+msvcVersion+`\LIB\SPECTRE\`+targetA+`;%LIBPATH%"`,
		)
	}

	if flags.LlvmClang {
		base = append(base, "",
			`@if exist "%VCINSTALLDIR%\Tools\LLVM\`+targetA+`" set "LLVM_PATH=%VCINSTALLDIR%\Tools\LLVM\`+targetA+`"`,
			`@if exist "%LLVM_PATH%\BIN" set "PATH=%LLVM_PATH%\BIN;%PATH%"`,
			`@if exist "%LLVM_PATH%\LIB\CLANG\17\LIB" set "LIB=%LLVM_PATH%\LIB\CLANG\17\LIB;%LIB%"`,
			`@if exist "%LLVM_PATH%\LIB\CLANG\17\LIB" set "LIBPATH=%LLVM_PATH%\LIB\CLANG\17\LIB;%LIBPATH%"`,
			`@if exist "%LLVM_PATH%\LIB\CLANG\17\INCLUDE" set "INCLUDE=%LLVM_PATH%\LIB\CLANG\17\INCLUDE;%INCLUDE%"`)
	}

	return strings.Join(base, "\n")
}
