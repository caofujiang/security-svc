.PHONY: build clean tool lint help
export VERSION=1.1
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



linux_build_arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64  go build  -ldflags="-s -w " -trimpath  -o  security-svc-arm  main.go


#
#linux_build_arm:
#	make build_with_linux_arm upx package
#linux_image_arm:
#	#docker build   --build-arg GO_VERSION=1.20.5   -t caofujiang/secruity-svc:$(VERSION) .
#	docker build   -t caofujiang/secruity-svc:$(VERSION) .

build_linux_arm_with_arg:
	docker run --rm --privileged multiarch/qemu-user-static:register --reset
	docker run --rm \
		-v $(shell echo -n ${GOPATH}):/go \
		-w /go/src/github.com/chaosblade-io/chaosblade \
		-v ~/.m2/repository:/root/.m2/repository \
		-v $(shell pwd):/go/src/github.com/chaosblade-io/chaosblade \
		caofujiang/chaosblade-build-arm:latest




help:
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"
	@echo "make lint: golint ./..."
	@echo "make clean: remove object files and cached files"
