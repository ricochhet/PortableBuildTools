package download

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	aflag "github.com/ricochhet/portablebuildtools/flag"
	acopy "github.com/ricochhet/portablebuildtools/thirdparty/copy"
	"github.com/tidwall/gjson"
)

func GetCRTD(payloads []string, destx64, destx86, destarm, destarm64 string, flags *aflag.Flags) error {
	for _, payload := range payloads {
		packages := gjson.Parse(payload).Array()
		for _, pkg := range packages {
			url := pkg.Get("url").String()
			sha256 := pkg.Get("sha256").String()
			fileName := pkg.Get("fileName").String()

			if _, err := File(url, sha256, fileName, flags.Downloads, fileName); err != nil {
				fmt.Println("Error downloading CRTD package:", err)
				continue
			}
		}
	}

	if err := extractMSI(flags, filepath.Join(flags.Downloads, "vc_RuntimeDebug.msi"), flags.Downloads); err != nil {
		return err
	}

	dlls, err := os.ReadDir(filepath.Join(flags.Downloads, "System64"))
	if err != nil {
		return err
	}

	for _, dll := range dlls {
		paths := []copyCRTDPath{
			{dest: destx64, flags: flags},
			{dest: destx86, flags: flags},
		}

		if flags.DownloadARMTargets {
			paths = append(paths,
				copyCRTDPath{dest: destarm, flags: flags},
				copyCRTDPath{dest: destarm64, flags: flags},
			)
		}

		if err := copyCRTDToPaths(dll, paths); err != nil {
			return err
		}
	}

	return nil
}

type copyCRTDPath struct {
	dest  string
	flags *aflag.Flags
}

func copyCRTDToPaths(dll fs.DirEntry, paths []copyCRTDPath) error {
	for _, path := range paths {
		if err := copyCRTD(dll, path.dest, path.flags); err != nil {
			return err
		}
	}

	return nil
}

func copyCRTD(dirEntry fs.DirEntry, target string, flags *aflag.Flags) error {
	return acopy.Copy(filepath.Join(filepath.Join(flags.Downloads, "System64"), dirEntry.Name()), filepath.Join(target, dirEntry.Name()))
}
