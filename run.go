//nolint:mnd // wontfix
package main

import (
	"path/filepath"

	"github.com/AllenDang/giu"
	"github.com/ricochhet/minicommon/logger"
	"github.com/ricochhet/minicommon/zip"
	"github.com/ricochhet/portablebuildtools/download"
	aflag "github.com/ricochhet/portablebuildtools/flag"
	"github.com/ricochhet/portablebuildtools/internal"
)

var working bool //nolint:gochecknoglobals // wontfix

func run() {
	errCh := make(chan error, 3)

	working = true

	go func() {
		logger.SharedLogger.Info("... DOWNLOADING")

		go runWerr(errCh)

		for err := range errCh {
			if err != nil {
				logger.SharedLogger.Errorf("ERROR: %v", err)

				working = false
			}
		}

		working = false

		logger.SharedLogger.Info("... DONE")

		giu.Update()
	}()
}

func runWerr(errCh chan<- error) { //nolint:funlen,cyclop // wontfix
	defer close(errCh)

	if _, _, err := internal.FindMsiExtract(); err != nil {
		errCh <- err

		return
	}

	msvcPackages := aflag.SetPackages(flags, flags.SetMsvcPackages, aflag.MsvcPackages(flags))
	sdkPackages := aflag.SetPackages(flags, flags.SetWinSdkPackages, aflag.WinSdkPackages(flags))

	cwd, err := internal.CreateDirectories(flags)
	if err != nil {
		errCh <- err

		return
	}

	flags.TmpPath = filepath.Join(cwd, flags.TmpPath)
	flags.TmpCrtd = filepath.Join(cwd, flags.TmpCrtd)
	flags.TmpDia = filepath.Join(cwd, flags.TmpDia)
	flags.Dest = filepath.Join(cwd, flags.Dest)
	msvcPackages, sdkPackages = aflag.AppendOptionals(msvcPackages, sdkPackages, flags)

	vsManifestJSON, err := download.GetManifest(flags)
	if err != nil {
		errCh <- err

		return
	}

	payloads, crtd, dia, sdk := download.GetPackages(flags, vsManifestJSON, msvcPackages)

	if err := download.GetPayloads(flags, payloads); err != nil {
		errCh <- err

		return
	}

	if err := download.GetWinSdk(flags, sdk, sdkPackages); err != nil {
		errCh <- err

		return
	}

	msvcv, err := internal.GetMsvcVersion(flags)
	if err != nil {
		errCh <- err

		return
	}

	destx64 := filepath.Join(flags.Dest, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetx64)
	destx86 := filepath.Join(flags.Dest, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetx86)
	destarm := filepath.Join(flags.Dest, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetarm)
	destarm64 := filepath.Join(flags.Dest, "VC", "Tools", "MSVC", msvcv, "bin", "Host"+flags.Host, flags.Targetarm64)

	if err := download.GetCrtd(crtd, destx64, destx86, destarm, destarm64, flags); err != nil {
		errCh <- err

		return
	}

	if err := download.GetDiaSdk(dia, destx64, destx86, destarm, destarm64, flags); err != nil {
		errCh <- err

		return
	}

	if err := internal.RemoveVcTipsTelemetry(flags); err != nil {
		errCh <- err

		return
	}

	if err := internal.CleanHostDirectory(flags); err != nil {
		errCh <- err

		return
	}

	if err := internal.WriteEnvironment(flags); err != nil {
		errCh <- err

		return
	}

	if err := internal.CopyInstances(flags); err != nil {
		errCh <- err

		return
	}

	if flags.Zip {
		if err := zip.Zip(flags.Dest, flags.DestZip); err != nil {
			errCh <- err

			return
		}
	}
}

func writeEnvironments() {
	errCh := make(chan error, 3)

	working = true

	go func() {
		logger.SharedLogger.Info("... WRITING")

		go writeEnvironmentsWerr(errCh)

		for err := range errCh {
			if err != nil {
				logger.SharedLogger.Errorf("ERROR: %v", err)
			}
		}

		working = false

		logger.SharedLogger.Info("... DONE")

		giu.Update()
	}()
}

func writeEnvironmentsWerr(errCh chan<- error) {
	if flags.WriteEnvironment {
		if err := internal.WriteEnvironment(flags); err != nil {
			errCh <- err

			return
		}

		errCh <- nil

		return
	}
}
