_ != mkdir -p .make

GO     ?= go
DEVCTL ?= $(GO) tool devctl
DOCKER ?= docker
GINKGO ?= $(GO) tool ginkgo

GO_SRC != $(DEVCTL) list --go
GO_SRC += main.go # Bug in devctl I think...

build: bin/pulumi2crd
test: .make/ginkgo-run
tidy: go.sum
docker: .make/docker-build

bin/pulumi2crd: go.mod ${GO_SRC}
	$(GO) build -o $@

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

.make/docker-build: Dockerfile .dockerignore ${GO_SRC}
	$(DOCKER) build . -f $<
	@touch $@

.make/ginkgo-run: ${GO_SRC}
	$(GINKGO) $(sort $(dir $?))
	@touch $@
