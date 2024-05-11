package download

import (
	"fmt"
	"path/filepath"

	"github.com/ricochhet/sdkstandalone/extract"
	aflag "github.com/ricochhet/sdkstandalone/flag"
	"github.com/tidwall/gjson"
)

func Getpayloads(f *aflag.Flags, payloads []string) {
	for _, item := range payloads {
		pkgs := gjson.Parse(item).Array()
		for _, pkg := range pkgs {
			url := gjson.Get(pkg.String(), "url").String()
			sha256 := gjson.Get(pkg.String(), "sha256").String()
			fileName := gjson.Get(pkg.String(), "fileName").String()
			if _, err := Downloadprogress(url, sha256, fileName, f.DOWNLOADS, fileName); err != nil {
				fmt.Println("Error downloading MSVC package:", err)
				continue
			} else {
				fpath := filepath.Join(f.DOWNLOADS, fileName)
				fmt.Println("Extracting: ", fpath)
				extract.Extractpackage(fpath, f.OUTPUT)
				break
			}
		}
	}
}
