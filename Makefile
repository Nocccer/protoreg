lint:
	golangci-lint-v2 run

lint-fix:
	golangci-lint-v2 run --fix

generate:
	go generate ./...

test: generate
	go test -v ./...