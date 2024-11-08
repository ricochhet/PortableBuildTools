LDFLAGS=-X 'main.buildDate=$(shell date)' -X 'main.gitHash=$(shell git rev-parse --short HEAD)' -X 'main.buildOn=$(shell go version)' -w -s -H=windowsgui

GO_BUILD=go build -trimpath -ldflags "$(LDFLAGS)"

.PHONY: all fmt lint test deadcode portablebuildtools msiextract postbuild create_release clean

all: portablebuildtools msiextract postbuild

fmt:
	gofumpt -l -w .

mod:
	go get -u
	go mod tidy


lint:
	golangci-lint run

test:
	go test ./...

deadcode:
	deadcode ./...

syso:
	windres portablebuildtools.rc -O coff -o portablebuildtools.syso

msiextract:
	cargo build --release --manifest-path msiextract/Cargo.toml

portablebuildtools:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 $(GO_BUILD) -o portablebuildtools.exe

postbuild:
	mkdir -p build/
	cp portablebuildtools.exe build/
	cp msiextract/target/release/msiextract.exe build/
	cp -r _Instances/ build/_Instances/

create_release:
	mkdir -p .releases/$(shell git rev-parse HEAD)
	cp portablebuildtools.exe .releases/$(shell git rev-parse HEAD)/
	cp msiextract/target/release/msiextract.exe .releases/$(shell git rev-parse HEAD)/
	cp -r _Instances/ .releases/$(shell git rev-parse HEAD)/_Instances/
	cd .releases/$(shell git rev-parse HEAD)
	find .releases/$(shell git rev-parse HEAD) -type f -exec sha256sum {} \; > .releases/$(shell git rev-parse HEAD)/portablebuildtools-sha256
	tar -czf .releases/PortableBuildTools.tar.gz -C .releases/$(shell git rev-parse HEAD) .

clean:
	rm -f portablebuildtools-windows.exe