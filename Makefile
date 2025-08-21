# Makefile for nocommit releases

PKG      := github.com/waleedbaroudi/nocommit
VERSION  ?= $(shell git describe --tags --always --dirty)
DIST     := dist

all: clean build checksums

clean:
	rm -rf $(DIST)

# macOS binaries
build: $(DIST)/nocommit_darwin_amd64.tar.gz \
       $(DIST)/nocommit_darwin_arm64.tar.gz

$(DIST)/nocommit_darwin_amd64.tar.gz:
	mkdir -p $(DIST)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath \
		-ldflags "-s -w -X $(PKG)/cmd.Version=$(VERSION)" \
		-o $(DIST)/nocommit .
	tar -C $(DIST) -czf $@ nocommit
	rm $(DIST)/nocommit

$(DIST)/nocommit_darwin_arm64.tar.gz:
	mkdir -p $(DIST)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath \
		-ldflags "-s -w -X $(PKG)/cmd.Version=$(VERSION)" \
		-o $(DIST)/nocommit .
	tar -C $(DIST) -czf $@ nocommit
	rm $(DIST)/nocommit

checksums:
	cd $(DIST) && (sha256sum nocommit_darwin_*.tar.gz || shasum -a 256 nocommit_darwin_*.tar.gz) > SHA256SUMS.txt

release: all
	@echo "✅ Artifacts ready in $(DIST)/"
	@ls -lh $(DIST)
