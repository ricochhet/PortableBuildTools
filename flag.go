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

package main

import (
	"flag"

	aflag "github.com/ricochhet/portablebuildtools/flag"
)

var (
	flags    *aflag.Flags = Newflag()    //nolint:gochecknoglobals // ...
	defaults              = aflag.Flags{ //nolint:gochecknoglobals // ...
		Version:           false,
		MsvcVer:           "14.39.17.9",
		WinSdkVer:         "Win11SDK_10.0.22621",
		Dest:              "build/sdk_standalone",
		DestZip:           "build/sdk_standalone.zip",
		TmpPath:           "build/downloads",
		TmpCrtd:           "build/downloads/crtd",
		TmpDia:            "build/downloads/dia",
		Host:              "x64",
		SetMsvcPackages:   "",
		SetWinSdkPackages: "",
		SpectreLibs:       false,
		ArmTargets:        false,
		LlvmClang:         false,
		UnitTest:          false,
		Cmake:             false,
		MfcAtl:            false,
		ManifestURL:       "https://aka.ms/vs/17/release/channel", // https://aka.ms/vs/17/pre/channel
		Targetx64:         "x64",
		Targetx86:         "x86",
		Targetarm:         "arm",
		Targetarm64:       "arm64",
		WriteEnvironment:  false,
		Verbose:           false,
		Zip:               false,
	}
)

func Newflag() *aflag.Flags {
	return &defaults
}

//nolint:gochecknoinits,lll // cli flags only
func init() {
	flag.BoolVar(&flags.Version, "v", false, "Print the current version")
	flag.StringVar(&flags.MsvcVer, "msvcv", defaults.MsvcVer, "Specify MSVC version")
	flag.StringVar(&flags.WinSdkVer, "sdkv", defaults.WinSdkVer, "Specify Windows SDK identifier")
	flag.StringVar(&flags.Dest, "dest", defaults.Dest, "Specify destination folder")
	flag.StringVar(&flags.DestZip, "dest-zip", defaults.DestZip, "Specify zip archive destination folder")
	flag.StringVar(&flags.TmpPath, "tmp", defaults.TmpPath, "Specify temporary download files folder")
	flag.StringVar(&flags.TmpCrtd, "tmp-crtd", defaults.TmpCrtd, "Specify temporary download files folder for CRTD")
	flag.StringVar(&flags.TmpDia, "tmp-dia", defaults.TmpDia, "Specify temporary download files folder for DIA SDK")
	flag.StringVar(&flags.Host, "host", defaults.Host, "Specify host architecture (x64 or x86)")
	flag.StringVar(&flags.SetMsvcPackages, "msvc-packages", defaults.SetMsvcPackages, "Specify a list file of MSVC packages to download")
	flag.StringVar(&flags.SetWinSdkPackages, "sdk-packages", defaults.SetWinSdkPackages, "Specify a list file of Windows SDK packages to download")
	flag.BoolVar(&flags.SpectreLibs, "spectre-libs", defaults.SpectreLibs, "Download Spectre libraries")
	flag.BoolVar(&flags.ArmTargets, "arm-targets", defaults.ArmTargets, "Download ARM targets")
	flag.BoolVar(&flags.LlvmClang, "llvm-clang", defaults.LlvmClang, "Download LLVM Clang")
	flag.BoolVar(&flags.UnitTest, "unittest", defaults.UnitTest, "Download UnitTest framework")
	flag.BoolVar(&flags.Cmake, "cmake", defaults.Cmake, "Download Cmake build tools")
	flag.BoolVar(&flags.MfcAtl, "mfc-atl", defaults.MfcAtl, "Download MFC/ATL libraries")
	flag.StringVar(&flags.ManifestURL, "manifest-url", defaults.ManifestURL, "Specify VS manifest url")
	flag.BoolVar(&flags.WriteEnvironment, "write-env", defaults.WriteEnvironment, "Write environment variable batch scripts")
	flag.BoolVar(&flags.Verbose, "verbose", defaults.Verbose, "Verbose logging")
	flag.BoolVar(&flags.Zip, "zip", defaults.Zip, "Create zip archive after download")
	flag.Parse()
}
