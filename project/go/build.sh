LDFLAGS="-X 'main.buildDate=$(date)' -X 'main.gitHash=$(git rev-parse HEAD)' -X 'main.buildOn=$(go version)' -w -s "

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o PortableBuildTools.exe -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o PortableBuildTools-linux -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o PortableBuildTools-linux-arm64 -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o PortableBuildTools-darwin -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o PortableBuildTools-darwin-arm64 -trimpath -ldflags "${LDFLAGS}"

# sha256
sha256sum PortableBuildTools* > PortableBuildTools-sha256
cat PortableBuildTools-sha256

# chmod 
chmod +x PortableBuildTools-*

# gzip
gzip --best PortableBuildTools*