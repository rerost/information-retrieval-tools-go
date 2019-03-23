PHONY: vendor
vendor:
	go get github.com/izumin5210/gex/cmd/gex
	gex --build

PHONY: test
test: vendor
	gotestcover -race -coverprofile=coverage.txt -v ./... 

PHONY: coverage
coverage:
	go tool cover -html=coverage.txt

PHONY: bench
bench: vendor
	go test -bench=. -race -v ./...

PHONY: profile
profile:
	: > mkdir pprof 
	# interleaving
	: > mkdir pprof/interleaving
	go test -bench=BenchmarkPerform -benchmem -o pprof/interleaving/test.bin -cpuprofile pprof/interleaving/cpu.out ./interleaving/interleaving_bench_test.go
	go tool pprof --svg pprof/interleaving/test.bin pprof/interleaving/cpu.out > pprof/interleaving/profile.svg
