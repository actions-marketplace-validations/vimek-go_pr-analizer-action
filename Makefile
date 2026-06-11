.PHONY: format lint test-coverage-out

init:
	go mod download
	go install github.com/segmentio/golines@latest
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.12.2
	go install github.com/vektra/mockery/v2@latest


format:
	golines --base-formatter="goimports" -w -m 120 .
	gofumpt -w .

lint: ## Run Go linter locally
	golangci-lint version
	golangci-lint -c ".golangci.yml" run ./...

test-coverage-out:
	ENV=test go test -race -coverprofile=profile.cov -covermode=atomic ./...
	go tool cover -func profile.cov