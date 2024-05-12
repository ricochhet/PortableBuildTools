package download

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	aflag "github.com/ricochhet/portablebuildtools/flag"
	"github.com/ricochhet/simpledownload"
	"github.com/ricochhet/simplezip"
	"github.com/tidwall/gjson"
)

var errNotVsixFile = errors.New("file is not a vsix file")

func GetPayloads(flags *aflag.Flags, payloads []string) error {
	for _, payload := range payloads {
		packages := gjson.Parse(payload).Array()
		for _, pkg := range packages {
			url := gjson.Get(pkg.String(), "url").String()
			sha256 := gjson.Get(pkg.String(), "sha256").String()
			fileName := gjson.Get(pkg.String(), "fileName").String()

			if err := simpledownload.File(url, sha256, fileName, flags.Downloads); err != nil {
				fmt.Println("Error downloading MSVC package:", err)
				continue
			}

			fpath := filepath.Join(flags.Downloads, fileName)

			fmt.Println("Extracting: ", fpath)

			if err := vsix(fpath, flags.Output); err != nil {
				return err
			}

			break
		}
	}

	return nil
}

func vsix(fpath, destpath string) error {
	if !strings.HasSuffix(fpath, ".vsix") {
		return errNotVsixFile
	}

	return simplezip.UnzipByPrefixWithMessenger(fpath, destpath, "Contents", simplezip.DefaultUnzipMessenger())
}
