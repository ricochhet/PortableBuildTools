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

	"github.com/ricochhet/minicommon/zip"
	"github.com/ricochhet/portablebuildtools/download"
	aflag "github.com/ricochhet/portablebuildtools/flag"
	"github.com/ricochhet/portablebuildtools/internal"
)

//nolint:cyclop,gocyclo // wintfix
func main() {
	if _, _, err := internal.FindMsiExtract(); err != nil {
		panic(err)
	}

	msvcPackages := aflag.SetPackages(flags, flags.SetMsvcPackages, aflag.MsvcPackages(flags))
	sdkPackages := aflag.SetPackages(flags, flags.SetWinSdkPackages, aflag.WinSdkPackages(flags))

	cwd, err := internal.CreateDirectories(flags)
	if err != nil {
		panic(err)
	}

	flags.TmpPath = filepath.Join(cwd, flags.TmpPath)
	flags.TmpCrtd = filepath.Join(cwd, flags.TmpCrtd)
	flags.TmpDia = filepath.Join(cwd, flags.TmpDia)
	flags.Dest = filepath.Join(cwd, flags.Dest)
	msvcPackages, sdkPackages = aflag.AppendOptionals(msvcPackages, sdkPackages, flags)

	if flags.WriteEnvironment {
		if err := internal.WriteEnvironment(flags); err != nil {
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

	if err := download.GetWinSdk(flags, sdk, sdkPackages); err != nil {
		panic(err)
	}

	msvcv, err := internal.GetMsvcVersion(flags)
	if err != nil {
		panic(err)
	}

	destx64 := filepath.Join(flags.Dest, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetx64)
	destx86 := filepath.Join(flags.Dest, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetx86)
	destarm := filepath.Join(flags.Dest, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetarm)
	destarm64 := filepath.Join(flags.Dest, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetarm64)

	if err := download.GetCrtd(crtd, destx64, destx86, destarm, destarm64, flags); err != nil {
		panic(err)
	}

	if err := download.GetDiaSdk(dia, destx64, destx86, destarm, destarm64, flags); err != nil {
		panic(err)
	}

	if err := internal.RemoveVcTipsTelemetry(flags); err != nil {
		panic(err)
	}

	if err := internal.CleanHostDirectory(flags); err != nil {
		panic(err)
	}

	if err := internal.WriteEnvironment(flags); err != nil {
		panic(err)
	}

	if flags.Zip {
		if err := zip.Zip(flags.Dest, flags.DestZip); err != nil {
			panic(err)
		}
	}
}
