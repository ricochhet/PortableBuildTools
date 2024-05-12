package main

import (
	"flag"

	aflag "github.com/ricochhet/sdkstandalone/flag"
)

var (
	flags    *aflag.Flags = Newflag()    //nolint:gochecknoglobals // flags need to be accessed heavily.
	defaults              = aflag.Flags{ //nolint:gochecknoglobals // ^^^^^
		MsvcVer:             "14.39.17.9",
		MsvcVerLocal:        "14.39.33519",
		WinSDKVer:           "Win11SDK_10.0.22621",
		Output:              "build/sdk_standalone",
		Downloads:           "build/downloads",
		DownloadsCRTD:       "build/downloads/crtd",
		DownloadsDIA:        "build/downloads/dia",
		Host:                "x64",
		SetMSVCPackages:     "",
		SetWinSDKPackages:   "",
		DownloadSpectreLibs: false,
		DownloadARMTargets:  false,
		DownloadLLVMClang:   false,
		DownloadUnitTest:    false,
		DownloadCmake:       false,
		ManifestURL:         "https://aka.ms/vs/17/release/channel",
		ManifestPreviewURL:  "https://aka.ms/vs/17/pre/channel",
		Targetx64:           "x64",
		Targetx86:           "x86",
		Targetarm:           "arm",
		Targetarm64:         "arm64",
		RewriteVars:         false,
		MSIExtractVerbose:   false,
	}
)

func Newflag() *aflag.Flags {
	return &aflag.Flags{} //nolint:exhaustruct // intentionally leave struct non-exhausted.
}

//nolint:gochecknoinits,lll // cli flags only
func init() {
	flags := defaults
	flag.StringVar(&flags.MsvcVer, "msvc", defaults.MsvcVer, "Specify MSVC version")
	flag.StringVar(&flags.MsvcVerLocal, "msvcv", defaults.MsvcVerLocal, "Specify secondary MSVC version")
	flag.StringVar(&flags.WinSDKVer, "sdkv", defaults.WinSDKVer, "Specify Windows SDK identifier")
	flag.StringVar(&flags.Output, "output", defaults.Output, "Specify output folder")
	flag.StringVar(&flags.Downloads, "downloads", defaults.Downloads, "Specify temporary download files folder")
	flag.StringVar(&flags.DownloadsCRTD, "downloads-crtd", defaults.DownloadsCRTD, "Specify temporary download files folder for CRTD")
	flag.StringVar(&flags.DownloadsDIA, "downloads-dia", defaults.DownloadsDIA, "Specify temporary download files folder for DIA SDK")
	flag.StringVar(&flags.Host, "host", defaults.Host, "Specify host architecture (x64 or x86)")
	flag.StringVar(&flags.SetMSVCPackages, "msvc-packages", defaults.SetMSVCPackages, "Specify a list file of MSVC packages to download")
	flag.StringVar(&flags.SetWinSDKPackages, "sdk-packages", defaults.SetWinSDKPackages, "Specify a list file of Windows SDK packages to download")
	flag.BoolVar(&flags.DownloadSpectreLibs, "download-spectre-libs", defaults.DownloadSpectreLibs, "Download Spectre libraries")
	flag.BoolVar(&flags.DownloadARMTargets, "download-arm-targets", defaults.DownloadARMTargets, "Download ARM targets")
	flag.BoolVar(&flags.DownloadLLVMClang, "download-llvm-clang", defaults.DownloadLLVMClang, "Download LLVM Clang")
	flag.BoolVar(&flags.DownloadUnitTest, "download-unittest", defaults.DownloadUnitTest, "Download UnitTest framework")
	flag.BoolVar(&flags.DownloadCmake, "download-cmake", defaults.DownloadCmake, "Download Cmake build tools")
	flag.StringVar(&flags.ManifestURL, "manifest-url", defaults.ManifestURL, "Specify VS manifest url")
	flag.StringVar(&flags.ManifestPreviewURL, "manifest-preview-url", defaults.ManifestPreviewURL, "Specify VS preview manifest url")
	flag.BoolVar(&flags.RewriteVars, "rewrite-vars", defaults.RewriteVars, "Rewrite environment variable batch scripts")
	flag.BoolVar(&flags.MSIExtractVerbose, "msiexec-verbose", defaults.MSIExtractVerbose, "Verbose output for rust-msiexec")
	flag.Parse()
}
