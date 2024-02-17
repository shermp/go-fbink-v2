.PHONY: default
default: build.docker ;

.PHONY: build
build: export GOOS=linux
build: export GOARCH=arm
build: export CGO_ENABLED=1
build: export CC=${CROSS_COMPILE}-gcc
build: export CXX=${CROSS_COMPILE}-g++
build:
	go build -o build/app main.go

build.docker:
	docker build -t go-fbink-build .
	docker run -it --rm -v $(PWD)/..:/go/src/app --workdir /go/src/app/example go-fbink-build make build
