package main

import (
	"path/filepath"

	"github.com/ricochhet/portablebuildtools/download"
	"github.com/ricochhet/portablebuildtools/extract"
	aflag "github.com/ricochhet/portablebuildtools/flag"
)

//nolint:cyclop // composed of err checking. . .not complex
func main() {
	msvcPackages := aflag.SetPackages(flags, flags.SetMSVCPackages, aflag.MSVCPackages(flags))
	sdkPackages := aflag.SetPackages(flags, flags.SetWinSDKPackages, aflag.WinSDKPackages(flags))

	cwd, err := download.CreateDirectories(flags)
	if err != nil {
		panic(err)
	}

	flags.Downloads = filepath.Join(cwd, flags.Downloads)
	flags.DownloadsCRTD = filepath.Join(cwd, flags.DownloadsCRTD)
	flags.DownloadsDIA = filepath.Join(cwd, flags.DownloadsDIA)
	flags.Output = filepath.Join(cwd, flags.Output)
	msvcPackages, sdkPackages = aflag.AppendOptionals(msvcPackages, sdkPackages, flags)

	if flags.RewriteVars {
		if err := download.WriteVars(flags); err != nil {
			panic(err)
		}

		return
	}

	vsManifestJSON, err := download.GetManifest(flags)
	if err != nil {
		panic(err)
	}

	payloads, crtd, dia, sdk := download.GetPackages(flags, vsManifestJSON, msvcPackages)
	if err := download.GetPayloads(flags, payloads); err != nil {
		panic(err)
	}

	if err := download.GetWinSDK(flags, sdk, sdkPackages); err != nil {
		panic(err)
	}

	msvcv, err := download.GetMSVCVersion(flags)
	if err != nil {
		panic(err)
	}

	destx64 := filepath.Join(flags.Output, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetx64)
	destx86 := filepath.Join(flags.Output, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetx86)
	destarm := filepath.Join(flags.Output, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetarm)
	destarm64 := filepath.Join(flags.Output, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetarm64)

	if err := download.GetCRTD(crtd, destx64, destx86, destarm, destarm64, flags); err != nil {
		panic(err)
	}

	if err := download.GetDIASDK(dia, destx64, destx86, destarm, destarm64, flags); err != nil {
		panic(err)
	}

	if err := download.RemoveVCTipsTelemetry(flags); err != nil {
		panic(err)
	}

	if err := download.CleanHostDirectory(flags); err != nil {
		panic(err)
	}

	if err := download.WriteVars(flags); err != nil {
		panic(err)
	}

	if flags.CreateZip {
		if err := extract.Zip(flags.Output, flags.OutputZip); err != nil {
			panic(err)
		}
	}
}
