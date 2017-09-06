test:
	overalls -project=github.com/vitalyisaev2/buildgraph -covermode=count -concurrency=2
	go tool cover -func=./overalls.coverprofile

clean:
	rm *coverprofile || true

build:
	go build .

.PHONY: test
