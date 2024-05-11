package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/ricochhet/sdkstandalone/download"
	aflag "github.com/ricochhet/sdkstandalone/flag"
)

func main() {
	f := aflag.Flags{
		MSVC_VERSION:          "14.39.17.9",
		MSVC_VERSION_LOCAL:    "14.39.33519",
		SDK_PID:               "Win11SDK_10.0.22621",
		OUTPUT:                "build/sdk_standalone",
		DOWNLOADS:             "build/downloads",
		DOWNLOADS_CRTD:        "build/downloads/crtd",
		DOWNLOADS_DIA:         "build/downloads/dia",
		HOST:                  "x64",
		DOWNLOAD_SPECTRE_LIBS: false,
		DOWNLOAD_ARM_TARGETS:  false,
		MANIFEST_URL:          "https://aka.ms/vs/17/release/channel",
		MANIFEST_PREVIEW_URL:  "https://aka.ms/vs/17/pre/channel",
		TARGETX64:             "x64",
		TARGETX86:             "x86",
		TARGETARM:             "arm",
		TARGETARM64:           "arm64",
		REWRITE_VARS:          false,
	}

	msvcVersion := flag.String("msvc", f.MSVC_VERSION, "Specify MSVC version")
	msvcVersionLocal := flag.String("msvcv", f.MSVC_VERSION_LOCAL, "Specificy secondary MSVC version")
	sdkPid := flag.String("sdkv", f.SDK_PID, "Specify Windows SDK identifier")
	host := flag.String("host", f.HOST, "Specify host architecture (x64 or x86)")
	downloadSpectreLibs := flag.Bool("download-spectre-libs", f.DOWNLOAD_SPECTRE_LIBS, "Download Spectre libraries")
	downloadArmTargets := flag.Bool("download-arm-targets", f.DOWNLOAD_ARM_TARGETS, "Download ARM targets")
	output := flag.String("output", f.OUTPUT, "Specify output folder")
	downloads := flag.String("downloads", f.DOWNLOADS, "Specify temporary download files folder")
	rewriteVars := flag.Bool("rewrite-vars", f.REWRITE_VARS, "Rewrite environment variable batch scripts")
	flag.Parse()

	f.MSVC_VERSION = *msvcVersion
	f.MSVC_VERSION_LOCAL = *msvcVersionLocal
	f.SDK_PID = *sdkPid
	f.HOST = *host
	f.DOWNLOAD_SPECTRE_LIBS = *downloadSpectreLibs
	f.DOWNLOAD_ARM_TARGETS = *downloadArmTargets
	f.OUTPUT = *output
	f.DOWNLOADS = *downloads
	f.REWRITE_VARS = *rewriteVars

	msvcPackages := aflag.Msvcpackages(&f)
	msvcARMPackages := aflag.Msvcarmpackages(&f)
	msvcSpectrePackages := aflag.Msvcspectrepackages(&f)
	msvcARMSpectrePackages := aflag.Msvcarmspectrepackages(&f)
	sdkPackages := aflag.Sdkpackages(&f)
	sdkARMPackages := aflag.Sdkarmpackages(&f)

	wd, err := download.Createdirectories(&f)
	if err != nil {
		panic(err)
	}
	f.DOWNLOADS = filepath.Join(wd, f.DOWNLOADS)
	f.DOWNLOADS_CRTD = filepath.Join(wd, f.DOWNLOADS_CRTD)
	f.DOWNLOADS_DIA = filepath.Join(wd, f.DOWNLOADS_DIA)
	f.OUTPUT = filepath.Join(wd, f.OUTPUT)

	if f.DOWNLOAD_ARM_TARGETS {
		msvcPackages = append(msvcPackages, msvcARMPackages...)
		sdkPackages = append(sdkPackages, sdkARMPackages...)
	}

	if f.DOWNLOAD_SPECTRE_LIBS {
		msvcPackages = append(msvcPackages, msvcSpectrePackages...)
		if f.DOWNLOAD_ARM_TARGETS {
			msvcPackages = append(msvcPackages, msvcARMSpectrePackages...)
		}
	}

	if f.REWRITE_VARS {
		sdkv, err := download.Getwinsdkversion(&f)
		if err != nil {
			panic(err)
		}

		os.WriteFile(filepath.Join(f.OUTPUT, "set_vars64.bat"), []byte(aflag.SetX64(f.MSVC_VERSION_LOCAL, sdkv, f.TARGETX64, f.TARGETX86, f.HOST)), 0644)
		os.WriteFile(filepath.Join(f.OUTPUT, "set_vars32.bat"), []byte(aflag.SetX86(f.MSVC_VERSION_LOCAL, sdkv, f.TARGETX86, f.TARGETX64, f.HOST)), 0644)
		os.WriteFile(filepath.Join(f.OUTPUT, "set_vars_arm64.bat"), []byte(aflag.SetX64(f.MSVC_VERSION_LOCAL, sdkv, f.TARGETARM64, f.TARGETARM, f.HOST)), 0644)
		os.WriteFile(filepath.Join(f.OUTPUT, "set_vars_arm32.bat"), []byte(aflag.SetX86(f.MSVC_VERSION_LOCAL, sdkv, f.TARGETARM, f.TARGETARM64, f.HOST)), 0644)
		return
	}

	vsmanifest, err := download.Getmanifest(&f)
	if err != nil {
		panic(err)
	}

	payloads, crtd, dia, sdk := download.Getpackages(&f, vsmanifest, msvcPackages)
	download.Getpayloads(&f, payloads)
	err = download.Getwinsdk(&f, sdk, sdkPackages)
	if err != nil {
		panic(err)
	}

	dstX64 := filepath.Join(f.OUTPUT, "VC", "Tools", "MSVC", f.MSVC_VERSION_LOCAL, "bin", "Host"+f.HOST, f.TARGETX64)
	dstX86 := filepath.Join(f.OUTPUT, "VC", "Tools", "MSVC", f.MSVC_VERSION_LOCAL, "bin", "Host"+f.HOST, f.TARGETX86)
	dstARM := filepath.Join(f.OUTPUT, "VC", "Tools", "MSVC", f.MSVC_VERSION_LOCAL, "bin", "Host"+f.HOST, f.TARGETARM)
	dstARM64 := filepath.Join(f.OUTPUT, "VC", "Tools", "MSVC", f.MSVC_VERSION_LOCAL, "bin", "Host"+f.HOST, f.TARGETARM64)
	err = download.Getcrtd(crtd, dstX64, dstX86, dstARM, dstARM64, &f)
	if err != nil {
		panic(err)
	}

	err = download.Getdiasdk(dia, dstX64, dstX86, dstARM, dstARM64, &f)
	if err != nil {
		panic(err)
	}

	err = download.Removetelemetry(&f)
	if err != nil {
		panic(err)
	}

	err = download.Cleanhost(&f)
	if err != nil {
		panic(err)
	}

	sdkv, err := download.Getwinsdkversion(&f)
	if err != nil {
		panic(err)
	}

	os.WriteFile(filepath.Join(f.OUTPUT, "set_vars64.bat"), []byte(aflag.SetX64(f.MSVC_VERSION_LOCAL, sdkv, f.TARGETX64, f.TARGETX86, f.HOST)), 0644)
	os.WriteFile(filepath.Join(f.OUTPUT, "set_vars32.bat"), []byte(aflag.SetX86(f.MSVC_VERSION_LOCAL, sdkv, f.TARGETX86, f.TARGETX64, f.HOST)), 0644)
	os.WriteFile(filepath.Join(f.OUTPUT, "set_vars_arm64.bat"), []byte(aflag.SetX64(f.MSVC_VERSION_LOCAL, sdkv, f.TARGETARM64, f.TARGETARM, f.HOST)), 0644)
	os.WriteFile(filepath.Join(f.OUTPUT, "set_vars_arm32.bat"), []byte(aflag.SetX86(f.MSVC_VERSION_LOCAL, sdkv, f.TARGETARM, f.TARGETARM64, f.HOST)), 0644)
}
