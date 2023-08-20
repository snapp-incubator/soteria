FROM alpine:3.18

# hadolint ignore=DL3018
RUN echo "https://repo.snapp.tech/repository/alpine/v3.18/main" > /etc/apk/repositories && \
  echo "https://repo.snapp.tech/repository/alpine/v3.18/community" >> /etc/apk/repositories && \
  apk --no-cache --update add ca-certificates tzdata && \
  mkdir /app

COPY ./soteria /app
WORKDIR /app

CMD ["/app/soteria", "serve"]
