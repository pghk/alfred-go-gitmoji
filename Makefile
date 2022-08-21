.PHONY: test build

test:
	alfred_workflow_data=../test/_data \
	alfred_workflow_cache=../test/_cache \
	alfred_workflow_bundleid=alfred-go-gitmoji.test \
	go test ./...

clean:
	go clean
	-rm -r ./build/*

build:
	make test
	make clean
	go build -o build ./cmd/gitmoji
	cp ./configs/* build/
	cp -r ./assets/* build/
