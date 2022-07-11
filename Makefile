.PHONY: build
build:
	go build -o four-keys cmd/four-keys/main.go

.PHONY: run
run:
	make build
	./four-keys

.PHONY: test
test:
	go test -count=1 ./...
