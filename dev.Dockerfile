FROM golang:alpine
ADD . /go/src/app
WORKDIR /go/src/app
CMD ["go", "run", "cmd/soteria/soteria.go", "serve"]
