check: fmt lint update-readme test

# run tests
test:
	go test ./...

# format source code
fmt:
	golangci-lint fmt

# lint source code
lint:
	golangci-lint run --tests=false

[private]
update-readme:
	./update-usage.sh

# run specified example
example EXAMPLE:
	go test ./examples/{{ EXAMPLE }} -v -count=1 -tags example
