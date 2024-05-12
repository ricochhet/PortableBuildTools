package download

import (
	"slices"
	"strings"

	aflag "github.com/ricochhet/sdkstandalone/flag"
	"github.com/tidwall/gjson"
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
			case strings.ToLower("Win11SDK_10.0.22621"):
				sdkPayloads = gjson.Get(pkg.String(), "payloads").Array()
			case "microsoft.visualcpp.runtimedebug.14":
				if gjson.Get(pkg.String(), "chip").String() == flags.Host {
					crtdPayloads = append(crtdPayloads, gjson.Get(pkg.String(), "payloads").String())
				}
			case "microsoft.visualc.140.dia.sdk.msi":
				diaPayloads = append(diaPayloads, gjson.Get(pkg.String(), "payloads").String())
			}
		}
	}

	return payloads, crtdPayloads, diaPayloads, sdkPayloads
}
