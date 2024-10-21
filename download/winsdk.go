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
	"bytes"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ricochhet/minicommon/download"
	aflag "github.com/ricochhet/portablebuildtools/flag"
	"github.com/ricochhet/portablebuildtools/internal"
	"github.com/tidwall/gjson"
)

func GetWinSdk(flags *aflag.Flags, packages []gjson.Result, winsdkpackages []string) error {
	installerPrefix := "Installers\\"
	installers := []string{}
	cabinets := []string{}

	for _, pkg := range packages {
		name := pkg.Get("fileName").String()
		url := pkg.Get("url").String()
		sha256 := pkg.Get("sha256").String()

		if slices.Contains(winsdkpackages, name) {
			fileName := strings.TrimPrefix(name, installerPrefix)
			installer, err := download.FileWithBytesValidated(url, sha256, fileName, flags.TmpPath)
			if err != nil { //nolint:wsl // gofumpt conflict
				fmt.Println("Error downloading Windows SDK package:", err)
				continue
			}

			installers = append(installers, filepath.Join(flags.TmpPath, fileName))
			cabinets = append(cabinets, getMSICabinets(installer)...)
		}
	}

	for _, pkg := range packages {
		name := pkg.Get("fileName").String()
		url := pkg.Get("url").String()
		sha256 := pkg.Get("sha256").String()

		if slices.Contains(cabinets, strings.TrimPrefix(name, installerPrefix)) {
			fileName := strings.TrimPrefix(name, installerPrefix)

			if err := download.FileValidated(url, sha256, fileName, flags.TmpPath); err != nil {
				fmt.Println("Error downloading cab:", err)
				continue
			}
		}
	}

	for _, installer := range installers {
		err := internal.ExtractMsi(flags, installer, flags.Dest)
		if err != nil {
			return err
		}
	}

	return nil
}

//nolint:mnd // represent *.cab location in msi
func getMSICabinets(installerMsi []byte) []string {
	cabs := []string{}
	index := 0

	for {
		foundIndex := bytes.Index(installerMsi[index:], []byte(".cab"))
		if foundIndex < 0 {
			break
		}

		index += foundIndex + 4
		start := index - 36

		if start < 0 {
			start = 0
		}

		cabs = append(cabs, string(installerMsi[start:index]))
	}

	return cabs
}
