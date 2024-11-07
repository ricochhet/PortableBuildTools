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

package download

import (
	"io/fs"
	"os"
	"path/filepath"

	acopy "github.com/otiai10/copy"
	"github.com/ricochhet/minicommon/charmbracelet"
	"github.com/ricochhet/minicommon/download"
	aflag "github.com/ricochhet/portablebuildtools/flag"
	"github.com/ricochhet/portablebuildtools/internal"
	"github.com/tidwall/gjson"
)

func GetCrtd(payloads []string, destx64, destx86, destarm, destarm64 string, flags *aflag.Flags) error {
	for _, payload := range payloads {
		packages := gjson.Parse(payload).Array()
		for _, pkg := range packages {
			url := pkg.Get("url").String()
			sha256 := pkg.Get("sha256").String()
			fileName := pkg.Get("fileName").String()

			if err := download.FileValidated(url, sha256, fileName, flags.TmpPath); err != nil {
				charmbracelet.SharedLogger.Errorf("Error downloading CRTD package: %v", err)
				continue
			}
		}
	}

	if err := internal.ExtractMsi(flags, filepath.Join(flags.TmpPath, "vc_RuntimeDebug.msi"), flags.TmpPath); err != nil {
		return err
	}

	dlls, err := os.ReadDir(filepath.Join(flags.TmpPath, "System64"))
	if err != nil {
		return err
	}

	for _, dll := range dlls {
		paths := []copyCrtdPath{
			{dest: destx64, flags: flags},
			{dest: destx86, flags: flags},
		}

		if flags.ArmTargets {
			paths = append(paths,
				copyCrtdPath{dest: destarm, flags: flags},
				copyCrtdPath{dest: destarm64, flags: flags},
			)
		}

		if err := copyCrtdToPaths(dll, paths); err != nil {
			return err
		}
	}

	return nil
}

type copyCrtdPath struct {
	dest  string
	flags *aflag.Flags
}

func copyCrtdToPaths(dll fs.DirEntry, paths []copyCrtdPath) error {
	for _, path := range paths {
		if err := CopyCrtd(dll, path.dest, path.flags); err != nil {
			return err
		}
	}

	return nil
}

func CopyCrtd(dirEntry fs.DirEntry, target string, flags *aflag.Flags) error {
	return acopy.Copy(filepath.Join(filepath.Join(flags.TmpPath, "System64"), dirEntry.Name()), filepath.Join(target, dirEntry.Name()))
}
