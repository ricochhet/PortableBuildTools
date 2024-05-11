package download

import (
	"fmt"
	"path/filepath"

	"github.com/ricochhet/sdkstandalone/extract"
	aflag "github.com/ricochhet/sdkstandalone/flag"
	"github.com/tidwall/gjson"
)

func Getpayloads(f *aflag.Flags, payloads []string) {
	for _, payload := range payloads {
		packages := gjson.Parse(payload).Array()
		for _, pkg := range packages {
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
