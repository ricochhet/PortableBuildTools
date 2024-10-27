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

func AppendOptionals(msvcPackages, sdkPackages []string, flags *Flags) ([]string, []string) { //nolint:cyclop // ...
	msvcAppend := msvcPackages
	sdkAppend := sdkPackages

	if flags.LlvmClang {
		msvcAppend = append(msvcAppend, LlvmClangPackages()...)
	}

	if flags.Cmake {
		msvcAppend = append(msvcAppend, CmakePackages()...)
	}

	if flags.UnitTest {
		msvcAppend = append(msvcAppend, UnitTestPackages()...)
	}

	if flags.ArmTargets {
		msvcAppend = append(msvcAppend, MsvcArmPackages(flags)...)
		sdkAppend = append(sdkAppend, WinSdkArmPackages(flags)...)
	}

	if flags.SpectreLibs {
		msvcAppend = append(msvcAppend, MsvcSpectrePackages(flags)...)
		if flags.ArmTargets {
			msvcAppend = append(msvcAppend, MsvcArmSpectrePackages(flags)...)
		}
	}

	if flags.MfcAtl {
		msvcAppend = append(msvcAppend, MfcAtlPackages(flags)...)

		if flags.ArmTargets {
			msvcAppend = append(msvcAppend, MfcAtlArmPackages(flags)...)
		}

		if flags.SpectreLibs {
			msvcAppend = append(msvcAppend, MfcAtlSpectrePackages(flags)...)

			if flags.ArmTargets {
				msvcAppend = append(msvcAppend, MfcAtlSpectreArmPackages(flags)...)
			}
		}
	}

	if flags.Vcpkg {
		msvcAppend = append(msvcAppend, VcpkgPackages(flags)...)
	}

	if flags.Msbuild {
		msvcAppend = append(msvcAppend, MsbuildPackages(flags)...)

		if flags.ArmTargets {
			msvcAppend = append(msvcAppend, MsbuildArmPackages(flags)...)
		}
	}

	return msvcAppend, sdkAppend
}

func MsvcPackages(flags *Flags) []string {
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
		// MSVC Setup
		"microsoft.visualstudio.setup.configuration",
	}
}

func MsvcArmPackages(flags *Flags) []string {
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

func MsvcSpectrePackages(flags *Flags) []string {
	return []string{
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.spectre.base", flags.MsvcVer, flags.Targetx64),
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.spectre.base", flags.MsvcVer, flags.Targetx86),
	}
}

func MsvcArmSpectrePackages(flags *Flags) []string {
	return []string{
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.spectre.base", flags.MsvcVer, flags.Targetarm),
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.spectre.base", flags.MsvcVer, flags.Targetarm64),
	}
}

func WinSdkPackages(flags *Flags) []string {
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

func WinSdkArmPackages(flags *Flags) []string {
	return []string{
		fmt.Sprintf("Windows SDK Desktop Headers %s-%s_en-us.msi", flags.Targetarm64, flags.Targetx86),
		fmt.Sprintf("Windows SDK Desktop Headers %s-%s_en-us.msi", flags.Targetarm, flags.Targetx86),
		fmt.Sprintf("Windows SDK Desktop Libs %s-%s_en-us.msi", flags.Targetarm64, flags.Targetx86),
		fmt.Sprintf("Windows SDK Desktop Libs %s-%s_en-us.msi", flags.Targetarm, flags.Targetx86),
		fmt.Sprintf("Windows SDK ARM Desktop Tools-%s_en-us.msi", flags.Targetx86),
		fmt.Sprintf("Windows SDK Desktop Tools %s-%s_en-us.msi", flags.Targetarm64, flags.Targetx86),
	}
}

func LlvmClangPackages() []string {
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

func MfcAtlPackages(flags *Flags) []string {
	return []string{
		// MFC
		fmt.Sprintf("microsoft.vc.%s.mfc.headers.base", flags.MsvcVer),
		fmt.Sprintf("microsoft.vc.%s.mfc.%s.base", flags.MsvcVer, flags.Targetx64),
		fmt.Sprintf("microsoft.vc.%s.mfc.%s.base", flags.MsvcVer, flags.Targetx86),
		// MFC MBCS (Multibyte Character Set)
		fmt.Sprintf("microsoft.vc.%s.mfc.mbcs.base", flags.MsvcVer),
		fmt.Sprintf("microsoft.vc.%s.mfc.mbcs.%s.base", flags.MsvcVer, flags.Targetx64),
		// MFC source
		fmt.Sprintf("microsoft.vc.%s.mfc.source.base", flags.MsvcVer),
		// ATL
		fmt.Sprintf("microsoft.vc.%s.atl.headers.base", flags.MsvcVer),
		fmt.Sprintf("microsoft.vc.%s.atl.%s.base", flags.MsvcVer, flags.Targetx64),
		fmt.Sprintf("microsoft.vc.%s.atl.%s.base", flags.MsvcVer, flags.Targetx86),
		// ATL source
		fmt.Sprintf("microsoft.vc.%s.atl.source.base", flags.MsvcVer),
	}
}

func MfcAtlSpectrePackages(flags *Flags) []string {
	return []string{
		fmt.Sprintf("microsoft.vc.%s.mfc.%s.spectre.base", flags.MsvcVer, flags.Targetx64),
		fmt.Sprintf("microsoft.vc.%s.mfc.%s.spectre.base", flags.MsvcVer, flags.Targetx86),

		fmt.Sprintf("microsoft.vc.%s.mfc.mbcs.spectre.base", flags.MsvcVer),
		fmt.Sprintf("microsoft.vc.%s.mfc.mbcs.%s.spectre.base", flags.MsvcVer, flags.Targetx64),

		fmt.Sprintf("microsoft.vc.%s.atl.%s.spectre.base", flags.MsvcVer, flags.Targetx64),
		fmt.Sprintf("microsoft.vc.%s.atl.%s.spectre.base", flags.MsvcVer, flags.Targetx86),
	}
}

func MfcAtlSpectreArmPackages(flags *Flags) []string {
	return []string{
		fmt.Sprintf("microsoft.vc.%s.mfc.%s.spectre.base", flags.MsvcVer, flags.Targetarm64),
		fmt.Sprintf("microsoft.vc.%s.mfc.%s.spectre.base", flags.MsvcVer, flags.Targetarm),

		fmt.Sprintf("microsoft.vc.%s.mfc.mbcs.%s.spectre.base", flags.MsvcVer, flags.Targetarm64),
		fmt.Sprintf("microsoft.vc.%s.mfc.mbcs.%s.spectre.base", flags.MsvcVer, flags.Targetarm),

		fmt.Sprintf("microsoft.vc.%s.atl.%s.spectre.base", flags.MsvcVer, flags.Targetarm64),
		fmt.Sprintf("microsoft.vc.%s.atl.%s.spectre.base", flags.MsvcVer, flags.Targetarm),
	}
}

func MfcAtlArmPackages(flags *Flags) []string {
	return []string{
		// MFC
		fmt.Sprintf("microsoft.vc.%s.mfc.%s.base", flags.MsvcVer, flags.Targetarm64),
		fmt.Sprintf("microsoft.vc.%s.mfc.%s.base", flags.MsvcVer, flags.Targetarm),
		// MFC MBCS (Multibyte Character Set)
		fmt.Sprintf("microsoft.vc.%s.mfc.mbcs.%s.base", flags.MsvcVer, flags.Targetarm64),
		fmt.Sprintf("microsoft.vc.%s.mfc.mbcs.%s.base", flags.MsvcVer, flags.Targetarm),
		// ATL
		fmt.Sprintf("microsoft.vc.%s.atl.%s.base", flags.MsvcVer, flags.Targetarm64),
		fmt.Sprintf("microsoft.vc.%s.atl.%s.base", flags.MsvcVer, flags.Targetarm),
	}
}

func MsbuildPackages(flags *Flags) []string {
	return []string{
		"microsoft.build",
		"microsoft.build.dependencies",

		"microsoft.visualc.140.msbuild.base.msi",
		"microsoft.visualc.140.msbuild.base.msi.resources",

		fmt.Sprintf("microsoft.visualc.140.msbuild.%s.msi", flags.Targetx64),
		fmt.Sprintf("microsoft.visualc.140.msbuild.%s.msi", flags.Targetx86),

		"microsoft.visualstudio.vc.msbuild.base",
		"microsoft.visualstudio.vc.msbuild.base.resources",

		"microsoft.visualstudio.vc.msbuild.base.uwp",

		"microsoft.visualstudio.vc.msbuild.llvm",
		"microsoft.visualstudio.vc.msbuild.llvm.resources",

		"microsoft.visualstudio.vc.msbuild.v150.base",
		"microsoft.visualstudio.vc.msbuild.v150.base.resources",
		"microsoft.visualstudio.vc.msbuild.v150.uwp",

		"microsoft.visualstudio.vc.msbuild.v150." + flags.Targetx64,
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.v150.%s.v141", flags.Targetx64),
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.v150.%s.v141_xp", flags.Targetx64),

		"microsoft.visualstudio.vc.msbuild.v150.%s" + flags.Targetx86,
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.v150.%s.v141", flags.Targetx86),
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.v150.%s.v141_xp", flags.Targetx86),

		"microsoft.visualstudio.vc.msbuild.v170.base",
		"microsoft.visualstudio.vc.msbuild.v170.base.resources",

		"microsoft.visualstudio.vc.msbuild.v170." + flags.Targetx64,
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.v170.%s.uwp", flags.Targetx64),
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.v170.%s.v143", flags.Targetx64),

		"microsoft.visualstudio.vc.msbuild.v170." + flags.Targetx86,
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.v170.%s.uwp", flags.Targetx86),
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.v170.%s.v143", flags.Targetx86),

		"microsoft.visualstudio.vc.msbuild." + flags.Targetx64,
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.%s.uwp", flags.Targetx64),
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.%s.v142", flags.Targetx64),

		"microsoft.visualstudio.vc.msbuild." + flags.Targetx86,
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.%s.uwp", flags.Targetx86),
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.%s.v142", flags.Targetx86),
	}
}

func MsbuildArmPackages(flags *Flags) []string {
	return []string{
		fmt.Sprintf("microsoft.visualc.140.msbuild.%s.msi", flags.Targetarm),

		"microsoft.visualstudio.vc.msbuild." + flags.Targetarm64,
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.%s.uwp", flags.Targetarm64),
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.%s.v142", flags.Targetarm64),

		"microsoft.visualstudio.vc.msbuild." + flags.Targetarm,
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.%s.uwp", flags.Targetarm),
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.%s.v142", flags.Targetarm),

		"microsoft.visualstudio.vc.msbuild.v150." + flags.Targetarm64,
		"microsoft.visualstudio.vc.msbuild.v150." + flags.Targetarm,

		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.v150.%s.v141", flags.Targetarm64),
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.v150.%s.v141", flags.Targetarm),

		"microsoft.visualstudio.vc.msbuild.v170." + flags.Targetarm64,
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.v170.%s.uwp", flags.Targetarm64),
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.v170.%s.v143", flags.Targetarm64),

		"microsoft.visualstudio.vc.msbuild.v170." + flags.Targetarm,
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.v170.%s.uwp", flags.Targetarm),
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.v170.%s.v143", flags.Targetarm),

		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.v170.%sec", flags.Targetarm64),
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.v170.%sec.uwp", flags.Targetarm64),
		fmt.Sprintf("microsoft.visualstudio.vc.msbuild.v170.%sec.v143", flags.Targetarm64),
	}
}

func VcpkgPackages(_ *Flags) []string {
	return []string{
		"microsoft.visualstudio.vc.vcpkg",
	}
}
