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
	"os"
	"path/filepath"

	acopy "github.com/otiai10/copy"
	"github.com/ricochhet/minicommon/charmbracelet"
	"github.com/ricochhet/minicommon/download"
	"github.com/ricochhet/minicommon/filesystem"
	aflag "github.com/ricochhet/portablebuildtools/flag"
	"github.com/ricochhet/portablebuildtools/internal"
	"github.com/tidwall/gjson"
)

//nolint:cyclop // wontfix
func GetPayloads(flags *aflag.Flags, payloads []string) error {
	for _, payload := range payloads {
		packages := gjson.Parse(payload).Array()
		tmpPath := flags.TmpPath
		cabinetCount := 0
		storedMsi := ""

	outer:
		for _, pkg := range packages {
			url := gjson.Get(pkg.String(), "url").String()
			sha256 := gjson.Get(pkg.String(), "sha256").String()
			fileName := gjson.Get(pkg.String(), "fileName").String()

			if fileName == "payload.vsix" {
				fileName = sha256 + ".vsix"
			}

			if filepath.Ext(fileName) == ".msi" {
				tmpPath = filepath.Join(flags.TmpPath, filesystem.GetFileName(fileName))
			} else if filepath.Ext(fileName) == ".vsix" {
				tmpPath = flags.TmpPath
			}

			if err := download.FileValidated(url, sha256, fileName, tmpPath); err != nil {
				charmbracelet.SharedLogger.Errorf("Error downloading MSVC package: %v", err)
				continue
			}

			fpath := filepath.Join(tmpPath, fileName)

			switch filepath.Ext(fpath) {
			case ".vsix":
				charmbracelet.SharedLogger.Infof("Extracting: %s", fpath)
				if err := internal.ExtractVsix(fpath, flags.Dest); err != nil {
					return err
				}
				break outer
			case ".msi":
				cabinetCount = 1
				storedMsi = fpath
			case ".cab":
				cabinetCount++
			default:
				charmbracelet.SharedLogger.Warnf("Unknown file format: %s, %s", filepath.Ext(fpath), fpath)
			}

			if cabinetCount >= len(packages) {
				charmbracelet.SharedLogger.Infof("Extracting: %s", fpath)

				if err := internal.ExtractMsi(flags, storedMsi, flags.Dest); err != nil {
					return err
				}

				break
			}
		}
	}

	if err := moveProgramData(flags); err != nil {
		return err
	}

	if err := moveProgramFiles(flags); err != nil {
		return err
	}

	return nil
}

func moveProgramData(flags *aflag.Flags) error {
	msiProgramData := filepath.Join(flags.Dest, "ProgramData")
	src := filepath.Join(msiProgramData, "Microsoft", "VisualStudio")
	dest := filepath.Join(flags.Dest, "VisualStudio")

	if filesystem.Exists(src) {
		if err := acopy.Copy(src, dest); err != nil {
			return err
		}
	}

	if err := filesystem.DeleteDirectory(msiProgramData); err != nil {
		return err
	}

	return nil
}

func moveProgramFiles(flags *aflag.Flags) error {
	msiProgramFiles := filepath.Join(flags.Dest, "Program Files")

	if !filesystem.Exists(msiProgramFiles) {
		return nil
	}

	dirs, err := os.ReadDir(msiProgramFiles)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if dir.Name() == "Microsoft Visual Studio 14.0" {
			src := filepath.Join(msiProgramFiles, dir.Name())

			subdirs, err := os.ReadDir(src)
			if err != nil {
				return err
			}

			for _, subdir := range subdirs {
				if err := acopy.Copy(filepath.Join(src, subdir.Name()), filepath.Join(flags.Dest, subdir.Name())); err != nil {
					return err
				}
			}
		} else {
			if err := acopy.Copy(filepath.Join(msiProgramFiles, dir.Name()), filepath.Join(flags.Dest, dir.Name())); err != nil {
				return err
			}
		}
	}

	if err := filesystem.DeleteDirectory(msiProgramFiles); err != nil {
		return err
	}

	return nil
}
