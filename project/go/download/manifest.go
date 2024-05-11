package download

import (
	"fmt"

	aflag "github.com/ricochhet/sdkstandalone/flag"
	"github.com/tidwall/gjson"
)

func Getmanifest(f *aflag.Flags) (string, error) {
	manifest := ""
	if b, err := Download(f.MANIFEST_URL); err == nil {
		manifest = string(b)
	} else {
		fmt.Println("Error downloading main manifest:", err)
		return "", err
	}

	channelItems := gjson.Get(manifest, "channelItems").Array()
	vschannelmanifest := ""
	for _, item := range channelItems {
		if gjson.Get(item.String(), "id").Str == "Microsoft.VisualStudio.Manifests.VisualStudio" {
			vschannelmanifest = item.String()
			break
		}
	}

	vsmanifestjson := ""
	payload := gjson.Get(vschannelmanifest, "payloads").Array()[0].String()
	if b, err := Download(gjson.Get(payload, "url").String()); err == nil {
		vsmanifestjson = string(b)
	} else {
		return "", err
	}

	return vsmanifestjson, nil
}
