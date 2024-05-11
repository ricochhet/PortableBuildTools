package download

import (
	"fmt"

	aflag "github.com/ricochhet/sdkstandalone/flag"
	"github.com/tidwall/gjson"
)

func Getmanifest(f *aflag.Flags) (string, error) {
	var manifest string
	if dl, err := Download(f.MANIFEST_URL); err == nil {
		manifest = string(dl)
	} else {
		fmt.Println("Error downloading main manifest:", err)
		return "", err
	}

	channel := gjson.Get(manifest, "channelItems").Array()
	var vs string
	for _, item := range channel {
		if gjson.Get(item.String(), "id").Str == "Microsoft.VisualStudio.Manifests.VisualStudio" {
			vs = item.String()
			break
		}
	}

	var vsmanifest string
	payload := gjson.Get(vs, "payloads").Array()[0].String()
	if dl, err := Download(gjson.Get(payload, "url").String()); err == nil {
		vsmanifest = string(dl)
	} else {
		return "", err
	}

	return vsmanifest, nil
}
