NAME:=piplup

#######################################
# VERSION SOURCE OF TRUTH FOR PROJECT #
#######################################
VERSION:=0.0.0

OUT:=./out
DEFAULT_INSTALLDIR:=/usr/bin
INSTALL:=install -p -m 0755
COMMIT:=$$(cat COMMIT 2> /dev/null || git describe --dirty=+WiP --always 2> /dev/null)

ifneq ("$(wildcard ./COMMIT))","")
	VENDOR = "-mod=vendor"
else
	VENDOR = ""
endif

.PHONY: build outdir install clean tarball fmt vet coverage


all: build

build: outdir
	go build -v $(VENDOR) -ldflags "-X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)'" -o $(OUT)/ ./cmd/...

install:
	test -z "$(DESTDIR)" && $(INSTALL) $(OUT)/$(NAME) $(DEFAULT_INSTALLDIR) || $(INSTALL) $(OUT)/$(NAME) $(DESTDIR)$(prefix)/bin/

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
	-rm -rf vendor/

tarball: outdir
	echo "$(COMMIT)" > $(OUT)/COMMIT
	git archive --format=tar.gz --prefix=$(NAME)-$(VERSION)/ -o $(OUT)/$(NAME)-$(VERSION).tar.gz --add-file $(OUT)/COMMIT HEAD
