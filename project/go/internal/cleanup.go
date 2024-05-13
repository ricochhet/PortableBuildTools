package internal

import (
	"os"
	"path/filepath"

	aflag "github.com/ricochhet/portablebuildtools/flag"
)

func RemoveVCTipsTelemetry(flags *aflag.Flags) error {
	vctipExe := "vctip.exe"
	msvcv, err := GetMSVCVersion(flags)
	if err != nil { //nolint:wsl // gofumpt conflict
		return err
	}

	paths := []string{
		filepath.Join(flags.Output, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetx64, vctipExe),
		filepath.Join(flags.Output, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetx86, vctipExe),
	}

	if flags.DownloadARMTargets {
		paths = append(paths,
			filepath.Join(flags.Output, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetarm, vctipExe),
			filepath.Join(flags.Output, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetarm64, vctipExe),
		)
	}

	for _, path := range paths {
		err := os.Remove(path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

func CleanHostDirectory(flags *aflag.Flags) error {
	targets := []string{flags.Targetx64, flags.Targetx86, flags.Targetarm, flags.Targetarm64}
	msvcv, err := GetMSVCVersion(flags)
	if err != nil { //nolint:wsl // gofumpt conflict
		return err
	}

	for _, arch := range targets {
		if arch != flags.Host {
			err := os.RemoveAll(filepath.Join(flags.Output, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+arch))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
