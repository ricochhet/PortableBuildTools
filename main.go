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
	"path/filepath"

	"github.com/ricochhet/portablebuildtools/download"
	aflag "github.com/ricochhet/portablebuildtools/flag"
	"github.com/ricochhet/portablebuildtools/internal"
	"github.com/ricochhet/simplezip"
)

//nolint:cyclop // composed of err checking. . .not complex
func main() {
	msvcPackages := aflag.SetPackages(flags, flags.SetMSVCPackages, aflag.MSVCPackages(flags))
	sdkPackages := aflag.SetPackages(flags, flags.SetWinSDKPackages, aflag.WinSDKPackages(flags))

	cwd, err := internal.CreateDirectories(flags)
	if err != nil {
		panic(err)
	}

	flags.Downloads = filepath.Join(cwd, flags.Downloads)
	flags.DownloadsCRTD = filepath.Join(cwd, flags.DownloadsCRTD)
	flags.DownloadsDIA = filepath.Join(cwd, flags.DownloadsDIA)
	flags.Output = filepath.Join(cwd, flags.Output)
	msvcPackages, sdkPackages = aflag.AppendOptionals(msvcPackages, sdkPackages, flags)

	if flags.RewriteVars {
		if err := internal.WriteVars(flags); err != nil {
			panic(err)
		}

		return
	}

	vsManifestJSON, err := download.GetManifest(flags)
	if err != nil {
		panic(err)
	}

	payloads, crtd, dia, sdk := download.GetPackages(flags, vsManifestJSON, msvcPackages)
	if err := download.GetPayloads(flags, payloads); err != nil {
		panic(err)
	}

	if err := download.GetWinSDK(flags, sdk, sdkPackages); err != nil {
		panic(err)
	}

	msvcv, err := internal.GetMSVCVersion(flags)
	if err != nil {
		panic(err)
	}

	destx64 := filepath.Join(flags.Output, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetx64)
	destx86 := filepath.Join(flags.Output, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetx86)
	destarm := filepath.Join(flags.Output, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetarm)
	destarm64 := filepath.Join(flags.Output, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetarm64)

	if err := download.GetCRTD(crtd, destx64, destx86, destarm, destarm64, flags); err != nil {
		panic(err)
	}

	if err := download.GetDIASDK(dia, destx64, destx86, destarm, destarm64, flags); err != nil {
		panic(err)
	}

	if err := internal.RemoveVCTipsTelemetry(flags); err != nil {
		panic(err)
	}

	if err := internal.CleanHostDirectory(flags); err != nil {
		panic(err)
	}

	if err := internal.WriteVars(flags); err != nil {
		panic(err)
	}

	if flags.CreateZip {
		if err := simplezip.Zip(flags.Output, flags.OutputZip); err != nil {
			panic(err)
		}
	}
}
