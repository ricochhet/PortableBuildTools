LDFLAGS=-X 'main.buildDate=$(shell date)' -X 'main.gitHash=$(shell git rev-parse HEAD)' -X 'main.buildOn=$(shell go version)' -w -s

GO_BUILD=go build -trimpath -ldflags "$(LDFLAGS)"

.PHONY: all fmt lint test deadcode portablebuildtools msiextract postbuild clean

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

msiextract:
	cargo build --release --manifest-path msiextract/Cargo.toml

portablebuildtools:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO_BUILD) -o portablebuildtools.exe

postbuild:
	mkdir -p build/
	cp portablebuildtools.exe build/
	cp msiextract/target/release/msiextract.exe build/

clean:
	rm -f portablebuildtools-windows.exe