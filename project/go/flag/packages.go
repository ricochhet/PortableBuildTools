package flag

import "fmt"

func Maybeappend(msvcPackages, sdkPackages []string, f *Flags) ([]string, []string) {
	msvcAppend := msvcPackages
	sdkAppend := sdkPackages
	if f.DOWNLOAD_LLVM_CLANG {
		msvcAppend = append(msvcAppend, Llvmclangpackages(f)...)
	}

	if f.DOWNLOAD_CMAKE {
		msvcAppend = append(msvcAppend, Cmakepackages(f)...)
	}

	if f.DOWNLOAD_UNITTEST {
		msvcAppend = append(msvcAppend, Unittestpackages(f)...)
	}

	if f.DOWNLOAD_ARM_TARGETS {
		msvcAppend = append(msvcAppend, Msvcarmpackages(f)...)
		sdkAppend = append(sdkAppend, Sdkarmpackages(f)...)
	}

	if f.DOWNLOAD_SPECTRE_LIBS {
		msvcAppend = append(msvcAppend, Msvcspectrepackages(f)...)
		if f.DOWNLOAD_ARM_TARGETS {
			msvcAppend = append(msvcAppend, Msvcarmspectrepackages(f)...)
		}
	}

	return msvcAppend, sdkAppend
}

func Msvcpackages(f *Flags) []string {
	return []string{
		// MSVC vcvars
		"microsoft.visualstudio.vc.vcvars",
		"microsoft.visualstudio.vc.devcmd",
		"microsoft.visualcpp.tools.core.x86",
		fmt.Sprintf("microsoft.visualcpp.tools.host%s.target%s", f.HOST, f.TARGETX64),
		fmt.Sprintf("microsoft.visualcpp.tools.host%s.target%s", f.HOST, f.TARGETX86),
		// MSVC binaries x64
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.base", f.MSVC_VERSION, f.HOST, f.TARGETX64),
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.res.base", f.MSVC_VERSION, f.HOST, f.TARGETX64),
		// MSVC binaries x86
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.base", f.MSVC_VERSION, f.HOST, f.TARGETX86),
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.res.base", f.MSVC_VERSION, f.HOST, f.TARGETX86),
		// MSVC headers
		fmt.Sprintf("microsoft.vc.%s.crt.headers.base", f.MSVC_VERSION),
		// MSVC Libs x64
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.base", f.MSVC_VERSION, f.TARGETX64),
		fmt.Sprintf("microsoft.vc.%s.crt.%s.store.base", f.MSVC_VERSION, f.TARGETX64),
		// MSVC Libs x86
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.base", f.MSVC_VERSION, f.TARGETX86),
		fmt.Sprintf("microsoft.vc.%s.crt.%s.store.base", f.MSVC_VERSION, f.TARGETX86),
		// MSVC runtime source
		fmt.Sprintf("microsoft.vc.%s.crt.source.base", f.MSVC_VERSION),
		// ASAN
		fmt.Sprintf("microsoft.vc.%s.asan.headers.base", f.MSVC_VERSION),
		fmt.Sprintf("microsoft.vc.%s.asan.%s.base", f.MSVC_VERSION, f.TARGETX64),
		fmt.Sprintf("microsoft.vc.%s.asan.%s.base", f.MSVC_VERSION, f.TARGETX86),
		// MSVC tools
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.base", f.MSVC_VERSION, f.HOST, f.TARGETX64),
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.base", f.MSVC_VERSION, f.HOST, f.TARGETX86),
		// DIA SDK
		"microsoft.visualcpp.dia.sdk",
		// MSVC redist
		fmt.Sprintf("microsoft.visualcpp.crt.redist.%s", f.TARGETX64),
		fmt.Sprintf("microsoft.visualcpp.crt.redist.%s", f.TARGETX86),
		// MSVC UnitTest
		"microsoft.visualstudio.vc.unittest.desktop.build.core",
		"microsoft.visualstudio.testtools.codecoverage",
		// MSVC Cli
		fmt.Sprintf("microsoft.vc.%s.cli.%s.base", f.MSVC_VERSION, f.HOST),
		// MSVC Store CRT
		fmt.Sprintf("microsoft.vc.%s.crt.%s.store.base", f.MSVC_VERSION, f.HOST),
		// MSVC Common Tools
		"microsoft.visualcpp.tools.common.utils",
		// MSVC Log
		"microsoft.visualstudio.log",
	}
}

func Msvcarmpackages(f *Flags) []string {
	return []string{
		fmt.Sprintf("microsoft.visualcpp.tools.host%s.target%s", f.HOST, f.TARGETARM),
		fmt.Sprintf("microsoft.visualcpp.tools.host%s.target%s", f.HOST, f.TARGETARM64),
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.base", f.MSVC_VERSION, f.HOST, f.TARGETARM),
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.res.base", f.MSVC_VERSION, f.HOST, f.TARGETARM),
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.base", f.MSVC_VERSION, f.HOST, f.TARGETARM64),
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.res.base", f.MSVC_VERSION, f.HOST, f.TARGETARM64),
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.base", f.MSVC_VERSION, f.TARGETARM),
		fmt.Sprintf("microsoft.vc.%s.crt.%s.store.base", f.MSVC_VERSION, f.TARGETARM),
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.base", f.MSVC_VERSION, f.TARGETARM64),
		fmt.Sprintf("microsoft.vc.%s.crt.%s.store.base", f.MSVC_VERSION, f.TARGETARM64),
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.base", f.MSVC_VERSION, f.HOST, f.TARGETARM),
		fmt.Sprintf("microsoft.vc.%s.tools.host%s.target%s.base", f.MSVC_VERSION, f.HOST, f.TARGETARM64),
	}
}

func Msvcspectrepackages(f *Flags) []string {
	return []string{
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.spectre.base", f.MSVC_VERSION, f.TARGETX64),
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.spectre.base", f.MSVC_VERSION, f.TARGETX86),
	}
}

func Msvcarmspectrepackages(f *Flags) []string {
	return []string{
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.spectre.base", f.MSVC_VERSION, f.TARGETARM),
		fmt.Sprintf("microsoft.vc.%s.crt.%s.desktop.spectre.base", f.MSVC_VERSION, f.TARGETARM64),
	}
}

func Sdkpackages(f *Flags) []string {
	return []string{
		// Windows SDK tools (like rc.exe & mt.exe)
		fmt.Sprintf("Installers\\Windows SDK for Windows Store Apps Tools-%s_en-us.msi", f.TARGETX86),
		// Windows SDK headers
		fmt.Sprintf("Installers\\Windows SDK for Windows Store Apps Headers-%s_en-us.msi", f.TARGETX86),
		fmt.Sprintf("Installers\\Windows SDK Desktop Headers %s-%s_en-us.msi", f.TARGETX86, f.TARGETX86),
		fmt.Sprintf("Installers\\Windows SDK Desktop Headers %s-x86_en-us.msi", f.TARGETX64),
		// Windows SDK libs
		fmt.Sprintf("Installers\\Windows SDK for Windows Store Apps Libs-%s_en-us.msi", f.TARGETX86),
		fmt.Sprintf("Installers\\Windows SDK Desktop Libs %s-%s_en-us.msi", f.TARGETX64, f.TARGETX86),
		fmt.Sprintf("Installers\\Windows SDK Desktop Libs %s-%s_en-us.msi", f.TARGETX86, f.TARGETX86),
		// Windows SDK tools
		fmt.Sprintf("Installers\\Windows SDK Desktop Tools %s-%s_en-us.msi", f.TARGETX64, f.TARGETX86),
		fmt.Sprintf("Installers\\Windows SDK Desktop Tools %s-%s_en-us.msi", f.TARGETX86, f.TARGETX86),
		// CRT headers & libs
		fmt.Sprintf("Installers\\Universal CRT Headers Libraries and Sources-%s_en-us.msi", f.TARGETX86),
		// CRT redist
		fmt.Sprintf("Installers\\Universal CRT Redistributable-%s_en-us.msi", f.TARGETX86),
		// Signing tools
		fmt.Sprintf("Installers\\Windows SDK Signing Tools-%s_en-us.msi", f.TARGETX86),
	}
}

func Sdkarmpackages(f *Flags) []string {
	return []string{
		fmt.Sprintf("Windows SDK Desktop Headers %s-%s_en-us.msi", f.TARGETARM64, f.TARGETX86),
		fmt.Sprintf("Windows SDK Desktop Headers %s-%s_en-us.msi", f.TARGETARM, f.TARGETX86),
		fmt.Sprintf("Windows SDK Desktop Libs %s-%s_en-us.msi", f.TARGETARM64, f.TARGETX86),
		fmt.Sprintf("Windows SDK Desktop Libs %s-%s_en-us.msi", f.TARGETARM, f.TARGETX86),
		fmt.Sprintf("Windows SDK ARM Desktop Tools-%s_en-us.msi", f.TARGETX86),
		fmt.Sprintf("Windows SDK Desktop Tools %s-%s_en-us.msi", f.TARGETARM64, f.TARGETX86),
	}
}

func Llvmclangpackages(f *Flags) []string {
	return []string{
		"microsoft.visualstudio.vc.llvm.base",
		"microsoft.visualstudio.vc.llvm.clang",
	}
}

func Cmakepackages(f *Flags) []string {
	return []string{
		"microsoft.visualstudio.vc.cmake",
	}
}

func Unittestpackages(f *Flags) []string {
	return []string{
		"microsoft.visualstudio.vc.unittest.desktop.build.core",
	}
}
