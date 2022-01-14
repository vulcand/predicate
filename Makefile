.PHONY: check test

default: check test

test:
	go test -v ./... -cover

check:
	golangci-lint run

sloccount:
	 find . -name "*.go" -print0 | xargs -0 wc -l
