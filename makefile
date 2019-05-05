build:
	GOOS=linux go build

format:
	go fmt $(go list ./... | grep -v /openapi/)

test:
	mkdir -p coverage
	go test ./... -v -cover -covermode=count -coverprofile=coverage/profile

coverage: test
	go tool cover -html=coverage/profile -o coverage/index.html

cleanup:
	rm -rf coverage
