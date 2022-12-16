DIST   ?= development
DEBUG  ?= --debug

.PHONY: bump

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

vet: build
	go vet ./...

lint: build
	golint ./...

benchmark: build
	go test -bench ./...

coverage: build
	go test -cover ./...

build-all: test vet
	env GOOS=linux   GOARCH=amd64       GOWORK=off go build -trimpath ./...
	env GOOS=linux   GOARCH=arm GOARM=7 GOWORK=off go build -trimpath ./...
	env GOOS=darwin  GOARCH=amd64       GOWORK=off go build -trimpath ./...
	env GOOS=windows GOARCH=amd64       GOWORK=off go build -trimpath ./...

release: build-all

publish: release
	echo "Releasing version $(VERSION)"
	gh release create "$(VERSION)" --draft --prerelease --title "$(VERSION)-beta" --notes-file release-notes.md

debug: build
	env GOOS=windows GOARCH=amd64       GOWORK=off go build -trimpath ./...

godoc:
	godoc -http=:80	-index_interval=60s
