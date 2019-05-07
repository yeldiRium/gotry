build:
	GOOS=linux go build

format:
	go fmt $(go list ./... | grep -v /openapi/)

test:
	mkdir -p coverage
	go test ./... -v -cover -covermode=count -coverprofile=coverage.txt
	go tool cover -html=coverage.txt -o coverage/index.html

publish_coverage: test
	curl -s https://codecov.io/bash | bash

cleanup:
	rm -rf coverage
