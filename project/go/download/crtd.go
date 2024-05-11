package download

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	aflag "github.com/ricochhet/sdkstandalone/flag"
	acopy "github.com/ricochhet/sdkstandalone/thirdparty/copy"
	"github.com/tidwall/gjson"
)

func Getcrtd(payloads []string, dstX64, dstX86, dstARM, dstARM64 string, f *aflag.Flags) error {
	for _, payload := range payloads {
		packages := gjson.Parse(payload).Array()
		for _, pkg := range packages {
			url := gjson.Get(pkg.String(), "url").String()
			sha256 := gjson.Get(pkg.String(), "sha256").String()
			fileName := gjson.Get(pkg.String(), "fileName").String()
			if _, err := Downloadprogress(url, sha256, fileName, f.DOWNLOADS, fileName); err != nil {
				fmt.Println("Error downloading CRTD package:", err)
				continue
			}
		}
	}
	err := Rustmsiexec(f, "./rust-msiexec.exe", filepath.Join(f.DOWNLOADS, "vc_RuntimeDebug.msi"), f.DOWNLOADS)
	if err != nil {
		return err
	}
	dlls, err := os.ReadDir(filepath.Join(f.DOWNLOADS, "System64"))
	if err != nil {
		return err
	}
	for _, dll := range dlls {
		err := copycrtd(dll, dstX64, f)
		if err != nil {
			return err
		}

		err = copycrtd(dll, dstX86, f)
		if err != nil {
			return err
		}

		if f.DOWNLOAD_ARM_TARGETS {
			err := copycrtd(dll, dstARM, f)
			if err != nil {
				return err
			}

			err = copycrtd(dll, dstARM64, f)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copycrtd(dirEntry fs.DirEntry, target string, f *aflag.Flags) error {
	return acopy.Copy(filepath.Join(filepath.Join(f.DOWNLOADS, "System64"), dirEntry.Name()), filepath.Join(target, dirEntry.Name()))
}
