PHONY: vendor
vendor:
	go get github.com/izumin5210/gex/cmd/gex
	gex --build

PHONY: test
test:
	gotestcover -race -coverprofile=coverage.txt -v ./... 

PHONY: coverage
coverage:
	go tool cover -html=coverage.txt
