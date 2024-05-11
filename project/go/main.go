package main

import (
	"flag"
	"path/filepath"

	"github.com/ricochhet/sdkstandalone/download"
	aflag "github.com/ricochhet/sdkstandalone/flag"
)

var defaults = aflag.Flags{
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

func main() {
	f := defaults
	flag.StringVar(&f.MSVC_VERSION, "msvc", defaults.MSVC_VERSION, "Specify MSVC version")
	flag.StringVar(&f.MSVC_VERSION_LOCAL, "msvcv", defaults.MSVC_VERSION_LOCAL, "Specificy secondary MSVC version")
	flag.StringVar(&f.SDK_PID, "sdkv", defaults.SDK_PID, "Specify Windows SDK identifier")
	flag.StringVar(&f.OUTPUT, "output", defaults.OUTPUT, "Specify output folder")
	flag.StringVar(&f.DOWNLOADS, "downloads", defaults.DOWNLOADS, "Specify temporary download files folder")
	flag.StringVar(&f.DOWNLOADS_CRTD, "downloads-crtd", defaults.DOWNLOADS_CRTD, "Specify temporary download files folder for CRTD")
	flag.StringVar(&f.DOWNLOADS_DIA, "downloads-dia", defaults.DOWNLOADS_DIA, "Specify temporary download files folder for DIA SDK")
	flag.StringVar(&f.HOST, "host", defaults.HOST, "Specify host architecture (x64 or x86)")
	flag.BoolVar(&f.DOWNLOAD_SPECTRE_LIBS, "download-spectre-libs", defaults.DOWNLOAD_SPECTRE_LIBS, "Download Spectre libraries")
	flag.BoolVar(&f.DOWNLOAD_ARM_TARGETS, "download-arm-targets", defaults.DOWNLOAD_ARM_TARGETS, "Download ARM targets")
	flag.StringVar(&f.MANIFEST_URL, "manifest-url", defaults.MANIFEST_URL, "Specify VS manifest url")
	flag.StringVar(&f.MANIFEST_PREVIEW_URL, "manifest-preview-url", defaults.MANIFEST_PREVIEW_URL, "Specify VS preview manifest url")
	flag.BoolVar(&f.REWRITE_VARS, "rewrite-vars", defaults.REWRITE_VARS, "Rewrite environment variable batch scripts")
	flag.Parse()

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
		err := download.Writevars(&f)
		if err != nil {
			panic(err)
		}
		return
	}

	vsmanifestjson, err := download.Getmanifest(&f)
	if err != nil {
		panic(err)
	}

	payloads, crtd, dia, sdk := download.Getpackages(&f, vsmanifestjson, msvcPackages)
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

	err = download.Writevars(&f)
	if err != nil {
		panic(err)
	}
}
