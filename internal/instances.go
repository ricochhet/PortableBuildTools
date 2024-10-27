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

	acopy "github.com/otiai10/copy"
	"github.com/ricochhet/minicommon/filesystem"
	aflag "github.com/ricochhet/portablebuildtools/flag"
)

var errNoInstances = errors.New("_Instances does not exist")

//nolint:lll // wontfix
func CopyInstances(flags *aflag.Flags) error {
	if !filesystem.Exists(filesystem.GetRelativePath("_Instances")) {
		return errNoInstances
	}

	if err := acopy.Copy(filesystem.GetRelativePath("_Instances"), filepath.Join(flags.Dest, "VisualStudio", "Packages", "_Instances")); err != nil {
		return err
	}

	return nil
}
