FROM alpine
ARG BUILD_DATE
ARG VCS_REF
ARG BUILD_VERSION
RUN echo -e "https://repo.snapp.tech/repository/alpine/v3.12/main\nhttps://repo.snapp.tech/repository/alpine/v3.12/community" > /etc/apk/repositories

RUN apk --no-cache --update add ca-certificates

RUN mkdir /app

COPY ./soteria /app
WORKDIR /app
ENV PRIVENT_BUILD_DATE=${BUILD_DATE}
ENV PRIVENT_VCS_REF=${VCS_REF}
ENV PRIVENT_BUILD_VERSION=${BUILD_VERSION}
CMD ["/app/soteria", "serve"]