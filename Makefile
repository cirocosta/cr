install:
	go install -v

fmt:
	go fmt
	cd ./lib && go fmt

test:
	cd ./lib && go test -v

.PHONY: fmt install test

