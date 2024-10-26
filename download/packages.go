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
	sdkPid    = "Win11SDK_10.0.22621"
	rtDbgPid  = "microsoft.visualcpp.runtimedebug.14"
	diaSdkPid = "microsoft.visualc.140.dia.sdk.msi"
)

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

			if fileType == "Vsix" && (language == "en-US" || language == "") {
				payloads = append(payloads, gjson.Get(pkg.String(), "payloads").String())
			}
		} else {
			switch pid {
			case strings.ToLower(sdkPid):
				sdkPayloads = gjson.Get(pkg.String(), "payloads").Array()
			case rtDbgPid:
				if gjson.Get(pkg.String(), "chip").String() == flags.Host {
					crtdPayloads = append(crtdPayloads, gjson.Get(pkg.String(), "payloads").String())
				}
			case diaSdkPid:
				diaPayloads = append(diaPayloads, gjson.Get(pkg.String(), "payloads").String())
			}
		}
	}

	return payloads, crtdPayloads, diaPayloads, sdkPayloads
}
