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

	aflag "github.com/ricochhet/portablebuildtools/flag"
)

var errMSIExtractMissing = errors.New("MSIExtract tool was not found")

func FindMsiExtract() (string, error) {
	var executable string
	if runtime.GOOS == "windows" {
		executable = "msiextract.exe"
	} else {
		executable = "msiextract"
	}

	if exists, err := aflag.IsFile(filepath.Join("./", executable)); err == nil && exists {
		return filepath.Join("./", executable), nil
	}

	if lookPath, err := exec.LookPath(executable); err == nil {
		return lookPath, nil
	}

	return "", errMSIExtractMissing
}

func ExtractMsi(flags *aflag.Flags, args ...string) error {
	path, err := FindMsiExtract()
	if err != nil {
		return err
	}

	if flags.Verbose {
		args = append(args, "-s")
		return Exec(path, args...)
	}

	return Exec(path, args...)
}
