default: test

# Builds github.com/bosh-oneandone-cpi for linux-amd64
build:
	go build  -o out/cpi github.com/bosh-oneandone-cpi/main

# Build cross-platform binaries
build-all:
	gox -output="out/cpi_{{.OS}}_{{.Arch}}" -ldflags="-X github.com/bosh-oneandone-cpi/oneandone/config.cpiRelease=`cat release 2>/dev/null`" github.com/bosh-oneandone-cpi/main

# Prepration for tests
get-deps:
	# Go lint tool
	go get github.com/golang/lint/golint

	# Simplify cross-compiling
	go get github.com/mitchellh/gox

	# Ginkgo and omega test tools
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega

# Cleans up directory and source code with gofmt
clean:
	go clean ./...

# Run gofmt on all code
fmt:
	gofmt -l -w .

# Run linter with non-stric checking
lint:
	@echo golint ./... | wc -l

# Vet code
vet:
	go tool vet $$(ls -d */ | grep -v vendor)

# Runs the unit tests with coverage
test: get-deps clean fmt lint vet build
	ginkgo -r -race -skipPackage=integration .

# Runs the integration tests from Concourse
testintci: get-deps
	ginkgo integration -slowSpecThreshold=500 -progress -nodes=3 -randomizeAllSpecs -randomizeSuites $(GINKGO_ARGS) -v

# Runs the integration tests with coverage
testint: check-proj get-deps clean fmt
	$(eval INTEGRATION_ADDRESS = $(shell gcloud --project=$(oneandone_PROJECT) compute addresses describe cfintegration --region=us-central1 | head  -n1 | cut -f2 -d' '))
    
	CPI_ASYNC_DELETE=true STEMCELL_URL=https://storage.oneandoneapis.com/bosh-cpi-artifacts/bosh-stemcell-3262.12-oneandone-kvm-ubuntu-trusty-go_agent-raw.tar.gz SERVICE_ACCOUNT=cfintegration@$(oneandone_PROJECT).iam.gserviceaccount.com oneandone_PROJECT=$(oneandone_PROJECT) EXTERNAL_STATIC_IP=$(INTEGRATION_ADDRESS) ginkgo integration -slowSpecThreshold=500 -progress -nodes=3 -randomizeAllSpecs -randomizeSuites $(GINKGO_ARGS) -v
