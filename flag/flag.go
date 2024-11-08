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

package flag

type Flags struct {
	Version           bool
	MsvcVer           string
	WinSdkVer         string
	Dest              string // output folder
	DestZip           string
	TmpPath           string // temporary download files
	TmpCrtd           string
	TmpDia            string
	Cwd               bool
	Host              string // or x86
	SetMsvcPackages   string
	SetWinSdkPackages string
	SpectreLibs       bool
	ArmTargets        bool
	LlvmClang         bool
	UnitTest          bool
	Cmake             bool
	MfcAtl            bool
	Vcpkg             bool
	Msbuild           bool
	ManifestURL       string
	Targetx64         string
	Targetx86         string
	Targetarm         string
	Targetarm64       string
	WriteEnvironment  bool
	Verbose           bool
	Zip               bool
}
