FROM golang:1.24-alpine3.21 AS builder

# hadolint ignore=DL3018
RUN apk --no-cache add git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/soteria
RUN go build -o /soteria

FROM alpine:3.20

# hadolint ignore=DL3018
RUN apk --no-cache add ca-certificates tzdata && \
  mkdir /app

COPY --from=builder /soteria /app
WORKDIR /app

ENTRYPOINT ["/app/soteria" ]
CMD ["serve"]
