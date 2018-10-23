.PHONY: test testci

test:
	go test ./... -v

testci:
	go test ./... -v -tags="docker"
