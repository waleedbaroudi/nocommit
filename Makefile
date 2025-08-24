# Makefile for nocommit releases

PKG      := github.com/waleedbaroudi/nocommit
VERSION  ?= $(shell git describe --tags --always --dirty)
DIST     := dist

.PHONY: all clean build checksums release

all: clean build checksums

clean:
	rm -rf $(DIST)

# Build all target artifacts
build: $(DIST)/nocommit_darwin_amd64.tar.gz \
       $(DIST)/nocommit_darwin_arm64.tar.gz \
       $(DIST)/nocommit_windows_amd64.zip \
       $(DIST)/nocommit_windows_arm64.zip

# ----------------------
# macOS (amd64)
$(DIST)/nocommit_darwin_amd64.tar.gz:
	mkdir -p $(DIST)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath \
		-ldflags "-s -w -X $(PKG)/cmd.Version=$(VERSION)" \
		-o $(DIST)/nocommit .
	tar -C $(DIST) -czf $@ nocommit
	rm $(DIST)/nocommit

# macOS (arm64)
$(DIST)/nocommit_darwin_arm64.tar.gz:
	mkdir -p $(DIST)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath \
		-ldflags "-s -w -X $(PKG)/cmd.Version=$(VERSION)" \
		-o $(DIST)/nocommit .
	tar -C $(DIST) -czf $@ nocommit
	rm $(DIST)/nocommit

# ----------------------
# Windows (amd64)
$(DIST)/nocommit_windows_amd64.zip:
	mkdir -p $(DIST)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath \
		-ldflags "-s -w -X $(PKG)/cmd.Version=$(VERSION)" \
		-o $(DIST)/nocommit.exe .
	cd $(DIST) && zip -9 nocommit_windows_amd64.zip nocommit.exe
	rm $(DIST)/nocommit.exe

# Windows (arm64)
$(DIST)/nocommit_windows_arm64.zip:
	mkdir -p $(DIST)
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -trimpath \
		-ldflags "-s -w -X $(PKG)/cmd.Version=$(VERSION)" \
		-o $(DIST)/nocommit.exe .
	cd $(DIST) && zip -9 nocommit_windows_arm64.zip nocommit.exe
	rm $(DIST)/nocommit.exe

# ----------------------
# Checksums for all artifacts
checksums:
	cd $(DIST) && (sha256sum nocommit_darwin_*.tar.gz nocommit_windows_*.zip || shasum -a 256 nocommit_darwin_*.tar.gz nocommit_windows_*.zip) > SHA256SUMS.txt

release: all
	@echo "✅ Artifacts ready in $(DIST)/"
	@ls -lh $(DIST)
