.PHONY: build
build:
	bash -c "go build -ldflags '-X main.version=$$(git tag --sort=-taggerdate | head -1)' -o four-keys cmd/four-keys/main.go"

.PHONY: run
run:
	make build
	./four-keys

.PHONY: test
test:
	go test -count=1 ./...
