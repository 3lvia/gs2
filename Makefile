# PHONY used to mitigate conflict with dir name test
.PHONY: test
test:
	go fmt ./...
	go generate ./...
	go vet ./...
	golint ./...
	go test ./...

bench:
	go test ./... -bench=.

profile:
	go test ./gs2 -bench=BenchmarkDecoder_LargeMeterReadingFile -memprofile memprofile.out -cpuprofile cpuprofile.out
	go tool pprof --pdf cpuprofile.out > cpu.pdf
	go tool pprof --pdf memprofile.out > mem.pdf

parse:
	go run cmd/parser/main.go -file=testdata/timeseries.gs2

generate:
	go run cmd/generate-gs2/main.go

clean:
	rm -f mem.pdf cpu.pdf memprofile.out cpuprofile.out gs2.test
