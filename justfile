allure: clean test
	allure generate allure-results --clean
	allure open allure-report

clean:
	rm -rf allure-results

test:
	go test . -v
