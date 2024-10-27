//nolint:mnd // wontfix
package main

import (
	"github.com/AllenDang/giu"
)

var sashPos1 float32 = 500 //nolint:gochecknoglobals // wontfix

func window() {
	giu.SingleWindow().Layout(
		giu.Condition(working, giu.Label("Working..."), nil),
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
	)
}

func Gui(ver string) {
	wnd := giu.NewMasterWindow("PortableBuildTools - %s"+ver, 840, 300, 0)
	wnd.Run(window)
}
