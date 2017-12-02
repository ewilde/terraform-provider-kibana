GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: build testacc

travisbuild: deps default

testacc: fmtcheck docker
	TF_ACC=1 go test -v ./kibana -run="TestAcc"

build: fmtcheck vet testacc
	@go install
	@mkdir -p ~/.terraform.d/plugins/
	@cp $(GOPATH)/bin/terraform-provider-kibana ~/.terraform.d/plugins/terraform-provider-kibana
	@echo "Build succeeded"

docker:
	cd docker/elasticsearch && docker build . -t elastic-local:6.0.0

build-gox: deps fmtcheck vet
	gox -osarch="linux/amd64 windows/amd64 darwin/amd64" \
	-output="pkg/{{.OS}}_{{.Arch}}/terraform-provider-kibana" .

release:
	go get github.com/goreleaser/goreleaser; \
    goreleaser; \

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

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile release docker