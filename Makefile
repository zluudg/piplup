NAME:=piplup

#######################################
# VERSION SOURCE OF TRUTH FOR PROJECT #
#######################################
VERSION:=0.0.0

OUT:=./out
DEFAULT_INSTALLDIR:=/usr/bin
INSTALL:=install -p -m 0755
COMMIT:=$$(cat COMMIT 2> /dev/null || git describe --dirty=+WiP --always 2> /dev/null)

.PHONY: build outdir install clean tarball fmt vet coverage container itest


all: build

build: outdir
	go build -v -ldflags "-X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)'" -o $(OUT)/ ./cmd/...

install:
	test -z "$(DESTDIR)" && $(INSTALL) $(OUT)/$(PROG) $(DEFAULT_INSTALLDIR) || $(INSTALL) $(OUT)/$(PROG) $(DESTDIR)$(prefix)/bin/

outdir:
	-mkdir -p $(OUT)

test:
	go test ./...

coverage: outdir
	go test -coverprofile=$(OUT)/coverage.out ./...
	go tool cover -html="$(OUT)/coverage.out" -o $(OUT)/coverage.html

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	-rm -rf $(OUT)
