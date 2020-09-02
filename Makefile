IMAGE_FLAG := $(shell git rev-parse --abbrev-ref HEAD | tr '/' '-') #get image flag from current branch name.
IMAGE_TAG  ?= registry.apps.private.teh-1.snappcloud.io/dispatching-staging/soteria:${IMAGE_FLAG}

dependencies:
	go mod vendor -v

compile: dependencies
	go build --ldflags "-linkmode external -extldflags '-static'" -mod vendor -v  cmd/soteria/soteria.go

test:
	go test --race -gcflags=-l -v -coverprofile .coverage.out.tmp ./...
	cat .coverage.out.tmp | grep -v "mock.go" > .coverage.out
	rm -rf .coverage.out.tmp
	go tool cover -func .coverage.out

html-report: test
	go tool cover -html=.coverage.out -o coverage.html

build-image: compile
	docker build -t ${IMAGE_TAG} .

push-image: build-image
    docker push ${IMAGE_TAG}

build-image-dev: compile
	docker build -t soteria:latest .
