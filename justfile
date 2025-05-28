check: fmt lint

# run tests
test:
	go test ./... -v

# format source code
fmt:
	golangci-lint fmt

# lint source code
lint:
	golangci-lint run --tests=false
