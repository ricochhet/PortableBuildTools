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
	"github.com/ricochhet/minicommon/charmbracelet"
	"github.com/ricochhet/minicommon/download"
	aflag "github.com/ricochhet/portablebuildtools/flag"
	"github.com/tidwall/gjson"
)

func GetManifest(flags *aflag.Flags) (string, error) {
	var manifest string
	if b, err := download.Download(flags.ManifestURL); err == nil {
		manifest = string(b)
	} else {
		charmbracelet.SharedLogger.Errorf("Error downloading main manifest: %v", err)
		return "", err
	}

	channelItems := gjson.Get(manifest, "channelItems").Array()
	vsChannelManifest := ""

	for _, item := range channelItems {
		if gjson.Get(item.String(), "id").Str == "Microsoft.VisualStudio.Manifests.VisualStudio" {
			vsChannelManifest = item.String()
			break
		}
	}

	var vsManifestJSON string

	payload := gjson.Get(vsChannelManifest, "payloads").Array()[0].String()
	if b, err := download.Download(gjson.Get(payload, "url").String()); err == nil {
		vsManifestJSON = string(b)
	} else {
		return "", err
	}

	return vsManifestJSON, nil
}
