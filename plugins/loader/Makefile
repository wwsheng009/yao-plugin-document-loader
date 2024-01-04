linux: clean
	CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -v -o ../docloader.so
.PHONY: clean
clean: 
	rm -rf .tmp
	rm -rf dist
	rm -rf ../docloader.so