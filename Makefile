linux: clean
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o ./yaoapp/plugins/docloader.so

windows: clean
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows GOHOSTOS=linux go build -o ./yaoapp/plugins/docloader.dll

.PHONY: clean
clean: 
	rm -rf .tmp
	rm -rf dist
	rm -rf ./yaoapp/plugins/docloader.*