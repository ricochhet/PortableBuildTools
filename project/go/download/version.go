package download

import (
	"errors"
	"path/filepath"

	aflag "github.com/ricochhet/portablebuildtools/flag"
)

var errNoVersionInDirectory = errors.New("no version in directory")

func GetMSVCVersion(flags *aflag.Flags) (string, error) {
	return getVersion(filepath.Join(flags.Output, "VC", "Tools", "MSVC"))
}

func GetWinSDKVersion(flags *aflag.Flags) (string, error) {
	return getVersion(filepath.Join(flags.Output, "Windows Kits", "10", "bin"))
}

func getVersion(apath string) (string, error) {
	versions, err := filepath.Glob(filepath.Join(apath, "*"))
	if err != nil {
		return "", err
	}

	if len(versions) == 0 {
		return "", errNoVersionInDirectory
	}

	return filepath.Base(versions[0]), nil
}
