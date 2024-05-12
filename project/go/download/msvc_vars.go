package download

import (
	"os"
	"path/filepath"

	aflag "github.com/ricochhet/sdkstandalone/flag"
)

func WriteVars(flags *aflag.Flags) error {
	sdkv, err := GetWinSDKVersion(flags)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(flags.Output, "set_vars64.bat"), x64(sdkv, flags.Targetx64, flags.Targetx86, flags), 0o600)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(flags.Output, "set_vars32.bat"), x86(sdkv, flags.Targetx86, flags.Targetx64, flags), 0o600)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(flags.Output, "set_vars_arm64.bat"), x64(sdkv, flags.Targetarm64, flags.Targetarm, flags), 0o600)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(flags.Output, "set_vars_arm32.bat"), x86(sdkv, flags.Targetarm, flags.Targetarm64, flags), 0o600)
	if err != nil {
		return err
	}

	return nil
}

func x64(sdkv, targetA, targetB string, flags *aflag.Flags) []byte {
	return []byte(aflag.NewMSVCX64Vars(flags.MsvcVerLocal, sdkv, targetA, targetB, flags.Host))
}

func x86(sdkv, targetA, targetB string, flags *aflag.Flags) []byte {
	return []byte(aflag.NewMSVCX64Vars(flags.MsvcVerLocal, sdkv, targetA, targetB, flags.Host))
}
