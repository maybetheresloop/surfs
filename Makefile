BIN_DIR=bin
MODULE=surfs

.PHONY: all
all: surfs-cli surfs-block surfs-meta

surfs-cli:
	go build -o bin/$@ -v ${MODULE}/cmd/cli

surfs-block:
	go build -o bin/$@ -v ${MODULE}/cmd/block

surfs-meta:
	go build -o bin/$@ -v ${MODULE}/cmd/meta

.PHONY: clean test

clean:
	rm -rf bin/

test:
	go test ${MODULE}/...