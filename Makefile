.PHONY: test testci

test:
	go test -v `go list ./...|grep -v vendor`
