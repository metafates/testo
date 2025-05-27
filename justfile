check: fmt lint

allure: clean test
	allure generate allure-results --clean
	allure open allure-report

clean:
	rm -rf allure-results

test:
	go test ./... -v

# format source code
fmt:
	golangci-lint fmt

# lint source code
lint:
	golangci-lint run --tests=false
