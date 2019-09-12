.DEFAULT_GOAL := default
TEST_PATH ?= "TestAcc"
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
$(eval REMAINDER := $$$(ELK_VERSION))

default: build test

check-elk-version:
ifndef ELK_VERSION
	$(error ELK_VERSION is undefined)
endif

check-kibana-type:
ifndef KIBANA_TYPE
	$(error KIBANA_TYPE is undefined)
endif

travisbuild: deps default

test: fmtcheck docker-build
	TF_ACC=1 go test -v ./kibana -run $(TEST_PATH)

build: fmtcheck vet test
	@go install
	@mkdir -p ~/.terraform.d/plugins/
	@cp $(GOPATH)/bin/terraform-provider-kibana ~/.terraform.d/plugins/terraform-provider-kibana
	@echo "Build succeeded"

docker-build: check-elk-version check-kibana-type
	$(eval MAIN_VERSION := $(shell echo $(ELK_VERSION) | head -c 3))

	@echo building docker ELK_VERSION:$(ELK_VERSION) KIBANA_TYPE: $(KIBANA_TYPE) MAIN_VERSION: $(MAIN_VERSION)
	@if [ "$(ELK_VERSION)" = "./..." ]; then \
		echo "ERROR: Set ELK_VERSION to a specific version. For example,"; \
		echo "  make docker-build"; \
		exit 1; \
	fi
	@if [ "$(KIBANA_TYPE)" = "KibanaTypeVanilla" ]; then \
		cd docker/elasticsearch && docker build --no-cache --build-arg ELK_VERSION=$(ELK_VERSION) --build-arg ELK_PACK=$(ELK_PACK) --build-arg MAKELOGS_VERSION=$(MAKELOGS_VERSION) . -t elastic-local:$(ELK_VERSION); \
	fi

start-kibana: docker-build
	@sh -c "'$(CURDIR)/scripts/start-docker.sh'"

build-gox: deps fmtcheck vet
	gox -osarch="linux/amd64 windows/amd64 darwin/amd64" \
	-output="pkg/{{.OS}}_{{.Arch}}/terraform-provider-kibana" .

deps:
	go get -u golang.org/x/net/context; \
    go get -u github.com/mitchellh/gox; \

clean:
	rm -rf pkg/
fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

.PHONY: build test vet fmt fmtcheck errcheck vendor-status test-compile release docker-build start-kibana
