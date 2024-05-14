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
	"errors"
	"fmt"
	"path/filepath"

	aflag "github.com/ricochhet/portablebuildtools/flag"
	"github.com/ricochhet/portablebuildtools/internal"
	acopy "github.com/ricochhet/portablebuildtools/thirdparty/copy"
	"github.com/ricochhet/simpledownload"
	"github.com/tidwall/gjson"
)

var errUnknownHostArch = errors.New("unknown host architecture")

func GetDIASDK(payloads []string, destx64, destx86, destarm, destarm64 string, flags *aflag.Flags) error {
	msdia140dll := "msdia140.dll"

	for _, payload := range payloads {
		packages := gjson.Parse(payload).Array()
		for _, pkg := range packages {
			url := pkg.Get("url").String()
			sha256 := pkg.Get("sha256").String()
			fileName := pkg.Get("fileName").String()

			if err := simpledownload.FileValidated(url, sha256, fileName, flags.DownloadsDIA); err != nil {
				fmt.Println("Error downloading DIA SDK package:", err)
				continue
			}
		}
	}

	if err := internal.ExtractMSI(flags, filepath.Join(flags.DownloadsDIA, "VC_diasdk.msi"), flags.DownloadsDIA); err != nil {
		return err
	}

	var msdia string

	switch flags.Host {
	case flags.Targetx64:
		msdia = msdia140dll
	case flags.Targetx86:
		msdia = "amd64/" + msdia140dll
	default:
		return errUnknownHostArch
	}

	paths := []copyDIAPath{
		{dest: filepath.Join(destx64, msdia), flags: flags},
		{dest: filepath.Join(destx86, msdia), flags: flags},
	}

	if flags.DownloadARMTargets {
		paths = append(paths,
			copyDIAPath{dest: filepath.Join(destarm, msdia), flags: flags},
			copyDIAPath{dest: filepath.Join(destarm64, msdia), flags: flags},
		)
	}

	if err := copyMSDIADllToPaths(msdia, paths); err != nil {
		return err
	}

	return nil
}

type copyDIAPath struct {
	dest  string
	flags *aflag.Flags
}

func copyMSDIADllToPaths(msdia string, paths []copyDIAPath) error {
	for _, path := range paths {
		if err := copyMSDIADLL(msdia, path.dest, path.flags); err != nil {
			return err
		}
	}

	return nil
}

func copyMSDIADLL(msdia, target string, flags *aflag.Flags) error {
	return acopy.Copy(filepath.Join(flags.DownloadsDIA, "Program Files", "Microsoft Visual Studio 14.0", "DIA SDK", "bin", msdia), target)
}
