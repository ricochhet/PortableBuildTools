LDFLAGS="-X 'main.buildDate=$(date)' -X 'main.gitHash=$(git rev-parse HEAD)' -X 'main.buildOn=$(go version)' -w -s "

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o downloader.exe -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o downloader-linux -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o downloader-linux-arm64 -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o downloader-darwin -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o downloader-darwin-arm64 -trimpath -ldflags "${LDFLAGS}"

# sha256
sha256sum downloader* > downloader-sha256
cat downloader-sha256

# chmod 
chmod +x downloader-*

# gzip
gzip --best downloader*