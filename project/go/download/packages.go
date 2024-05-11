package download

import (
	"slices"
	"strings"

	aflag "github.com/ricochhet/sdkstandalone/flag"
	"github.com/tidwall/gjson"
)

func Getpackages(f *aflag.Flags, manifest string, msvcPackages []string) ([]string, []string, []string, []gjson.Result) {
	packages := gjson.Get(manifest, "packages").Array()
	var payloads []string
	var crtd []string
	var dia []string
	var sdk []gjson.Result
	for _, item := range packages {
		pid := strings.ToLower(gjson.Get(item.String(), "id").String())
		if slices.Contains(msvcPackages, pid) {
			t := gjson.Get(item.String(), "type").String()
			l := gjson.Get(item.String(), "language").String()
			if t == "Vsix" && (l == "en-US" || l == "") {
				payloads = append(payloads, gjson.Get(item.String(), "payloads").String())
			}
		} else if pid == strings.ToLower("Win11SDK_10.0.22621") {
			sdk = gjson.Get(item.String(), "payloads").Array()
		} else if pid == "microsoft.visualcpp.runtimedebug.14" && gjson.Get(item.String(), "chip").String() == f.HOST {
			crtd = append(crtd, gjson.Get(item.String(), "payloads").String())
		} else if pid == "microsoft.visualc.140.dia.sdk.msi" {
			dia = append(dia, gjson.Get(item.String(), "payloads").String())
		}
	}

	return payloads, crtd, dia, sdk
}
