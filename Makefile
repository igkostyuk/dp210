all: lint test 
# ==============================================================================
# Lint all files in the project

lint:
	@golangci-lint run -c .golangci.yml

# ==============================================================================
# Running tests within the local computer

test:
	go test ./... -count=1

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

