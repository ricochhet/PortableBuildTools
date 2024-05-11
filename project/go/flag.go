package main

import (
	"flag"

	aflag "github.com/ricochhet/sdkstandalone/flag"
)

var f *aflag.Flags = Newflag()
var defaults = aflag.Flags{
	MSVC_VERSION:          "14.39.17.9",
	MSVC_VERSION_LOCAL:    "14.39.33519",
	SDK_PID:               "Win11SDK_10.0.22621",
	OUTPUT:                "build/sdk_standalone",
	DOWNLOADS:             "build/downloads",
	DOWNLOADS_CRTD:        "build/downloads/crtd",
	DOWNLOADS_DIA:         "build/downloads/dia",
	HOST:                  "x64",
	SET_MSVC_PACKAGES:     "",
	SET_WINSDK_PACKAGES:   "",
	DOWNLOAD_SPECTRE_LIBS: false,
	DOWNLOAD_ARM_TARGETS:  false,
	DOWNLOAD_LLVM_CLANG:   false,
	DOWNLOAD_UNITTEST:     false,
	DOWNLOAD_CMAKE:        false,
	MANIFEST_URL:          "https://aka.ms/vs/17/release/channel",
	MANIFEST_PREVIEW_URL:  "https://aka.ms/vs/17/pre/channel",
	TARGETX64:             "x64",
	TARGETX86:             "x86",
	TARGETARM:             "arm",
	TARGETARM64:           "arm64",
	REWRITE_VARS:          false,
	MSIEXEC_VERBOSE:       false,
}

func Newflag() *aflag.Flags {
	return &aflag.Flags{}
}

func init() {
	f := defaults
	flag.StringVar(&f.MSVC_VERSION, "msvc", defaults.MSVC_VERSION, "Specify MSVC version")
	flag.StringVar(&f.MSVC_VERSION_LOCAL, "msvcv", defaults.MSVC_VERSION_LOCAL, "Specify secondary MSVC version")
	flag.StringVar(&f.SDK_PID, "sdkv", defaults.SDK_PID, "Specify Windows SDK identifier")
	flag.StringVar(&f.OUTPUT, "output", defaults.OUTPUT, "Specify output folder")
	flag.StringVar(&f.DOWNLOADS, "downloads", defaults.DOWNLOADS, "Specify temporary download files folder")
	flag.StringVar(&f.DOWNLOADS_CRTD, "downloads-crtd", defaults.DOWNLOADS_CRTD, "Specify temporary download files folder for CRTD")
	flag.StringVar(&f.DOWNLOADS_DIA, "downloads-dia", defaults.DOWNLOADS_DIA, "Specify temporary download files folder for DIA SDK")
	flag.StringVar(&f.HOST, "host", defaults.HOST, "Specify host architecture (x64 or x86)")
	flag.StringVar(&f.SET_MSVC_PACKAGES, "msvc-packages", defaults.SET_MSVC_PACKAGES, "Specify a list file of MSVC packages to download")
	flag.StringVar(&f.SET_WINSDK_PACKAGES, "sdk-packages", defaults.SET_WINSDK_PACKAGES, "Specify a list file of Windows SDK packages to download")
	flag.BoolVar(&f.DOWNLOAD_SPECTRE_LIBS, "download-spectre-libs", defaults.DOWNLOAD_SPECTRE_LIBS, "Download Spectre libraries")
	flag.BoolVar(&f.DOWNLOAD_ARM_TARGETS, "download-arm-targets", defaults.DOWNLOAD_ARM_TARGETS, "Download ARM targets")
	flag.BoolVar(&f.DOWNLOAD_LLVM_CLANG, "download-llvm-clang", defaults.DOWNLOAD_LLVM_CLANG, "Download LLVM Clang")
	flag.BoolVar(&f.DOWNLOAD_UNITTEST, "download-unittest", defaults.DOWNLOAD_UNITTEST, "Download UnitTest framework")
	flag.BoolVar(&f.DOWNLOAD_CMAKE, "download-cmake", defaults.DOWNLOAD_CMAKE, "Download Cmake build tools")
	flag.StringVar(&f.MANIFEST_URL, "manifest-url", defaults.MANIFEST_URL, "Specify VS manifest url")
	flag.StringVar(&f.MANIFEST_PREVIEW_URL, "manifest-preview-url", defaults.MANIFEST_PREVIEW_URL, "Specify VS preview manifest url")
	flag.BoolVar(&f.REWRITE_VARS, "rewrite-vars", defaults.REWRITE_VARS, "Rewrite environment variable batch scripts")
	flag.BoolVar(&f.MSIEXEC_VERBOSE, "msiexec-verbose", defaults.MSIEXEC_VERBOSE, "Verbose output for rust-msiexec")
	flag.Parse()
}
