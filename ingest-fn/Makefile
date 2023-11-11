.PHONY: build

build-IngestFunction:
	curl -s -L https://go.dev/dl/go1.21.4.linux-arm64.tar.gz -o go.tar.gz
	tar -xf go.tar.gz
	GOOS=linux GOARCH=arm64 ./go/bin/go build -o bootstrap
	cp ./bootstrap $(ARTIFACTS_DIR)/.
	cp template.gohtml $(ARTIFACTS_DIR)/.
	rm -rf go && rm go.tar.gz
