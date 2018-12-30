VERSION	:=	$(shell cat ./VERSION)

install:
	go install -v

fmt:
	go fmt ./...

test:
	cd ./lib && go test -v

release:
	git tag -a $(VERSION) -m "Release" || true
	git push origin $(VERSION)
	goreleaser --rm-dist

.PHONY: fmt install test release

