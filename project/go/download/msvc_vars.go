package download

import (
	"os"
	"path/filepath"

	aflag "github.com/ricochhet/sdkstandalone/flag"
)

func WriteVars(flags *aflag.Flags) error {
	msvcv, err := GetMSVCVersion(flags)
	if err != nil {
		return err
	}

	sdkv, err := GetWinSDKVersion(flags)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(flags.Output, "set_vars64.bat"), x64(msvcv, sdkv, flags.Targetx64, flags.Targetx86, flags), 0o600)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(flags.Output, "set_vars32.bat"), x86(msvcv, sdkv, flags.Targetx86, flags.Targetx64, flags), 0o600)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(flags.Output, "set_vars_arm64.bat"), x64(msvcv, sdkv, flags.Targetarm64, flags.Targetarm, flags), 0o600)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(flags.Output, "set_vars_arm32.bat"), x86(msvcv, sdkv, flags.Targetarm, flags.Targetarm64, flags), 0o600)
	if err != nil {
		return err
	}

	return nil
}

func x64(msvcv, sdkv, targetA, targetB string, flags *aflag.Flags) []byte {
	return []byte(aflag.NewMSVCX64Vars(msvcv, sdkv, targetA, targetB, flags.Host))
}

func x86(msvcv, sdkv, targetA, targetB string, flags *aflag.Flags) []byte {
	return []byte(aflag.NewMSVCX86Vars(msvcv, sdkv, targetA, targetB, flags.Host))
}
