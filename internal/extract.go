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
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ricochhet/minicommon/zip"
	aflag "github.com/ricochhet/portablebuildtools/flag"
)

var (
	errMSIExtractMissing = errors.New("MSIExtract tool was not found")
	errNotVsixFile       = errors.New("file is not a vsix file")
)

func FindMsiExtract() (string, bool, error) {
	var executable string
	if runtime.GOOS == "windows" {
		executable = "msiextract.exe"
	} else {
		executable = "msiextract"
	}

	if exists, err := aflag.IsFile(filepath.Join("./", executable)); err == nil && exists {
		return filepath.Join("./", executable), true, nil
	}

	if lookPath, err := exec.LookPath(executable); err == nil {
		return lookPath, false, nil
	}

	return "", false, errMSIExtractMissing
}

func ExtractMsi(flags *aflag.Flags, args ...string) error {
	path, rel, err := FindMsiExtract()
	if err != nil {
		return err
	}

	if flags.Verbose {
		args = append(args, "-s")
		return Exec(path, rel, args...)
	}

	return Exec(path, rel, args...)
}

func ExtractVsix(fpath, destpath string) error {
	if !strings.HasSuffix(fpath, ".vsix") {
		return errNotVsixFile
	}

	return zip.UnzipByPrefixWithMessenger(fpath, destpath, "Contents", zip.DefaultUnzipMessenger())
}
