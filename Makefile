SHELL := /bin/bash

# ---- Project ----
BIN_NAME ?= semver
CMD_PKG  ?= .               # your main is at repo root (per your gox logs)
DIST     ?= dist

# Trim CR/LF so "1.2.3\n" works
VERSION  := $(shell tr -d '\r\n' < VERSION 2>/dev/null || echo dev)

# Target matrix (gox)
ARCHES ?= amd64 arm64
OSES ?= linux darwin windows
OSARCHS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Build flags (override to inject version vars if you have them)
LDFLAGS ?= -s -w
GOFLAGS ?= -trimpath -buildvcs=false
CGO_ENABLED ?= 0

.PHONY: all deps fmt vet test clean ensure_gox build dist checksum release debug

all: build

deps:
	@go mod download

fmt:
	@echo "gofmt..."
	@ret=0 && for d in $$(go list -f '{{.Dir}}' ./...); do \
		gofmt -l -w $$d/*.go || ret=$$? ; \
	done ; exit $$ret

vet:
	@echo "go vet..."
	@go vet ./...

test:
	@go test ./...

clean:
	@echo "Cleaning $(DIST)..."
	@go clean ./... || true
	@rm -rf "$(DIST)"

ensure_gox:
	@command -v gox >/dev/null 2>&1 || { \
		echo "Installing gox..."; \
		go install github.com/mitchellh/gox@latest; \
	}

# Build for current platform only
build: clean deps fmt vet test
	@mkdir -p "$(DIST)"
	CGO_ENABLED=$(CGO_ENABLED) go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o "$(DIST)/$(BIN_NAME)" $(CMD_PKG)
	@echo "Built $(DIST)/$(BIN_NAME)"

dist: clean deps fmt vet test ensure_gox
	@mkdir -p "$(DIST)"
	@echo ">> Cross-compiling with gox for OS=[${OSES}] ARCH=[${ARCHES}]"
	CGO_ENABLED=$(CGO_ENABLED) gox \
		-os="$(OSES)" \
		-arch="$(ARCHES)" \
		-ldflags "$(LDFLAGS)" \
		-output "$(DIST)/$(BIN_NAME)-{{.OS}}-{{.Arch}}" \
		$(CMD_PKG) || true

	@set -e; BUILT=0; SKIPPED=0; \
	for os in $(OSES); do \
	  for arch in $(ARCHES); do \
	    EXT=$$( [ "$$os" = "windows" ] && echo ".exe" || echo "" ); \
	    OUT="$(BIN_NAME)-$${os}-$${arch}$${EXT}"; \
	    if [ ! -f "$(DIST)/$$OUT" ]; then \
	      echo "⚠️  Skipping $$os/$$arch (no artifact produced; not an error)"; \
	      SKIPPED=$$((SKIPPED+1)); \
	      continue; \
	    fi; \
	    PKG="$(BIN_NAME)_$(VERSION)_$${os}_$${arch}"; \
	    echo ">> Packaging $$PKG"; \
	    mkdir -p "$(DIST)/$$PKG"; \
	    mv "$(DIST)/$$OUT" "$(DIST)/$$PKG/$(BIN_NAME)$${EXT}"; \
	    [ -f README.md ] && cp README.md "$(DIST)/$$PKG/" || true; \
	    [ -f LICENSE ] && cp LICENSE "$(DIST)/$$PKG/" || true; \
	    if [ "$$os" = "windows" ]; then \
	      (cd "$(DIST)" && zip -qr "$$PKG.zip" "$$PKG"); \
	    else \
	      (cd "$(DIST)" && tar -czf "$$PKG.tar.gz" "$$PKG"); \
	    fi; \
	    rm -rf "$(DIST)/$$PKG"; \
	    BUILT=$$((BUILT+1)); \
	  done; \
	done; \
	echo "Packaged $$BUILT target(s); skipped $$SKIPPED."

	$(MAKE) checksum

	@# ---- Friendly success summary (portable) ----
	@{ \
	  echo ""; \
	  count=$$(find "$(DIST)" -maxdepth 1 -type f \( -name '*.tar.gz' -o -name '*.zip' \) | wc -l | tr -d '[:space:]'); \
	  if [ "$$count" -gt 0 ]; then \
	    echo "✅ Success: built $$count artifact(s) for version $(VERSION)"; \
	    echo "   Output directory: $(DIST)/"; \
	    find "$(DIST)" -maxdepth 1 -type f \( -name '*.tar.gz' -o -name '*.zip' \) -print | sed 's/^/   - /'; \
	    if [ -f "$(DIST)/$(BIN_NAME)_$(VERSION)_SHA256SUMS.txt" ]; then \
	      echo "   Checksums: $(DIST)/$(BIN_NAME)_$(VERSION)_SHA256SUMS.txt"; \
	    fi; \
	    echo "   Upload with: gh release upload \"$(VERSION)\" $(DIST)/* --clobber"; \
	    echo ""; \
	  else \
	    echo "❌ No artifacts were produced. See logs above."; \
	    exit 1; \
	  fi; \
	}


checksum:
	@cd "$(DIST)"; \
	if command -v shasum >/dev/null 2>&1; then \
	  : > "$(BIN_NAME)_$(VERSION)_SHA256SUMS.txt"; \
	  for f in *.tar.gz *.zip; do [ -e "$$f" ] && shasum -a 256 "$$f" >> "$(BIN_NAME)_$(VERSION)_SHA256SUMS.txt"; done; \
	elif command -v sha256sum >/dev/null 2>&1; then \
	  : > "$(BIN_NAME)_$(VERSION)_SHA256SUMS.txt"; \
	  for f in *.tar.gz *.zip; do [ -e "$$f" ] && sha256sum "$$f" >> "$(BIN_NAME)_$(VERSION)_SHA256SUMS.txt"; done; \
	else \
	  echo "No shasum/sha256sum found; skipping checksums"; \
	fi; \
	echo "Checksums: $(DIST)/$(BIN_NAME)_$(VERSION)_SHA256SUMS.txt"

# Upload everything in ./dist to the GitHub Release tag == VERSION
release:
	@gh release upload "$(VERSION)" $(DIST)/* --clobber
	@gh release view "$(VERSION)"

debug:
	$(info BIN_NAME=$(BIN_NAME))
	$(info CMD_PKG=$(CMD_PKG))
	$(info DIST=$(DIST))
	$(info VERSION=$(VERSION))
	$(info OSARCHS=$(OSARCHS))
	$(info LDFLAGS=$(LDFLAGS))
	@true
