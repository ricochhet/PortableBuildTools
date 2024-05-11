package download

import (
	"slices"
	"strings"

	aflag "github.com/ricochhet/sdkstandalone/flag"
	"github.com/tidwall/gjson"
)

func Getpackages(f *aflag.Flags, manifest string, msvcPackages []string) ([]string, []string, []string, []gjson.Result) {
	packages := gjson.Get(manifest, "packages").Array()
	payloads := []string{}
	crtdPayloads := []string{}
	diaPayloads := []string{}
	sdkPayloads := []gjson.Result{}
	for _, pkg := range packages {
		pid := strings.ToLower(gjson.Get(pkg.String(), "id").String())
		if slices.Contains(msvcPackages, pid) {
			t := gjson.Get(pkg.String(), "type").String()
			l := gjson.Get(pkg.String(), "language").String()
			if t == "Vsix" && (l == "en-US" || l == "") {
				payloads = append(payloads, gjson.Get(pkg.String(), "payloads").String())
			}
		} else if pid == strings.ToLower("Win11SDK_10.0.22621") {
			sdkPayloads = gjson.Get(pkg.String(), "payloads").Array()
		} else if pid == "microsoft.visualcpp.runtimedebug.14" && gjson.Get(pkg.String(), "chip").String() == f.HOST {
			crtdPayloads = append(crtdPayloads, gjson.Get(pkg.String(), "payloads").String())
		} else if pid == "microsoft.visualc.140.dia.sdk.msi" {
			diaPayloads = append(diaPayloads, gjson.Get(pkg.String(), "payloads").String())
		}
	}

	return payloads, crtdPayloads, diaPayloads, sdkPayloads
}
