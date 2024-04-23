DIST   ?= development
DEBUG  ?= --debug

.DEFAULT_GOAL := debug
.PHONY: bump
.PHONY: vet
.PHONY: lint

all: test      \
	 benchmark \
     coverage

clean:
	go clean
	rm -rf bin

update:
	go get -u github.com/uhppoted/uhppote-core@master
	go get -u golang.org/x/sys
	go mod tidy

update-release:
	go get -u github.com/uhppoted/uhppote-core
	go get -u golang.org/x/sys
	go mod tidy

format: 
	go fmt ./...

build: format
	go build -trimpath ./...

test: build
	go test ./...

benchmark: build
	go test -bench ./...

coverage: build
	go test -cover ./...

vet: 
	go vet ./...

lint:
	env GOOS=darwin  GOARCH=amd64 staticcheck ./...
	env GOOS=linux   GOARCH=amd64 staticcheck ./...
	env GOOS=windows GOARCH=amd64 staticcheck ./...

vuln:
	govulncheck ./...

build-all: test vet lint
	env GOOS=linux   GOARCH=amd64       GOWORK=off go build -trimpath ./...
	env GOOS=linux   GOARCH=arm GOARM=7 GOWORK=off go build -trimpath ./...
	env GOOS=darwin  GOARCH=amd64       GOWORK=off go build -trimpath ./...
	env GOOS=windows GOARCH=amd64       GOWORK=off go build -trimpath ./...

release: build-all

publish: release
	echo "Releasing version $(VERSION)"
	gh release create "$(VERSION)" --draft --prerelease --title "$(VERSION)-beta" --notes-file release-notes.md

debug: build
	go test -v -run TestDeviceUnmarshal ./config/...

godoc:
	godoc -http=:80	-index_interval=60s
