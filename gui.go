/*
 * PortableBuildTools
 * Copyright (C) 2024 PortableBuildTools contributors
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

//nolint:mnd // wontfix
package main

import (
	"bytes"
	"io"
	"log"

	"github.com/AllenDang/giu"
	"github.com/ricochhet/minicommon/logger"
)

var (
	sashPos1 float32      = 500 //nolint:gochecknoglobals // wontfix
	sashPos2 float32      = 300 //nolint:gochecknoglobals // wontfix
	buf      bytes.Buffer       //nolint:gochecknoglobals // wontfix
)

func window() {
	giu.SingleWindow().Layout(
		giu.Condition(working, giu.Label("Working..."), nil),
		giu.SplitLayout(giu.DirectionHorizontal, &sashPos2,
			giu.Layout{
				giu.SplitLayout(giu.DirectionVertical, &sashPos1,
					giu.Layout{
						giu.InputText(&flags.MsvcVer).Label("MSVC Version"),
						giu.InputText(&flags.WinSdkVer).Label("Windows SDK Version"),
						giu.InputText(&flags.Dest).Label("Output"),
						giu.InputText(&flags.DestZip).Label("Output (Zip)"),
						giu.InputText(&flags.TmpPath).Label("Temp Path"),
						giu.InputText(&flags.TmpCrtd).Label("Temp Path (CRTD)"),
						giu.InputText(&flags.TmpDia).Label("Temp Path (DIA)"),
						giu.InputText(&flags.Host).Label("Host"),
					},
					giu.Layout{
						giu.Checkbox("Enable Spectre Libraries", &flags.SpectreLibs),
						giu.Checkbox("Enable ARM/ARM64 Targets", &flags.ArmTargets),
						giu.Checkbox("Enable LLVM Clang", &flags.LlvmClang),
						giu.Checkbox("Enable UnitTest Framework", &flags.UnitTest),
						giu.Checkbox("Enable Cmake", &flags.Cmake),
						giu.Checkbox("Enable MFC ATL", &flags.MfcAtl),
						giu.Checkbox("Enable VCPKG", &flags.Vcpkg),
						giu.Checkbox("Enable MSBuild", &flags.Msbuild),
						giu.Checkbox("Create zip after download", &flags.Zip),
						giu.Separator(),
						giu.Row(
							giu.Button("Write Environment").OnClick(writeEnvironments).Disabled(working),
							giu.Button("Download").OnClick(run).Disabled(working),
						),
					},
				),
			}, giu.Layout{
				giu.Child().Layout(
					giu.Label(buf.String()),
				),
			}),
	)
}

func Gui(ver string, logfile io.Writer) {
	logger.SharedLogger = logger.NewLogger(4, logger.InfoLevel, io.MultiWriter(logfile, &buf), log.Lshortfile|log.LstdFlags)
	logger.SharedLogger.Info("Initialized!")

	wnd := giu.NewMasterWindow("PortableBuildTools - "+ver, 840, 500, 0)
	wnd.Run(window)
}
