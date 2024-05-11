package main

import (
	"path/filepath"

	"github.com/ricochhet/sdkstandalone/download"
	aflag "github.com/ricochhet/sdkstandalone/flag"
)

func main() {
	msvcPackages := aflag.Setpackages(f, f.SET_MSVC_PACKAGES, aflag.Msvcpackages(f))
	sdkPackages := aflag.Setpackages(f, f.SET_WINSDK_PACKAGES, aflag.Sdkpackages(f))

	wd, err := download.Createdirectories(f)
	if err != nil {
		panic(err)
	}
	f.DOWNLOADS = filepath.Join(wd, f.DOWNLOADS)
	f.DOWNLOADS_CRTD = filepath.Join(wd, f.DOWNLOADS_CRTD)
	f.DOWNLOADS_DIA = filepath.Join(wd, f.DOWNLOADS_DIA)
	f.OUTPUT = filepath.Join(wd, f.OUTPUT)
	msvcPackages, sdkPackages = aflag.Maybeappend(msvcPackages, sdkPackages, f)

	if f.REWRITE_VARS {
		err := download.Writevars(f)
		if err != nil {
			panic(err)
		}
		return
	}

	vsmanifestjson, err := download.Getmanifest(f)
	if err != nil {
		panic(err)
	}

	payloads, crtd, dia, sdk := download.Getpackages(f, vsmanifestjson, msvcPackages)
	download.Getpayloads(f, payloads)
	err = download.Getwinsdk(f, sdk, sdkPackages)
	if err != nil {
		panic(err)
	}

	dstX64 := filepath.Join(f.OUTPUT, "VC", "Tools", "MSVC", f.MSVC_VERSION_LOCAL, "bin", "Host"+f.HOST, f.TARGETX64)
	dstX86 := filepath.Join(f.OUTPUT, "VC", "Tools", "MSVC", f.MSVC_VERSION_LOCAL, "bin", "Host"+f.HOST, f.TARGETX86)
	dstARM := filepath.Join(f.OUTPUT, "VC", "Tools", "MSVC", f.MSVC_VERSION_LOCAL, "bin", "Host"+f.HOST, f.TARGETARM)
	dstARM64 := filepath.Join(f.OUTPUT, "VC", "Tools", "MSVC", f.MSVC_VERSION_LOCAL, "bin", "Host"+f.HOST, f.TARGETARM64)
	err = download.Getcrtd(crtd, dstX64, dstX86, dstARM, dstARM64, f)
	if err != nil {
		panic(err)
	}

	err = download.Getdiasdk(dia, dstX64, dstX86, dstARM, dstARM64, f)
	if err != nil {
		panic(err)
	}

	err = download.Removetelemetry(f)
	if err != nil {
		panic(err)
	}

	err = download.Cleanhost(f)
	if err != nil {
		panic(err)
	}

	err = download.Writevars(f)
	if err != nil {
		panic(err)
	}
}
