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
	"os"
	"path/filepath"

	aflag "github.com/ricochhet/portablebuildtools/flag"
)

func CreateDirectories(flags *aflag.Flags) (string, error) {
	directories := []string{flags.Downloads, flags.DownloadsCRTD, flags.DownloadsDIA, flags.Output}
	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return "", err
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return wd, nil
}

func WriteVars(flags *aflag.Flags) error {
	msvcv, err := GetMSVCVersion(flags)
	if err != nil {
		return err
	}

	sdkv, err := GetWinSDKVersion(flags)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(flags.Output, "set_vars64.bat"), x64(msvcv, sdkv, flags.Targetx64, flags.Targetx86, flags), 0o600)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(flags.Output, "set_vars32.bat"), x86(msvcv, sdkv, flags.Targetx86, flags.Targetx64, flags), 0o600)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(flags.Output, "set_vars_arm64.bat"), x64(msvcv, sdkv, flags.Targetarm64, flags.Targetarm, flags), 0o600)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(flags.Output, "set_vars_arm32.bat"), x86(msvcv, sdkv, flags.Targetarm, flags.Targetarm64, flags), 0o600)
	if err != nil {
		return err
	}

	return nil
}

func x64(msvcv, sdkv, targetA, targetB string, flags *aflag.Flags) []byte {
	return []byte(aflag.NewMSVCX64Vars(msvcv, sdkv, targetA, targetB, flags.Host, flags))
}

func x86(msvcv, sdkv, targetA, targetB string, flags *aflag.Flags) []byte {
	return []byte(aflag.NewMSVCX86Vars(msvcv, sdkv, targetA, targetB, flags.Host, flags))
}
