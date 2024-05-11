package download

import (
	"fmt"
	"path/filepath"

	aflag "github.com/ricochhet/sdkstandalone/flag"
	"github.com/ricochhet/sdkstandalone/process"
	acopy "github.com/ricochhet/sdkstandalone/thirdparty/copy"
	"github.com/tidwall/gjson"
)

func Getdiasdk(dia []string, dstX64, dstX86, dstARM, dstARM64 string, f *aflag.Flags) error {
	for _, item := range dia {
		pkgs := gjson.Parse(item).Array()
		for _, pkg := range pkgs {
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
		panic(err)
	}
	var msdia string
	if f.HOST == f.TARGETX64 {
		msdia = "msdia140.dll"
	} else if f.HOST == f.TARGETX86 {
		msdia = "amd64/msdia140.dll"
	} else {
		return fmt.Errorf("unknown host architecture")
	}

	diaTargetX64 := filepath.Join(dstX64, "msdia140.dll")
	diaTargetX86 := filepath.Join(dstX86, "msdia140.dll")
	diaTargetARM := filepath.Join(dstARM, "msdia140.dll")
	diaTargetARM64 := filepath.Join(dstARM64, "msdia140.dll")

	err = acopy.Copy(filepath.Join(f.DOWNLOADS_DIA, "Program Files", "Microsoft Visual Studio 14.0", "DIA SDK", "bin", msdia), diaTargetX64)
	if err != nil {
		return err
	}
	err = acopy.Copy(filepath.Join(f.DOWNLOADS_DIA, "Program Files", "Microsoft Visual Studio 14.0", "DIA SDK", "bin", msdia), diaTargetX86)
	if err != nil {
		return err
	}
	if f.DOWNLOAD_ARM_TARGETS {
		err = acopy.Copy(filepath.Join(f.DOWNLOADS_DIA, "Program Files", "Microsoft Visual Studio 14.0", "DIA SDK", "bin", msdia), diaTargetARM)
		if err != nil {
			return err
		}
		err = acopy.Copy(filepath.Join(f.DOWNLOADS_DIA, "Program Files", "Microsoft Visual Studio 14.0", "DIA SDK", "bin", msdia), diaTargetARM64)
		if err != nil {
			return err
		}
	}

	return err
}
