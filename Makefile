# Makefile — macOS, Linux, Windows release artifacts only

PKG      := github.com/waleedbaroudi/nocommit
VERSION  ?= $(shell git describe --tags --always --dirty)
DIST     := dist
LDFLAGS  := -s -w -X $(PKG)/cmd.Version=$(VERSION)

.PHONY: all clean build checksums release

all: release

release: clean build checksums
	@echo "✅ Artifacts in $(DIST)/"
	@ls -lh $(DIST)

clean:
	rm -rf $(DIST)

# -------------------------------------------------------
# Build everything
build: \
  $(DIST)/nocommit_darwin_amd64.tar.gz \
  $(DIST)/nocommit_darwin_arm64.tar.gz \
  $(DIST)/nocommit_linux_amd64.tar.gz  \
  $(DIST)/nocommit_linux_arm64.tar.gz  \
  $(DIST)/nocommit_windows_amd64.zip   \
  $(DIST)/nocommit_windows_arm64.zip

# -------------------- macOS ----------------------------
$(DIST)/nocommit_darwin_amd64.tar.gz:
	mkdir -p $(DIST)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST)/nocommit .
	tar -C $(DIST) -czf $@ nocommit
	rm -f $(DIST)/nocommit

$(DIST)/nocommit_darwin_arm64.tar.gz:
	mkdir -p $(DIST)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST)/nocommit .
	tar -C $(DIST) -czf $@ nocommit
	rm -f $(DIST)/nocommit

# -------------------- Linux ------------------
$(DIST)/nocommit_linux_amd64.tar.gz:
	mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST)/nocommit .
	tar -C $(DIST) -czf $@ nocommit
	rm -f $(DIST)/nocommit

$(DIST)/nocommit_linux_arm64.tar.gz:
	mkdir -p $(DIST)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST)/nocommit .
	tar -C $(DIST) -czf $@ nocommit
	rm -f $(DIST)/nocommit

# -------------------- Windows --------------------------
$(DIST)/nocommit_windows_amd64.zip:
	mkdir -p $(DIST)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST)/nocommit.exe .
	cd $(DIST) && zip -9 nocommit_windows_amd64.zip nocommit.exe
	rm -f $(DIST)/nocommit.exe

$(DIST)/nocommit_windows_arm64.zip:
	mkdir -p $(DIST)
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST)/nocommit.exe .
	cd $(DIST) && zip -9 nocommit_windows_arm64.zip nocommit.exe
	rm -f $(DIST)/nocommit.exe

# -------------------- Checksums ------------------------
checksums:
	cd $(DIST) && (sha256sum nocommit_darwin_*.tar.gz nocommit_linux_*.tar.gz nocommit_windows_*.zip || shasum -a 256 nocommit_darwin_*.tar.gz nocommit_linux_*.tar.gz nocommit_windows_*.zip) > SHA256SUMS.txt
