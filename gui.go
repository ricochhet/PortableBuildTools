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

	"github.com/AllenDang/giu"
	"github.com/ricochhet/minicommon/charmbracelet"
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
						giu.Style().SetDisabled(working).To(giu.InputText(&flags.MsvcVer).Label("MSVC Version")),
						giu.Style().SetDisabled(working).To(giu.InputText(&flags.WinSdkVer).Label("Windows SDK Version")),
						giu.Style().SetDisabled(working).To(giu.InputText(&flags.Dest).Label("Output")),
						giu.Style().SetDisabled(working).To(giu.InputText(&flags.DestZip).Label("Output (Zip)")),
						giu.Style().SetDisabled(working).To(giu.InputText(&flags.TmpPath).Label("Temp Path")),
						giu.Style().SetDisabled(working).To(giu.InputText(&flags.TmpCrtd).Label("Temp Path (CRTD)")),
						giu.Style().SetDisabled(working).To(giu.InputText(&flags.TmpDia).Label("Temp Path (DIA)")),
						giu.Style().SetDisabled(working).To(giu.Tooltip("Uncheck to manually specify absolute paths outside of the CWD.")),
						giu.Style().SetDisabled(working).To(giu.Checkbox("Use CWD", &flags.Cwd)),
						giu.Style().SetDisabled(working).To(giu.InputText(&flags.Host).Label("Host")),
					},
					giu.Layout{
						giu.Style().SetDisabled(working).To(giu.Checkbox("Enable Spectre Libraries", &flags.SpectreLibs)),
						giu.Style().SetDisabled(working).To(giu.Checkbox("Enable ARM/ARM64 Targets", &flags.ArmTargets)),
						giu.Style().SetDisabled(working).To(giu.Checkbox("Enable LLVM Clang", &flags.LlvmClang)),
						giu.Style().SetDisabled(working).To(giu.Checkbox("Enable UnitTest Framework", &flags.UnitTest)),
						giu.Style().SetDisabled(working).To(giu.Checkbox("Enable Cmake", &flags.Cmake)),
						giu.Style().SetDisabled(working).To(giu.Checkbox("Enable MFC ATL", &flags.MfcAtl)),
						giu.Style().SetDisabled(working).To(giu.Checkbox("Enable VCPKG", &flags.Vcpkg)),
						giu.Style().SetDisabled(working).To(giu.Checkbox("Enable MSBuild", &flags.Msbuild)),
						giu.Style().SetDisabled(working).To(giu.Checkbox("Create zip after download", &flags.Zip)),
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
	charmbracelet.SharedLogger = charmbracelet.NewMultiLogger(logfile, &buf)
	charmbracelet.SharedLogger.Info("Initialized!")

	wnd := giu.NewMasterWindow("PortableBuildTools - "+ver, 840, 500, 0)
	wnd.Run(window)
}
