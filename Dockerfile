FROM alpine:3.14
ARG BUILD_DATE
ARG VCS_REF
ARG BUILD_VERSION

# hadolint ignore=DL3018
RUN echo "https://repo.snapp.tech/repository/alpine/v3.14/main" > /etc/apk/repositories && \
    echo "https://repo.snapp.tech/repository/alpine/v3.14/community" >> /etc/apk/repositories && \
    apk --no-cache --update add ca-certificates tzdata && \
    mkdir /app

COPY ./soteria /app
WORKDIR /app

ENV SOTERIA_BUILD_DATE=${BUILD_DATE}
ENV SOTERIA_VCS_REF=${VCS_REF}
ENV SOTERIA_BUILD_VERSION=${BUILD_VERSION}

CMD ["/app/soteria", "serve"]
