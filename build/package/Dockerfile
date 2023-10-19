FROM golang:1.19-alpine3.17 as builder

# hadolint ignore=DL3018
RUN apk --no-cache add git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/soteria
RUN go build -o /soteria

FROM alpine:3.17

# hadolint ignore=DL3018
RUN apk --no-cache add ca-certificates tzdata && \
  mkdir /app

COPY --from=builder /soteria /app
WORKDIR /app

ENTRYPOINT ["/app/soteria" ]
CMD ["serve"]
