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
	"strings"

	aflag "github.com/ricochhet/portablebuildtools/flag"
	"github.com/ricochhet/simpledownload"
	"github.com/ricochhet/simplezip"
	"github.com/tidwall/gjson"
)

var errNotVsixFile = errors.New("file is not a vsix file")

func GetPayloads(flags *aflag.Flags, payloads []string) error {
	for _, payload := range payloads {
		packages := gjson.Parse(payload).Array()
		for _, pkg := range packages {
			url := gjson.Get(pkg.String(), "url").String()
			sha256 := gjson.Get(pkg.String(), "sha256").String()
			fileName := gjson.Get(pkg.String(), "fileName").String()

			if fileName == "payload.vsix" {
				fileName = sha256 + ".vsix"
			}

			if err := simpledownload.FileValidated(url, sha256, fileName, flags.Downloads); err != nil {
				fmt.Println("Error downloading MSVC package:", err)
				continue
			}

			fpath := filepath.Join(flags.Downloads, fileName)

			fmt.Println("Extracting: ", fpath)

			if err := extractVsix(fpath, flags.Output); err != nil {
				return err
			}

			break
		}
	}

	return nil
}

func extractVsix(fpath, destpath string) error {
	if !strings.HasSuffix(fpath, ".vsix") {
		return errNotVsixFile
	}

	return simplezip.UnzipByPrefixWithMessenger(fpath, destpath, "Contents", simplezip.DefaultUnzipMessenger())
}
