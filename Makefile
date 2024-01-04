linux: clean
	CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -v -o ./yaoapp/plugins/docloader.so
.PHONY: clean
clean: 
	rm -rf .tmp
	rm -rf dist
	rm -rf ./yaoapp/plugins/docloader.so