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

# Get test coverage
coverage:
	go test -v -coverpkg=./... -coverprofile=profile.cov ./...
	go tool cover -func profile.cov

coverage-html: coverage
	go tool cover -html profile.cov
