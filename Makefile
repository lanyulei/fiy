PROJECT:=fiy

.PHONY: build
build:
	CGO_ENABLED=0 go build -o fiy main.go
build-sqlite:
	go build -tags sqlite3 -o fiy main.go
#.PHONY: test
#test:
#	go test -v ./... -cover

#.PHONY: docker
#docker:
#	docker build . -t fiy:latest
