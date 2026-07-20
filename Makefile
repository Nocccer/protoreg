clean-cache:
	protoreg -clean-cache

generate: install
	go generate ./...

install:
	go install .

lint:
	golangci-lint-v2 run

lint-fix:
	golangci-lint-v2 run --fix

test: generate
	go test -v ./...

test-clean: clean-cache generate
	go test -v ./...