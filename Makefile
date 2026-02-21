.PHONY: run dev build

run:
	go run ./cmd/spec-viewer serve --folder ./examples/specs
build:
	go build -o bin/spec-viewer ./cmd/spec-viewer
