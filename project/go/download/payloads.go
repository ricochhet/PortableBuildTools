package download

import (
	"fmt"
	"path/filepath"

	"github.com/ricochhet/sdkstandalone/extract"
	aflag "github.com/ricochhet/sdkstandalone/flag"
	"github.com/tidwall/gjson"
)

func GetPayloads(flags *aflag.Flags, payloads []string) error {
	for _, payload := range payloads {
		packages := gjson.Parse(payload).Array()
		for _, pkg := range packages {
			url := gjson.Get(pkg.String(), "url").String()
			sha256 := gjson.Get(pkg.String(), "sha256").String()
			fileName := gjson.Get(pkg.String(), "fileName").String()

			if _, err := File(url, sha256, fileName, flags.Downloads, fileName); err != nil {
				fmt.Println("Error downloading MSVC package:", err)
				continue
			}

			fpath := filepath.Join(flags.Downloads, fileName)

			fmt.Println("Extracting: ", fpath)

			if err := extract.Vsix(fpath, flags.Output); err != nil {
				return err
			}

			break
		}
	}

	return nil
}
