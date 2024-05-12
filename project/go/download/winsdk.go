package download

import (
	"bytes"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	aflag "github.com/ricochhet/portablebuildtools/flag"
	"github.com/tidwall/gjson"
)

func GetWinSDK(flags *aflag.Flags, packages []gjson.Result, winsdkpackages []string) error {
	installerPrefix := "Installers\\"
	installers := []string{}
	cabinets := []string{}

	for _, pkg := range packages {
		name := pkg.Get("fileName").String()
		url := pkg.Get("url").String()
		sha256 := pkg.Get("sha256").String()

		if slices.Contains(winsdkpackages, name) {
			fileName := strings.TrimPrefix(name, installerPrefix)
			installer, err := File(url, sha256, fileName, flags.Downloads, fileName)
			if err != nil { //nolint:wsl // gofumpt conflict
				fmt.Println("Error downloading Windows SDK package:", err)
				continue
			}

			installers = append(installers, filepath.Join(flags.Downloads, fileName))
			cabinets = append(cabinets, getMSICabinets(installer)...)
		}
	}

	for _, pkg := range packages {
		name := pkg.Get("fileName").String()
		url := pkg.Get("url").String()
		sha256 := pkg.Get("sha256").String()

		if slices.Contains(cabinets, strings.TrimPrefix(name, installerPrefix)) {
			fileName := strings.TrimPrefix(name, installerPrefix)

			if _, err := File(url, sha256, fileName, flags.Downloads, fileName); err != nil {
				fmt.Println("Error downloading cab:", err)
				continue
			}
		}
	}

	for _, installer := range installers {
		err := extractMSI(flags, installer, flags.Output)
		if err != nil {
			return err
		}
	}

	return nil
}

//nolint:mnd // represent *.cab location in msi
func getMSICabinets(installerMsi []byte) []string {
	cabs := []string{}
	index := 0

	for {
		foundIndex := bytes.Index(installerMsi[index:], []byte(".cab"))
		if foundIndex < 0 {
			break
		}

		index += foundIndex + 4
		start := index - 36

		if start < 0 {
			start = 0
		}

		cabs = append(cabs, string(installerMsi[start:index]))
	}

	return cabs
}
