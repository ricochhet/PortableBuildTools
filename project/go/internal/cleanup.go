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

package internal

import (
	"os"
	"path/filepath"

	aflag "github.com/ricochhet/portablebuildtools/flag"
)

func RemoveVCTipsTelemetry(flags *aflag.Flags) error {
	vctipExe := "vctip.exe"
	msvcv, err := GetMSVCVersion(flags)
	if err != nil { //nolint:wsl // gofumpt conflict
		return err
	}

	paths := []string{
		filepath.Join(flags.Output, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetx64, vctipExe),
		filepath.Join(flags.Output, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetx86, vctipExe),
	}

	if flags.DownloadARMTargets {
		paths = append(paths,
			filepath.Join(flags.Output, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetarm, vctipExe),
			filepath.Join(flags.Output, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetarm64, vctipExe),
		)
	}

	for _, path := range paths {
		err := os.Remove(path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

func CleanHostDirectory(flags *aflag.Flags) error {
	targets := []string{flags.Targetx64, flags.Targetx86, flags.Targetarm, flags.Targetarm64}
	msvcv, err := GetMSVCVersion(flags)
	if err != nil { //nolint:wsl // gofumpt conflict
		return err
	}

	for _, arch := range targets {
		if arch != flags.Host {
			err := os.RemoveAll(filepath.Join(flags.Output, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+arch))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
