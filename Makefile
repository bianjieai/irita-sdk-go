PACKAGES=$(shell go list ./...)
PACKAGES_UNITTEST=$(shell go list ./... | grep -v integration_test)
export GO111MODULE = on

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs goimports -w -local github.com/bianjieai/irita-sdk-go

test-unit:
	@go test -v $(PACKAGES_UNITTEST)

test-integration:
	cd integration_test/scripts/ && sh build.sh && sh start.sh
	sleep 2s
	@go test -v $(PACKAGES)
	cd integration_test/scripts/ && sh clean.sh

proto-gen:
	@./third_party/protocgen.sh