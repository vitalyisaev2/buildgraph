test:
	overalls -project=github.com/vitalyisaev2/buildgraph -covermode=count -concurrency=2
	go tool cover -func=./overalls.coverprofile

clean:
	find . -type f -name "*.coverprofile" -exec rm -rf {} \;

migrations:
	cd ./storage/postgres/migrations/ && \
	go-bindata -o ./bindata.go -pkg migrations . && \
	cd -

build:
	go build .

run: build
	./buildgraph -c config/example.yaml

.PHONY: test
