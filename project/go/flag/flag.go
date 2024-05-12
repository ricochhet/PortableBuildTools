package flag

type Flags struct {
	MsvcVer             string
	MsvcVerLocal        string
	WinSDKVer           string
	Output              string // output folder
	Downloads           string // temporary download files
	DownloadsCRTD       string
	DownloadsDIA        string
	Host                string // or x86
	SetMSVCPackages     string
	SetWinSDKPackages   string
	DownloadSpectreLibs bool
	DownloadARMTargets  bool
	DownloadLLVMClang   bool
	DownloadUnitTest    bool
	DownloadCmake       bool
	ManifestURL         string
	ManifestPreviewURL  string
	Targetx64           string
	Targetx86           string
	Targetarm           string
	Targetarm64         string
	RewriteVars         bool
	MSIExtractVerbose   bool
}
