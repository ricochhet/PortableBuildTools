package download

import (
	"errors"
	"fmt"
	"path/filepath"

	aflag "github.com/ricochhet/sdkstandalone/flag"
	"github.com/ricochhet/sdkstandalone/process"
	acopy "github.com/ricochhet/sdkstandalone/thirdparty/copy"
	"github.com/tidwall/gjson"
)

var msdia140dll = "msdia140.dll"

func Getdiasdk(payloads []string, dstX64, dstX86, dstARM, dstARM64 string, f *aflag.Flags) error {
	for _, payload := range payloads {
		packages := gjson.Parse(payload).Array()
		for _, pkg := range packages {
			url := gjson.Get(pkg.String(), "url").String()
			sha256 := gjson.Get(pkg.String(), "sha256").String()
			fileName := gjson.Get(pkg.String(), "fileName").String()
			if _, err := Downloadprogress(url, sha256, fileName, f.DOWNLOADS_DIA, fileName); err != nil {
				fmt.Println("Error downloading DIA SDK package:", err)
				continue
			}
		}
	}
	err := process.Exec("./rust-msiexec.exe", filepath.Join(f.DOWNLOADS_DIA, "VC_diasdk.msi"), f.DOWNLOADS_DIA)
	if err != nil {
		return err
	}
	msdia := ""
	if f.HOST == f.TARGETX64 {
		msdia = msdia140dll
	} else if f.HOST == f.TARGETX86 {
		msdia = "amd64/" + msdia140dll
	} else {
		return errors.New("unknown host architecture")
	}

	err = copymsdiadll(msdia, filepath.Join(dstX64, msdia140dll), f)
	if err != nil {
		return err
	}
	err = copymsdiadll(msdia, filepath.Join(dstX86, msdia140dll), f)
	if err != nil {
		return err
	}
	if f.DOWNLOAD_ARM_TARGETS {
		err = copymsdiadll(msdia, filepath.Join(dstARM, msdia140dll), f)
		if err != nil {
			return err
		}
		err = copymsdiadll(msdia, filepath.Join(dstARM64, msdia140dll), f)
		if err != nil {
			return err
		}
	}

	return err
}

func copymsdiadll(msdia, target string, f *aflag.Flags) error {
	return acopy.Copy(filepath.Join(f.DOWNLOADS_DIA, "Program Files", "Microsoft Visual Studio 14.0", "DIA SDK", "bin", msdia), target)
}
