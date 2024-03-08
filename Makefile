.PHONY: build clean tool lint help
export VERSION=1.0
all: build

build:
	@go build -v .

tool:
	go vet ./...; true
	gofmt -w .

lint:
	golint ./...

clean:
	rm -rf go-gin-example
	go clean -i .

linux_build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build  -ldflags="-s -w " -trimpath  -o  security-svc-linux  main.go
linux_image:
	docker build -t caofujiang/security-svc:$(VERSION) .

#
#linux_build_arm:
#	CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build  -ldflags="-s -w " -trimpath  -o  security-svc-linux  main.go
#linux_image_arm:
#	#docker build   --build-arg GO_VERSION=1.20.5   -t caofujiang/secruity-svc:$(VERSION) .
#	docker build   -t caofujiang/secruity-svc:$(VERSION) .

help:
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"
	@echo "make lint: golint ./..."
	@echo "make clean: remove object files and cached files"
