package download

import (
	"fmt"
	"os"
	"path/filepath"

	aflag "github.com/ricochhet/sdkstandalone/flag"
	"github.com/ricochhet/sdkstandalone/process"
	acopy "github.com/ricochhet/sdkstandalone/thirdparty/copy"
	"github.com/tidwall/gjson"
)

func Getcrtd(crtd []string, dstX64, dstX86, dstARM, dstARM64 string, f *aflag.Flags) error {
	for _, item := range crtd {
		pkgs := gjson.Parse(item).Array()
		for _, pkg := range pkgs {
			url := gjson.Get(pkg.String(), "url").String()
			sha256 := gjson.Get(pkg.String(), "sha256").String()
			fileName := gjson.Get(pkg.String(), "fileName").String()
			if _, err := Downloadprogress(url, sha256, fileName, f.DOWNLOADS, fileName); err != nil {
				fmt.Println("Error downloading CRTD package:", err)
				continue
			}
		}
	}
	err := process.Exec("./rust-msiexec.exe", filepath.Join(f.DOWNLOADS, "vc_RuntimeDebug.msi"), f.DOWNLOADS)
	if err != nil {
		return err
	}
	crtdGlob, err := os.ReadDir(filepath.Join(f.DOWNLOADS, "System64"))
	if err != nil {
		return err
	}
	for _, item := range crtdGlob {
		err := acopy.Copy(filepath.Join(filepath.Join(f.DOWNLOADS, "System64"), item.Name()), filepath.Join(dstX64, item.Name()))
		if err != nil {
			return err
		}

		err = acopy.Copy(filepath.Join(filepath.Join(f.DOWNLOADS, "System64"), item.Name()), filepath.Join(dstX86, item.Name()))
		if err != nil {
			return err
		}

		if f.DOWNLOAD_ARM_TARGETS {
			err := acopy.Copy(filepath.Join(filepath.Join(f.DOWNLOADS, "System64"), item.Name()), filepath.Join(dstARM, item.Name()))
			if err != nil {
				return err
			}

			err = acopy.Copy(filepath.Join(filepath.Join(f.DOWNLOADS, "System64"), item.Name()), filepath.Join(dstARM64, item.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
