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

import "fmt"

func AppendOptionals(msvcPackages, sdkPackages []string, flags *Flags) ([]string, []string) {
	msvcAppend := msvcPackages
	sdkAppend := sdkPackages

	if flags.DownloadLLVMClang {
		msvcAppend = append(msvcAppend, LLVMClangPackages()...)
	}

	if flags.DownloadCmake {
		msvcAppend = append(msvcAppend, CmakePackages()...)
	}

	if flags.DownloadUnitTest {
		msvcAppend = append(msvcAppend, UnitTestPackages()...)
	}

	if flags.DownloadARMTargets {
		msvcAppend = append(msvcAppend, MSVCARMPackages(flags)...)
		sdkAppend = append(sdkAppend, WinSDKARMPackages(flags)...)
	}

	if flags.DownloadSpectreLibs {
		msvcAppend = append(msvcAppend, MSVCSpectrePackages(flags)...)
		if flags.DownloadARMTargets {
			msvcAppend = append(msvcAppend, MSVCARMSpectrePackages(flags)...)
		}
	}

	return msvcAppend, sdkAppend
}

func MSVCPackages(flags *Flags) []string {
	return []string{
		// MSVC vcvars
		"microsoft.visualstudio.vc.vcvars",
		"microsoft.visualstudio.vc.devcmd",
		"microsoft.visualcpp.tools.core.x86",
		fmt.Sprintf("microsoft.visualcpp.tools.host%s.target%s", flags.Host, flags.Targetx64),
		fmt.Sprintf("microsoft.visualcpp.tools.host%s.target%s", flags.Host, flags.Targetx86),
		// MSVC binaries x64
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.base", flags.MsvcVer, flags.Host, flags.Targetx64),
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.res.base", flags.MsvcVer, flags.Host, flags.Targetx64),
		// MSVC binaries x86
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.base", flags.MsvcVer, flags.Host, flags.Targetx86),
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.res.base", flags.MsvcVer, flags.Host, flags.Targetx86),
		// MSVC headers
		fmt.Sprintf("microsoft.vc.%s.crt.headers.base", flags.MsvcVer),
		// MSVC Libs x64
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.base", flags.MsvcVer, flags.Targetx64),
		fmt.Sprintf("microsoft.vc.%s.crt.%s.store.base", flags.MsvcVer, flags.Targetx64),
		// MSVC Libs x86
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.base", flags.MsvcVer, flags.Targetx86),
		fmt.Sprintf("microsoft.vc.%s.crt.%s.store.base", flags.MsvcVer, flags.Targetx86),
		// MSVC runtime source
		fmt.Sprintf("microsoft.vc.%s.crt.source.base", flags.MsvcVer),
		// ASAN
		fmt.Sprintf("microsoft.vc.%s.asan.headers.base", flags.MsvcVer),
		fmt.Sprintf("microsoft.vc.%s.asan.%s.base", flags.MsvcVer, flags.Targetx64),
		fmt.Sprintf("microsoft.vc.%s.asan.%s.base", flags.MsvcVer, flags.Targetx86),
		// MSVC tools
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.base", flags.MsvcVer, flags.Host, flags.Targetx64),
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.base", flags.MsvcVer, flags.Host, flags.Targetx86),
		// DIA SDK
		"microsoft.visualcpp.dia.sdk",
		// MSVC redist
		"microsoft.visualcpp.crt.redist." + flags.Targetx64,
		"microsoft.visualcpp.crt.redist." + flags.Targetx86,
		// MSVC UnitTest
		"microsoft.visualstudio.vc.unittest.desktop.build.core",
		"microsoft.visualstudio.testtools.codecoverage",
		// MSVC Cli
		fmt.Sprintf("microsoft.vc.%s.cli.%s.base", flags.MsvcVer, flags.Host),
		// MSVC Store CRT
		fmt.Sprintf("microsoft.vc.%s.crt.%s.store.base", flags.MsvcVer, flags.Host),
		// MSVC Common Tools
		"microsoft.visualcpp.tools.common.utils",
		// MSVC Log
		"microsoft.visualstudio.log",
	}
}

func MSVCARMPackages(flags *Flags) []string {
	return []string{
		fmt.Sprintf("microsoft.visualcpp.tools.host%s.target%s", flags.Host, flags.Targetarm),
		fmt.Sprintf("microsoft.visualcpp.tools.host%s.target%s", flags.Host, flags.Targetarm64),
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.base", flags.MsvcVer, flags.Host, flags.Targetarm),
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.res.base", flags.MsvcVer, flags.Host, flags.Targetarm),
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.base", flags.MsvcVer, flags.Host, flags.Targetarm64),
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.res.base", flags.MsvcVer, flags.Host, flags.Targetarm64),
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.base", flags.MsvcVer, flags.Targetarm),
		fmt.Sprintf("microsoft.vc.%s.crt.%s.store.base", flags.MsvcVer, flags.Targetarm),
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.base", flags.MsvcVer, flags.Targetarm64),
		fmt.Sprintf("microsoft.vc.%s.crt.%s.store.base", flags.MsvcVer, flags.Targetarm64),
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.base", flags.MsvcVer, flags.Host, flags.Targetarm),
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.base", flags.MsvcVer, flags.Host, flags.Targetarm64),
	}
}

func MSVCSpectrePackages(flags *Flags) []string {
	return []string{
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.spectre.base", flags.MsvcVer, flags.Targetx64),
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.spectre.base", flags.MsvcVer, flags.Targetx86),
	}
}

func MSVCARMSpectrePackages(flags *Flags) []string {
	return []string{
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.spectre.base", flags.MsvcVer, flags.Targetarm),
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.spectre.base", flags.MsvcVer, flags.Targetarm64),
	}
}

func WinSDKPackages(flags *Flags) []string {
	return []string{
		// Windows SDK tools (like rc.exe & mt.exe)
		fmt.Sprintf("Installers\\Windows SDK for Windows Store Apps Tools-%s_en-us.msi", flags.Targetx86),
		// Windows SDK headers
		fmt.Sprintf("Installers\\Windows SDK for Windows Store Apps Headers-%s_en-us.msi", flags.Targetx86),
		fmt.Sprintf("Installers\\Windows SDK Desktop Headers %s-%s_en-us.msi", flags.Targetx86, flags.Targetx86),
		fmt.Sprintf("Installers\\Windows SDK Desktop Headers %s-x86_en-us.msi", flags.Targetx64),
		// Windows SDK libs
		fmt.Sprintf("Installers\\Windows SDK for Windows Store Apps Libs-%s_en-us.msi", flags.Targetx86),
		fmt.Sprintf("Installers\\Windows SDK Desktop Libs %s-%s_en-us.msi", flags.Targetx64, flags.Targetx86),
		fmt.Sprintf("Installers\\Windows SDK Desktop Libs %s-%s_en-us.msi", flags.Targetx86, flags.Targetx86),
		// Windows SDK tools
		fmt.Sprintf("Installers\\Windows SDK Desktop Tools %s-%s_en-us.msi", flags.Targetx64, flags.Targetx86),
		fmt.Sprintf("Installers\\Windows SDK Desktop Tools %s-%s_en-us.msi", flags.Targetx86, flags.Targetx86),
		// CRT headers & libs
		fmt.Sprintf("Installers\\Universal CRT Headers Libraries and Sources-%s_en-us.msi", flags.Targetx86),
		// CRT redist
		fmt.Sprintf("Installers\\Universal CRT Redistributable-%s_en-us.msi", flags.Targetx86),
		// Signing tools
		fmt.Sprintf("Installers\\Windows SDK Signing Tools-%s_en-us.msi", flags.Targetx86),
	}
}

func WinSDKARMPackages(flags *Flags) []string {
	return []string{
		fmt.Sprintf("Windows SDK Desktop Headers %s-%s_en-us.msi", flags.Targetarm64, flags.Targetx86),
		fmt.Sprintf("Windows SDK Desktop Headers %s-%s_en-us.msi", flags.Targetarm, flags.Targetx86),
		fmt.Sprintf("Windows SDK Desktop Libs %s-%s_en-us.msi", flags.Targetarm64, flags.Targetx86),
		fmt.Sprintf("Windows SDK Desktop Libs %s-%s_en-us.msi", flags.Targetarm, flags.Targetx86),
		fmt.Sprintf("Windows SDK ARM Desktop Tools-%s_en-us.msi", flags.Targetx86),
		fmt.Sprintf("Windows SDK Desktop Tools %s-%s_en-us.msi", flags.Targetarm64, flags.Targetx86),
	}
}

func LLVMClangPackages() []string {
	return []string{
		"microsoft.visualstudio.vc.llvm.base",
		"microsoft.visualstudio.vc.llvm.clang",
	}
}

func CmakePackages() []string {
	return []string{
		"microsoft.visualstudio.vc.cmake",
	}
}

func UnitTestPackages() []string {
	return []string{
		"microsoft.visualstudio.vc.unittest.desktop.build.core",
	}
}
