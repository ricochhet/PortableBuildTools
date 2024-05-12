SET LDFLAGS="-w -s"
go build -o PortableBuildTools.exe -trimpath -ldflags %LDFLAGS%