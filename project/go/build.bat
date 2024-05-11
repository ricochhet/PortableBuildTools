SET LDFLAGS="-w -s"
go build -o downloader.exe -trimpath -ldflags %LDFLAGS%