LDFLAGS="-X 'main.buildDate=$(date)' -X 'main.gitHash=$(git rev-parse HEAD)' -X 'main.buildOn=$(go version)' -w -s "

CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o portablebuildtools.exe -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o portablebuildtools-linux -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o portablebuildtools-linux-arm64 -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o portablebuildtools-darwin -trimpath -ldflags "${LDFLAGS}"
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o portablebuildtools-darwin-arm64 -trimpath -ldflags "${LDFLAGS}"

# sha256
sha256sum portablebuildtools* > portablebuildtools-sha256
cat portablebuildtools-sha256

# chmod 
chmod +x portablebuildtools-*

# gzip
gzip --best portablebuildtools*