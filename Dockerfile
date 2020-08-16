FROM golang:alpine AS build-env
LABEL maintainer="Alexander Makhinov <contact@monaxgt.com>" \
      repository="https://github.com/Rostelecom-CERT/maxanon"

COPY . /go/src/github.com/Rostelecom-CERT/maxanon

RUN apk add --no-cache git mercurial \
    && cd /go/src/github.com/Rostelecom-CERT/maxanon/service/maxanon \
    && go get -t . \
    && CGO_ENABLED=0 go build -ldflags="-s -w" \
                              -a \
                              -installsuffix static \
                              -o /maxanon

FROM alpine:3.12

RUN apk --update --no-cache add ca-certificates curl \
  && adduser -h /app -D app \
  && mkdir -p /app/data \
  && chown -R app /app

COPY --from=build-env /maxanon /app/maxanon

USER app

VOLUME /app/data

WORKDIR /app

ENTRYPOINT ["./maxanon"]