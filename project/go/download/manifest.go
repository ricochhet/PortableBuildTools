package download

import (
	"fmt"

	aflag "github.com/ricochhet/portablebuildtools/flag"
	"github.com/tidwall/gjson"
)

func GetManifest(flags *aflag.Flags) (string, error) {
	var manifest string
	if b, err := Download(flags.ManifestURL); err == nil {
		manifest = string(b)
	} else {
		fmt.Println("Error downloading main manifest:", err)
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
	if b, err := Download(gjson.Get(payload, "url").String()); err == nil {
		vsManifestJSON = string(b)
	} else {
		return "", err
	}

	return vsManifestJSON, nil
}
