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
	"slices"
	"strings"

	aflag "github.com/ricochhet/portablebuildtools/flag"
	"github.com/tidwall/gjson"
)

var (
	sdkPid    = "Win11SDK_10.0.22621"                 //nolint:gochecknoglobals // ...
	rtDbgPid  = "microsoft.visualcpp.runtimedebug.14" //nolint:gochecknoglobals // ...
	diaSdkPid = "microsoft.visualc.140.dia.sdk.msi"   //nolint:gochecknoglobals // ...
)

//nolint:cyclop // wontfix
func GetPackages(flags *aflag.Flags, manifest string, msvcpackages []string) ([]string, []string, []string, []gjson.Result) {
	packages := gjson.Get(manifest, "packages").Array()

	var (
		payloads     []string
		crtdPayloads []string
		diaPayloads  []string
		sdkPayloads  []gjson.Result
	)

	for _, pkg := range packages {
		pid := strings.ToLower(gjson.Get(pkg.String(), "id").String())
		if slices.Contains(msvcpackages, pid) {
			fileType := gjson.Get(pkg.String(), "type").String()
			language := gjson.Get(pkg.String(), "language").String()

			if (fileType == "Vsix" || fileType == "Msi") && (language == "en-US" || language == "" || language == "neutral") {
				payloads = append(payloads, gjson.Get(pkg.String(), "payloads").String())
			}
		} else {
			switch pid {
			case strings.ToLower(sdkPid):
				sdkPayloads = gjson.Get(pkg.String(), "payloads").Array()
			case rtDbgPid:
				chip := gjson.Get(pkg.String(), "chip").String()
				if chip == flags.Host || chip == "neutral" {
					crtdPayloads = append(crtdPayloads, gjson.Get(pkg.String(), "payloads").String())
				}
			case diaSdkPid:
				diaPayloads = append(diaPayloads, gjson.Get(pkg.String(), "payloads").String())
			}
		}
	}

	return payloads, crtdPayloads, diaPayloads, sdkPayloads
}
