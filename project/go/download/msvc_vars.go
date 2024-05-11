package download

import (
	"os"
	"path/filepath"

	aflag "github.com/ricochhet/sdkstandalone/flag"
)

func Writevars(f *aflag.Flags) error {
	sdkv, err := Getwinsdkversion(f)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(f.OUTPUT, "set_vars64.bat"), []byte(aflag.NewMsvcX64Vars(f.MSVC_VERSION_LOCAL, sdkv, f.TARGETX64, f.TARGETX86, f.HOST)), 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(f.OUTPUT, "set_vars32.bat"), []byte(aflag.NewMsvcX86Vars(f.MSVC_VERSION_LOCAL, sdkv, f.TARGETX86, f.TARGETX64, f.HOST)), 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(f.OUTPUT, "set_vars_arm64.bat"), []byte(aflag.NewMsvcX64Vars(f.MSVC_VERSION_LOCAL, sdkv, f.TARGETARM64, f.TARGETARM, f.HOST)), 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(f.OUTPUT, "set_vars_arm32.bat"), []byte(aflag.NewMsvcX86Vars(f.MSVC_VERSION_LOCAL, sdkv, f.TARGETARM, f.TARGETARM64, f.HOST)), 0644)
	if err != nil {
		return err
	}

	return nil
}
