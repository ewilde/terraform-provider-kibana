.DEFAULT_GOAL := default

TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
$(eval REMAINDER := $$$(ELK_VERSION))
MAIN_VERSION := $(shell echo $(ELK_VERSION) | head -c 3)

default: build test

build: fmtcheck errcheck vet
	go install

test: docker-build fmtcheck
	go test -v ./...

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

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./aws"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

docker-build:
	@if [ "$(ELK_VERSION)" = "./..." ]; then \
		echo "ERROR: Set ELK_VERSION to a specific version. For example,"; \
		echo "  make docker-build"; \
		exit 1; \
	fi
	@if [ "$(KIBANA_TYPE)" = "KibanaTypeVanilla" ]; then \
		cd docker/elasticsearch-$(MAIN_VERSION) && docker build . -t elastic-local:$(ELK_VERSION); \
	fi

kibana-start: docker-build
	@sh -c "'$(CURDIR)/scripts/start-kibana-$(ELK_VERSION).sh'"
.PHONY: build docker-build kibana-start test testacc vet fmt fmtcheck errcheck vendor-status test-compile
