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
	"path/filepath"
	"strings"

	"github.com/ricochhet/minicommon/download"
	"github.com/ricochhet/minicommon/logger"
	"github.com/ricochhet/minicommon/zip"
	aflag "github.com/ricochhet/portablebuildtools/flag"
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

			if err := download.FileValidated(url, sha256, fileName, flags.TmpPath); err != nil {
				logger.SharedLogger.Errorf("Error downloading MSVC package: %v", err)
				continue
			}

			fpath := filepath.Join(flags.TmpPath, fileName)

			logger.SharedLogger.Infof("Extracting: %s", fpath)

			if err := extractVsix(fpath, flags.Dest); err != nil {
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

	return zip.UnzipByPrefixWithMessenger(fpath, destpath, "Contents", zip.DefaultUnzipMessenger())
}
