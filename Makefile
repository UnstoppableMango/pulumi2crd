_ != mkdir -p .make

GO         ?= go
DEVCTL     ?= $(GO) tool devctl
DOCKER     ?= docker
DPRINT     ?= ${CURDIR}/bin/dprint
GINKGO     ?= $(GO) tool ginkgo
GOLINT     ?= $(GO) tool golangci-lint
GORELEASER ?= goreleaser

GO_SRC != $(DEVCTL) list --go
GO_SRC += main.go # Bug in devctl I think...

build: bin/pulumi2crd
test: .make/ginkgo-run
fmt format: .make/go-fmt .make/dprint-fmt
lint: .make/go-vet .make/golangci-lint-run
tidy: go.sum
docker: .make/docker-build

bin/pulumi2crd: go.mod ${GO_SRC}
	$(GO) build -o $@

bin/dprint: .versions/dprint | .make/dprint/install.sh
	DPRINT_INSTALL=${CURDIR} .make/dprint/install.sh $(shell $(DEVCTL) v dprint)
	@touch $@

%_suite_test.go: ## Bootstrap a Ginkgo test suite
	cd $(dir $@) && $(GINKGO) bootstrap
%_test.go: ## Generate a Ginkgo test
	cd $(dir $@) && $(GINKGO) generate $(notdir $@)

go.sum: go.mod ${GO_SRC}
	$(GO) mod tidy

go.work:
	$(GO) work init
	$(GO) work use .

go.work.sum: go.work
	$(GO) work sync

.envrc: hack/example.envrc ## Generate a recommended .envrc
	cp $< $@ && chmod a=,u=r $@

.make/docker-build: Dockerfile .dockerignore ${GO_SRC}
	$(DOCKER) build . -f $<
	@touch $@

.make/dprint/install.sh:
	@mkdir -p $(dir $@)
	curl -fsSL https://dprint.dev/install.sh -o $@
	@chmod +x $@

JSON_SRC := .dprint.json .github/renovate.json .vscode/extensions.json
MD_SRC   := README.md

.make/dprint-fmt: ${JSON_SRC} ${MD_SRC} | bin/dprint
	$(DPRINT) fmt --allow-no-files $?
	@touch $@

.make/ginkgo-run: ${GO_SRC}
	$(GINKGO) $(sort $(dir $?))
	@touch $@

.make/go-fmt: ${GO_SRC}
	$(GO) fmt $(addprefix ./,$(sort $(dir $?)))
	@touch $@

.make/go-vet: ${GO_SRC}
	$(GO) vet $(addprefix ./,$(sort $(dir $?)))
	@touch $@

.make/golangci-lint-run: ${GO_SRC}
	$(GOLINT) run
	@touch $@
