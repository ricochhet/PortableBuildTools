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
	"path/filepath"

	aflag "github.com/ricochhet/portablebuildtools/flag"
)

var errNoVersionInDirectory = errors.New("no version in directory")

func GetMsvcVersion(flags *aflag.Flags) (string, error) {
	return getVersion(filepath.Join(flags.Dest, "VC", "Tools", "MSVC"))
}

func GetWinSdkVersion(flags *aflag.Flags) (string, error) {
	return getVersion(filepath.Join(flags.Dest, "Windows Kits", "10", "bin"))
}

func getVersion(apath string) (string, error) {
	versions, err := filepath.Glob(filepath.Join(apath, "*"))
	if err != nil {
		return "", err
	}

	if len(versions) == 0 {
		return "", errNoVersionInDirectory
	}

	return filepath.Base(versions[0]), nil
}
