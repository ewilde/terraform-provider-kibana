.DEFAULT_GOAL := default

TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
$(eval REMAINDER := $$$(ELK_VERSION))
MAIN_VERSION := $(shell echo $(ELK_VERSION) | head -c 3)
GO111MODULE := off

default: build test

build: fmtcheck errcheck vet
	go install

test: docker-build fmtcheck
	go test -v $(go list ./... | grep -v /vendor/)

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@govendor status

docker-build:
	@if [ "$(ELK_VERSION)" = "./..." ]; then \
		echo "ERROR: Set ELK_VERSION to a specific version. For example,"; \
		echo "  make docker-build"; \
		exit 1; \
	fi
	@if [ "$(KIBANA_TYPE)" = "KibanaTypeVanilla" ]; then \
		cd docker/elasticsearch && docker build --build-arg ELK_VERSION=$(ELK_VERSION) --build-arg ELK_PACK=$(ELK_PACK) . -t elastic-local:$(ELK_VERSION); \
	fi

kibana-start: docker-build
	@sh -c "'$(CURDIR)/scripts/start-kibana-$(ELK_VERSION).sh'"
.PHONY: build docker-build kibana-start test testacc vet fmt fmtcheck errcheck vendor-status test-compile
