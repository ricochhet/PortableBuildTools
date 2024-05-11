package flag

type Flags struct {
	MSVC_VERSION          string
	MSVC_VERSION_LOCAL    string
	SDK_PID               string
	OUTPUT                string // output folder
	DOWNLOADS             string // temporary download files
	DOWNLOADS_CRTD        string
	DOWNLOADS_DIA         string
	HOST                  string // or x86
	SET_MSVC_PACKAGES     string
	SET_WINSDK_PACKAGES   string
	DOWNLOAD_SPECTRE_LIBS bool
	DOWNLOAD_ARM_TARGETS  bool
	DOWNLOAD_LLVM_CLANG   bool
	DOWNLOAD_UNITTEST     bool
	DOWNLOAD_CMAKE        bool
	MANIFEST_URL          string
	MANIFEST_PREVIEW_URL  string
	TARGETX64             string
	TARGETX86             string
	TARGETARM             string
	TARGETARM64           string
	REWRITE_VARS          bool
	MSIEXEC_VERBOSE       bool
}
