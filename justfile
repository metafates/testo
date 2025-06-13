check: fmt lint test

# run tests
test:
	go test ./...

# format source code
fmt:
	golangci-lint fmt

# lint source code
lint:
	golangci-lint run --tests=false

# run specified example
example EXAMPLE:
	go test ./examples/{{ EXAMPLE }} -v -count=1 -tags example
