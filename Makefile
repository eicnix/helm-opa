HELM_HOME ?= $(shell helm home)
HELM_PLUGIN_DIR ?= $(HELM_HOME)/plugins/helm-opa
HAS_DEP := $(shell command -v DEP;)
VERSION := $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' plugin.yaml)
DIST := $(CURDIR)/_dist
LDFLAGS := "-X main.version=${VERSION}"

.PHONY: install
install: bootstrap build
	cp opa $(HELM_PLUGIN_DIR)
	cp plugin.yaml $(HELM_PLUGIN_DIR)

.PHONY: hookInstall
hookInstall: bootstrap build

.PHONY: build
build:
	go build -o opa -ldflags $(LDFLAGS) ./main.go

.PHONY: dist
dist:
	mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 go build -o opa -ldflags $(LDFLAGS) ./main.go
	tar -zcvf $(DIST)/helm-opa-linux-$(VERSION).tgz opa README.md LICENSE.txt plugin.yaml
	GOOS=darwin GOARCH=amd64 go build -o opa -ldflags $(LDFLAGS) ./main.go
	tar -zcvf $(DIST)/helm-opa-macos-$(VERSION).tgz opa README.md LICENSE.txt plugin.yaml
	GOOS=windows GOARCH=amd64 go build -o opa.exe -ldflags $(LDFLAGS) ./main.go
	tar -zcvf $(DIST)/helm-opa-windows-$(VERSION).tgz opa.exe README.md LICENSE.txt plugin.yaml

.PHONY: bootstrap
bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
	dep ensure